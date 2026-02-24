package model

import "time"

type User struct {
	ID                string     `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	FullName          string     `gorm:"column:full_name"`
	Email             string     `gorm:"column:email"`
	Password          string     `gorm:"column:password"`
	Role              string     `gorm:"column:role"`
	EmployeeCode      *string    `gorm:"column:employee_code"`
	EmploymentStatus  *string    `gorm:"column:employment_status"`
	BirthPlace        *string    `gorm:"column:birth_place"`
	BirthDate         *time.Time `gorm:"column:birth_date"`
	Gender            *string    `gorm:"column:gender"`
	Address           *string    `gorm:"column:address"`
	PhoneNumber       *string    `gorm:"column:phone_number"`
	Position          *string    `gorm:"column:position"`
	Department        *string    `gorm:"column:department"`
	BankAccountNumber *string    `gorm:"column:bank_account_number"`
	ProfilePictureURL *string    `gorm:"column:profile_picture_url"`
	ProfilePictureID  *string    `gorm:"column:profile_picture_public_id"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
