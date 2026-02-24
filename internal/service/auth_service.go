package service

import (
	"errors"
	"strings"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
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
		if isDuplicateEmailError(err) {
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

	token, err := utils.GenerateAccessToken(config.GetEnv().AccessKey, utils.TokenClaims{
		Subject: user.ID,
		Email:   user.Email,
		Role:    user.Role,
		Exp:     time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

var (
	ErrEmailAlreadyRegistered = errors.New("Email is already registered")
	ErrInvalidCredentials     = errors.New("Email or password is invalid, please try again")
)

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

func isDuplicateEmailError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "1062")
}
