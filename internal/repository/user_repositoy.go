package repository

import (
	"context"
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserFilter struct {
	Search string
	Role   string
}

type UserRepository interface {
	GetAll(ctx context.Context, loadProfile bool, filter *UserFilter) ([]model.User, error)
	GetByID(ctx context.Context, id string, loadProfile bool) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User, profile *model.UserProfile) error
	Update(ctx context.Context, user *model.User, profile *model.UserProfile) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context, loadProfile bool, filter *UserFilter) ([]model.User, error) {
	var users []model.User
	query := r.db.WithContext(ctx)
	if loadProfile {
		query = query.Preload("Profile", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at DESC")
		})
	}
	if filter != nil {
		if filter.Search != "" {
			like := "%" + filter.Search + "%"
			query = query.Where("full_name LIKE ? OR email LIKE ?", like, like)
		}
		if filter.Role != "" {
			query = query.Where("role = ?", filter.Role)
		}
	}
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) GetByID(ctx context.Context, id string, loadProfile bool) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if loadProfile {
		var profile model.UserProfile
		if err := r.db.WithContext(ctx).
			Where("user_id = ?", user.ID).
			First(&profile).Error; err == nil {
			user.Profile = &profile
		}
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ?", strings.ToLower(strings.TrimSpace(email))).
		Take(&user).Error
	if err != nil {
		return nil, err
	}
	var profile model.UserProfile
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user.ID).
		Order("updated_at DESC").
		First(&profile).Error; err == nil {
		user.Profile = &profile
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if user.ID == "" {
			user.ID = uuid.NewString()
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if profile == nil {
			profile = &model.UserProfile{}
		}
		profile.UserID = user.ID
		if err := tx.Save(profile).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) Update(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).
			Where("id = ?", user.ID).
			Updates(map[string]any{
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  user.Password,
				"role":      user.Role,
			}).Error; err != nil {
			return err
		}

		if profile != nil {
			profile.UserID = user.ID
			if profile.ID == "" {
				var existing model.UserProfile
				if err := tx.Where("user_id = ?", user.ID).First(&existing).Error; err == nil {
					profile.ID = existing.ID
				}
			}
			if err := tx.Save(profile).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}
