package service

import (
	"context"
	"errors"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, payload request.AuthRegisterRequest) (*model.User, error)
	Login(ctx context.Context, payload request.AuthLoginRequest) (*model.User, string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
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
