package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, payload request.AuthRegisterRequest) (*model.User, error)
	Login(ctx context.Context, payload request.AuthLoginRequest) (*model.User, string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, payload request.ResetPasswordRequest) error
}

type authService struct {
	userRepo    repository.UserRepository
	userOtpRepo repository.UserOTPRepository
}

func NewAuthService(userRepo repository.UserRepository, userOtpRepo repository.UserOTPRepository) AuthService {
	return &authService{userRepo: userRepo, userOtpRepo: userOtpRepo}
}

func (s *authService) Register(ctx context.Context, payload request.AuthRegisterRequest) (*model.User, error) {
	if _, err := s.userRepo.GetByEmail(ctx, payload.Email); err == nil {
		return nil, ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.NewString(),
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: hashed,
		Role:     model.UserRoleUser,
	}

	// Create empty profile row so later updates are straightforward.
	if err := s.userRepo.Create(ctx, user, nil); err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, payload request.AuthLoginRequest) (*model.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if !utils.CheckPassword(payload.Password, user.Password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Jangan ungkap apakah email terdaftar atau tidak
		return nil
	}

	now := time.Now()
	record, err := s.userOtpRepo.GetByUserID(ctx, user.ID)
	if err != nil && !isNotFoundError(err) {
		return err
	}

	if record != nil {
		// Cek apakah masih freeze
		if record.FrozenUntil != nil && record.FrozenUntil.After(now) {
			remaining := int(record.FrozenUntil.Sub(now).Minutes()) + 1
			return common.BadRequestError(fmt.Sprintf("Terlalu banyak permintaan. Coba lagi dalam %d menit.", remaining))
		}
		// Reset count jika freeze sudah berakhir
		if record.FrozenUntil != nil && !record.FrozenUntil.After(now) {
			record.ResendCount = 0
			record.FrozenUntil = nil
		}
		// Freeze jika sudah 3x
		if record.ResendCount >= 3 {
			frozenUntil := now.Add(10 * time.Minute)
			record.FrozenUntil = &frozenUntil
			_ = s.userOtpRepo.Save(ctx, record)
			return common.BadRequestError("Terlalu banyak permintaan OTP. Coba lagi dalam 10 menit.")
		}
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := now.Add(10 * time.Minute)
	resendCount := 1
	if record != nil {
		resendCount = record.ResendCount + 1
	}

	newRecord := &model.UserOTP{
		UserID:      user.ID,
		OTPCode:     otp,
		ExpiresAt:   expiresAt,
		ResendCount: resendCount,
	}
	if record != nil {
		newRecord.ID = record.ID
	}

	if err := s.userOtpRepo.Save(ctx, newRecord); err != nil {
		return err
	}

	go s.sendOTPEmail(user, otp)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, payload request.ResetPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		return common.BadRequestError("Email tidak ditemukan")
	}

	record, err := s.userOtpRepo.GetByUserID(ctx, user.ID)
	if err != nil || record == nil {
		return common.BadRequestError("OTP tidak valid atau sudah kadaluarsa")
	}

	if time.Now().After(record.ExpiresAt) {
		return common.BadRequestError("OTP sudah kadaluarsa")
	}

	if record.OTPCode != payload.OTP {
		return common.BadRequestError("OTP tidak valid")
	}

	hashed, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	if err := s.userRepo.Update(ctx, user, nil); err != nil {
		return err
	}

	return s.userOtpRepo.DeleteByUserID(ctx, user.ID)
}

func (s *authService) sendOTPEmail(user *model.User, otp string) {
	env := config.GetEnv()
	_ = utils.SendEmail(utils.EmailParams{
		FromName:   env.SMTPFromName,
		FromEmail:  env.SMTPFromEmail,
		Password:   env.SMTPPassword,
		Host:       env.SMTPHost,
		Port:       env.SMTPPort,
		Encryption: utils.EncryptionType(env.SMTPEncryption),
		ToName:     user.FullName,
		ToEmail:    user.Email,
		Subject:    "Kode OTP Reset Password - King Royal",
		Template:   "templates/otp_email.html",
		Data: map[string]any{
			"name": user.FullName,
			"otp":  otp,
		},
	})
}

type DeletedAccountError struct {
	DeletedAt time.Time
	Email     string
}

func NewDeletedAccountError(deletedAt time.Time, email string) *DeletedAccountError {
	return &DeletedAccountError{DeletedAt: deletedAt, Email: email}
}

func (e *DeletedAccountError) Error() string {
	return "Account has been deleted"
}
