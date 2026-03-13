package service

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
)

type UserService interface {
	GetAllUsers() ([]model.User, error)
	CreateUser(payload model.User) (*model.User, error)
	GetUserByID(userID string) (*model.User, error)
	UpdateUser(userID string, payload model.User) (*model.User, error)
	DeleteUser(userID string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepository()}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) CreateUser(payload model.User) (*model.User, error) {
	if payload.Password != "" {
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			return nil, err
		}
		payload.Password = hashedPassword
	}

	user := &payload
	if err := s.userRepo.Create(user); err != nil {
		if isDuplicateError(err) {
			return nil, common.BadRequestError("Email is already registered")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.BadRequestError("User not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(userID string, payload model.User) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.BadRequestError("User not found")
		}
		return nil, err
	}

	if payload.Password != "" {
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			return nil, err
		}
		payload.Password = hashedPassword
	}

	applyUserUpdates(user, payload)
	if err := s.userRepo.Update(user); err != nil {
		if isDuplicateError(err) {
			return nil, common.BadRequestError("Email is already registered")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(userID string) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if isNotFoundError(err) {
			return common.BadRequestError("User not found")
		}
		return err
	}
	return s.userRepo.Delete(userID)
}

func applyUserUpdates(existing *model.User, payload model.User) {
	if payload.FullName != "" {
		existing.FullName = payload.FullName
	}
	if payload.Email != "" {
		existing.Email = payload.Email
	}
	if payload.Password != "" {
		existing.Password = payload.Password
	}
	if payload.Role != "" {
		existing.Role = payload.Role
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
	if payload.ProfilePictureURL != nil {
		existing.ProfilePictureURL = payload.ProfilePictureURL
	}
	if payload.ProfilePictureID != nil {
		existing.ProfilePictureID = payload.ProfilePictureID
	}
}
