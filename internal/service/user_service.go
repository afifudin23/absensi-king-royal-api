package service

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/google/uuid"
)

type UserService interface {
	GetAll(ctx context.Context, filter *repository.UserFilter) ([]model.User, error)
	ResetPassword(ctx context.Context, adminID, targetUserID, newPassword string) error
	Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error)
	UpdateProfile(ctx context.Context, userID string, payload request.UserUpdateProfileRequest) (*model.User, error)
	Delete(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID string, payload request.ChangePasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
	fileRepo repository.FileRepository
}

func NewUserService(userRepo repository.UserRepository, fileRepo repository.FileRepository) UserService {
	return &userService{userRepo: userRepo, fileRepo: fileRepo}
}

func (s *userService) GetAll(ctx context.Context, filter *repository.UserFilter) ([]model.User, error) {
	return s.userRepo.GetAll(ctx, true, filter)
}

func (s *userService) ResetPassword(ctx context.Context, adminID, targetUserID, newPassword string) error {
	if adminID == targetUserID {
		return common.BadRequestError("Cannot reset your own password via this endpoint")
	}
	user, err := s.userRepo.GetByID(ctx, targetUserID, false)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) || isNotFoundError(err) {
			return ErrUserNotFound
		}
		return err
	}
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.userRepo.Update(ctx, user, nil)
}

func generatePassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[n.Int64()]
	}
	return string(result)
}

func (s *userService) Create(ctx context.Context, payload request.UserCreateRequest) (*model.User, error) {
	rawPassword := generatePassword(8)

	hashedPassword, err := utils.HashPassword(rawPassword)
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
		BirthDate:         parseDateString(payload.BirthDate),
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

	go func() {
		env := config.GetEnv()
		if env == nil {
			return
		}
		err := utils.SendEmail(utils.EmailParams{
			FromName:   env.SMTPFromName,
			FromEmail:  env.SMTPFromEmail,
			Password:   env.SMTPPassword,
			Host:       env.SMTPHost,
			Port:       env.SMTPPort,
			Encryption: utils.EncryptionType(env.SMTPEncryption),
			ToName:     user.FullName,
			ToEmail:    user.Email,
			Subject:    "Selamat Datang - Akun Absensi King Royal",
			Template:   "templates/welcome_email.html",
			Data: map[string]any{
				"name":     user.FullName,
				"email":    user.Email,
				"password": rawPassword,
			},
		})
		if err != nil {
			log.Printf("failed to send welcome email to %s: %v", user.Email, err)
		}
	}()

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID, true)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, userID string, payload request.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID, true)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	profile := ensureUserProfile(user)
	applyUserUpdateRequest(user, profile, payload)

	if payload.ClearProfilePicture != nil && *payload.ClearProfilePicture {
		profile.ProfilePictureID = nil
		profile.ProfilePictureURL = nil
	} else if payload.ProfilePictureID != nil {
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
	user, err := s.userRepo.GetByID(ctx, userID, true)
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

func (s *userService) ChangePassword(ctx context.Context, userID string, payload request.ChangePasswordRequest) error {
	// Pastikan new_password dan confirm_password sama
	if payload.NewPassword != payload.ConfirmPassword {
		return common.BadRequestError("New password and confirm password do not match")
	}

	user, err := s.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		if isNotFoundError(err) {
			return ErrUserNotFound
		}
		return err
	}

	// Verifikasi password lama sebelum ganti
	if !utils.CheckPassword(payload.CurrentPassword, user.Password) {
		return common.BadRequestError("Current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user, nil)
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
	if payload.Email != nil {
		existing.Email = *payload.Email
	}
	if payload.Role != nil {
		existing.Role = model.UserRole(*payload.Role)
	}
	if payload.JoinedAt != nil {
		profile.JoinedAt = parseDateString(payload.JoinedAt)
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

func parseDateString(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, *s); err == nil {
			return &t
		}
	}
	return nil
}

func applyUserUpdate(
	profile *model.UserProfile,
	employeeCode *string,
	employmentStatus *model.UserEmploymentStatus,
	birthPlace *string,
	birthDate *string,
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
		profile.BirthDate = parseDateString(birthDate)
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
