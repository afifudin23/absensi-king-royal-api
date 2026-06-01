package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserOTP struct {
	ID          string     `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID      string     `gorm:"column:user_id;type:char(36);not null;uniqueIndex"`
	OTPCode     string     `gorm:"column:otp_code;type:varchar(6);not null"`
	ExpiresAt   time.Time  `gorm:"column:expires_at;type:datetime;not null"`
	ResendCount int        `gorm:"column:resend_count;type:int;not null;default:0"`
	FrozenUntil *time.Time `gorm:"column:frozen_until;type:datetime;null"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (u *UserOTP) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

func (UserOTP) TableName() string {
	return "user_otps"
}
