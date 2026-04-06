package service

import (
	"context"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/google/uuid"
)

type UserService interface {
	GetAll(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error)
	UpdateProfile(ctx context.Context, userID string, payload request.UserUpdateProfileRequest) (*model.User, error)
	Delete(ctx context.Context, userID string) error
	// SendEmail()
}

type userService struct {
	userRepo repository.UserRepository
	fileRepo repository.FileRepository
}

func NewUserService(userRepo repository.UserRepository, fileRepo repository.FileRepository) UserService {
	return &userService{userRepo: userRepo, fileRepo: fileRepo}
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAll(ctx, false)
}

func (s *userService) Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.NewString(),
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     payload.Role,
	}

	profile := &model.UserProfile{
		UserID:            user.ID,
		EmployeeCode:      payload.EmployeeCode,
		EmploymentStatus:  payload.EmploymentStatus,
		BirthPlace:        payload.BirthPlace,
		BirthDate:         payload.BirthDate,
		Gender:            payload.Gender,
		Address:           payload.Address,
		PhoneNumber:       payload.PhoneNumber,
		Position:          payload.Position,
		Department:        payload.Department,
		BankAccountNumber: payload.BankAccountNumber,
		BasicSalary:       payload.BasicSalary,
		PositionAllowance: payload.PositionAllowance,
		OtherAllowance:    payload.OtherAllowance,
	}

	if payload.ProfilePictureID != nil {
		if s.fileRepo == nil {
			return nil, common.InternalServerError()
		}

		file, err := s.fileRepo.GetByID(ctx, *payload.ProfilePictureID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid profile_picture_id")
			}
			return nil, err
		}
		if file.Type != model.FileTypeProfilePicture {
			return nil, common.BadRequestError("Invalid file type for profile picture")
		}

		url := file.FileURL
		profile.ProfilePictureID = payload.ProfilePictureID
		profile.ProfilePictureURL = &url
	}

	user.Profile = profile

	if err := s.userRepo.Create(ctx, user, profile); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	profile := ensureUserProfile(user)
	applyUserUpdateRequest(user, profile, payload)

	if payload.ProfilePictureID != nil {
		if s.fileRepo == nil {
			return nil, common.InternalServerError()
		}

		file, err := s.fileRepo.GetByID(ctx, *payload.ProfilePictureID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid profile_picture_id")
			}
			return nil, err
		}
		if file.Type != model.FileTypeProfilePicture {
			return nil, common.BadRequestError("Invalid file type for profile picture")
		}
		if file.UploadedBy != userID {
			return nil, common.ForbiddenError("File does not belong to current user")
		}

		url := file.FileURL
		profile.ProfilePictureID = payload.ProfilePictureID
		profile.ProfilePictureURL = &url
	}

	if err := s.userRepo.Update(ctx, user, profile); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	user.Profile = profile
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, payload request.UserUpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	profile := ensureUserProfile(user)

	if payload.Password != nil && *payload.Password != "" {
		hashedPassword, err := utils.HashPassword(*payload.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	applyUserUpdateProfileRequest(user, profile, payload)

	if payload.ProfilePictureID != nil {
		if s.fileRepo == nil {
			return nil, common.InternalServerError()
		}

		file, err := s.fileRepo.GetByID(ctx, *payload.ProfilePictureID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid profile_picture_id")
			}
			return nil, err
		}
		if file.Type != model.FileTypeProfilePicture {
			return nil, common.BadRequestError("Invalid file type for profile picture")
		}
		if file.UploadedBy != userID {
			return nil, common.ForbiddenError("File does not belong to current user")
		}

		url := file.FileURL
		profile.ProfilePictureID = payload.ProfilePictureID
		profile.ProfilePictureURL = &url
	}

	if err := s.userRepo.Update(ctx, user, profile); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	user.Profile = profile
	return user, nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	_, err := s.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		if isNotFoundError(err) {
			return ErrUserNotFound
		}
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}

func ensureUserProfile(user *model.User) *model.UserProfile {
	if user.Profile == nil {
		user.Profile = &model.UserProfile{UserID: user.ID}
	}
	return user.Profile
}

func applyUserUpdateRequest(existing *model.User, profile *model.UserProfile, payload request.UserUpdateRequest) {
	if payload.FullName != nil {
		existing.FullName = *payload.FullName
	}
	if payload.Role != nil {
		existing.Role = model.UserRole(*payload.Role)
	}

	applyUserUpdate(profile,
		payload.EmployeeCode,
		payload.EmploymentStatus,
		payload.BirthPlace,
		payload.BirthDate,
		payload.Gender,
		payload.Address,
		payload.PhoneNumber,
		payload.Position,
		payload.Department,
		payload.BankAccountNumber,
		payload.BasicSalary,
		payload.PositionAllowance,
		payload.OtherAllowance,
	)
}

func applyUserUpdateProfileRequest(existing *model.User, profile *model.UserProfile, payload request.UserUpdateProfileRequest) {
	if payload.FullName != nil {
		existing.FullName = *payload.FullName
	}
	if payload.Email != nil {
		existing.Email = *payload.Email
	}
	if payload.Role != nil {
		existing.Role = *payload.Role
	}

	applyUserUpdate(profile,
		payload.EmployeeCode,
		payload.EmploymentStatus,
		payload.BirthPlace,
		payload.BirthDate,
		payload.Gender,
		payload.Address,
		payload.PhoneNumber,
		payload.Position,
		payload.Department,
		payload.BankAccountNumber,
		nil,
		nil,
		nil,
	)
}

func applyUserUpdate(
	profile *model.UserProfile,
	employeeCode *string,
	employmentStatus *model.UserEmploymentStatus,
	birthPlace *string,
	birthDate *time.Time,
	gender *model.UserGender,
	address *string,
	phoneNumber *string,
	position *string,
	department *string,
	bankAccountNumber *string,
	basicSalary *float64,
	positionAllowance *float64,
	otherAllowance *float64,
) {
	// NOTE: keep request/response schema unchanged; only persistence changes.
	if employeeCode != nil {
		profile.EmployeeCode = employeeCode
	}
	if employmentStatus != nil {
		profile.EmploymentStatus = employmentStatus
	}
	if birthPlace != nil {
		profile.BirthPlace = birthPlace
	}
	if birthDate != nil {
		profile.BirthDate = birthDate
	}
	if gender != nil {
		profile.Gender = gender
	}
	if address != nil {
		profile.Address = address
	}
	if phoneNumber != nil {
		profile.PhoneNumber = phoneNumber
	}
	if position != nil {
		profile.Position = position
	}
	if department != nil {
		profile.Department = department
	}
	if bankAccountNumber != nil {
		profile.BankAccountNumber = bankAccountNumber
	}
	if basicSalary != nil {
		profile.BasicSalary = basicSalary
	}
	if positionAllowance != nil {
		profile.PositionAllowance = positionAllowance
	}
	if otherAllowance != nil {
		profile.OtherAllowance = otherAllowance
	}
}
