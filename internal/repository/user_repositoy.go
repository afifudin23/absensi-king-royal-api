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
	GetByID(id string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	Create(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id string) (bool, error)
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

func (r *userRepository) GetByID(id string) (model.User, error) {
	var user model.User
	err := r.db.
		Where("id = ?", id).
		Take(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.
		Where("email = ?", strings.ToLower(strings.TrimSpace(email))).
		Take(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	if strings.TrimSpace(user.Email) == "" {
		return model.User{}, errors.New("User is required")
	}
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) Update(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *userRepository) Delete(id string) (bool, error) {
	err := r.db.Delete(&model.User{}, "id = ?", id).Error
	return err == nil, err
}
