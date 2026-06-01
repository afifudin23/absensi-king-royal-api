package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string



const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID       string   `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	FullName string   `gorm:"column:full_name;type:varchar(255);not null"`
	Email    string   `gorm:"column:email;type:varchar(255);not null"`
	Password string   `gorm:"column:password;type:varchar(255);not null"`
	Role     UserRole `gorm:"column:role;type:enum('admin','user');not null;default:user"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null"`

	Profile *UserProfile `gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
