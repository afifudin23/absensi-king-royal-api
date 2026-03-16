package service

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
)

type UserService interface {
	GetAll(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error)
	UpdateProfile(ctx context.Context, userID string, payload request.UserUpdateProfileRequest) (*model.User, error)
	Delete(ctx context.Context, userID string) error
}

type userService struct {
	userRepo repository.UserRepository
	fileRepo repository.FileRepository
}

func NewUserService(userRepo repository.UserRepository, fileRepo repository.FileRepository) UserService {
	return &userService{userRepo: userRepo, fileRepo: fileRepo}
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	var profilePictureURL *string
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
		profilePictureURL = &url
	}

	user := &model.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     payload.Role,

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
		ProfilePictureURL: profilePictureURL,
		ProfilePictureID:  payload.ProfilePictureID,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	applyUserUpdateRequest(user, payload)

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
		user.ProfilePictureID = payload.ProfilePictureID
		user.ProfilePictureURL = &url
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, payload request.UserUpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if payload.Password != nil && *payload.Password != "" {
		hashedPassword, err := utils.HashPassword(*payload.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	applyUserUpdateProfileRequest(user, payload)

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
		user.ProfilePictureID = payload.ProfilePictureID
		user.ProfilePictureURL = &url
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if isNotFoundError(err) {
			return ErrUserNotFound
		}
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}

func applyUserUpdateRequest(existing *model.User, payload request.UserUpdateRequest) {
	if payload.FullName != nil {
		existing.FullName = *payload.FullName
	}
	if payload.Role != nil {
		existing.Role = model.UserRole(*payload.Role)
	}

	if payload.EmployeeCode != nil {
		existing.EmployeeCode = payload.EmployeeCode
	}
	if payload.EmploymentStatus != nil {
		existing.EmploymentStatus = payload.EmploymentStatus
	}
	if payload.BirthPlace != nil {
		existing.BirthPlace = payload.BirthPlace
	}
	if payload.BirthDate != nil {
		existing.BirthDate = payload.BirthDate
	}
	if payload.Gender != nil {
		existing.Gender = payload.Gender
	}
	if payload.Address != nil {
		existing.Address = payload.Address
	}
	if payload.PhoneNumber != nil {
		existing.PhoneNumber = payload.PhoneNumber
	}
	if payload.Position != nil {
		existing.Position = payload.Position
	}
	if payload.Department != nil {
		existing.Department = payload.Department
	}
	if payload.BankAccountNumber != nil {
		existing.BankAccountNumber = payload.BankAccountNumber
	}
	if payload.BasicSalary != nil {
		existing.BasicSalary = payload.BasicSalary
	}
}

func applyUserUpdateProfileRequest(existing *model.User, payload request.UserUpdateProfileRequest) {
	if payload.FullName != nil {
		existing.FullName = *payload.FullName
	}
	if payload.Email != nil {
		existing.Email = *payload.Email
	}
	if payload.Role != nil {
		existing.Role = *payload.Role
	}

	if payload.EmployeeCode != nil {
		existing.EmployeeCode = payload.EmployeeCode
	}
	if payload.EmploymentStatus != nil {
		existing.EmploymentStatus = payload.EmploymentStatus
	}
	if payload.BirthPlace != nil {
		existing.BirthPlace = payload.BirthPlace
	}
	if payload.BirthDate != nil {
		existing.BirthDate = payload.BirthDate
	}
	if payload.Gender != nil {
		existing.Gender = payload.Gender
	}
	if payload.Address != nil {
		existing.Address = payload.Address
	}
	if payload.PhoneNumber != nil {
		existing.PhoneNumber = payload.PhoneNumber
	}
	if payload.Position != nil {
		existing.Position = payload.Position
	}
	if payload.Department != nil {
		existing.Department = payload.Department
	}
	if payload.BankAccountNumber != nil {
		existing.BankAccountNumber = payload.BankAccountNumber
	}
}
