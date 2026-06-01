package repository

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserOTPRepository interface {
	GetByUserID(ctx context.Context, userID string) (*model.UserOTP, error)
	Save(ctx context.Context, otp *model.UserOTP) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type userOTPRepository struct {
	db *gorm.DB
}

func NewUserOTPRepository(db *gorm.DB) UserOTPRepository {
	return &userOTPRepository{db: db}
}

func (r *userOTPRepository) GetByUserID(ctx context.Context, userID string) (*model.UserOTP, error) {
	var otp model.UserOTP
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Take(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *userOTPRepository) Save(ctx context.Context, otp *model.UserOTP) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"otp_code", "expires_at", "resend_count", "frozen_until", "updated_at",
			}),
		}).
		Create(otp).Error
}

func (r *userOTPRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserOTP{}).Error
}
