package repository

import (
	"errors"
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: config.GetDB()}
}

func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.
		Where("id = ?", id).
		Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.
		Where("email = ?", strings.ToLower(strings.TrimSpace(email))).
		Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	if user == nil || strings.TrimSpace(user.Email) == "" {
		return errors.New("User is required")
	}
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}
