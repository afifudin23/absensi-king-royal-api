package service

import (
	"errors"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(payload request.AuthRegisterRequest) (*model.User, error)
	Login(payload request.AuthLoginRequest) (*model.User, string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService() AuthService {
	return &authService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *authService) Register(payload request.AuthRegisterRequest) (*model.User, error) {
	if _, err := s.userRepo.GetByEmail(payload.Email); err == nil {
		return nil, ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: string(hashed),
		Role:     "user",
	}

	user, err = s.userRepo.Create(user)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(payload request.AuthLoginRequest) (*model.User, string, error) {
	user, err := s.userRepo.GetByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if user.DeletedAt != nil {
		return nil, "", NewDeletedAccountError(*user.DeletedAt, user.Email)
	}

	if !utils.CheckPassword(payload.Password, user.Password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
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
