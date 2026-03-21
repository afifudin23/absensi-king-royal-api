package model

import "time"

type UserProfile struct {
	ID                string      `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID            string      `gorm:"column:user_id;type:char(36);not null"`
	EmployeeCode      *string     `gorm:"column:employee_code;type:varchar(100);null"`
	EmploymentStatus  *string     `gorm:"column:employment_status;type:enum('permanent','contract','internship','freelance');null"`
	BirthPlace        *string     `gorm:"column:birth_place;type:varchar(100);null"`
	BirthDate         *time.Time  `gorm:"column:birth_date;type:date;null"`
	Gender            *UserGender `gorm:"column:gender;type:enum('male','female','other');null"`
	Address           *string     `gorm:"column:address;type:text;null"`
	PhoneNumber       *string     `gorm:"column:phone_number;type:varchar(20);null"`
	Position          *string     `gorm:"column:position;type:varchar(100);null"`
	Department        *string     `gorm:"column:department;type:varchar(100);null"`
	BankAccountNumber *string     `gorm:"column:bank_account_number;type:varchar(100);null"`
	BasicSalary       *float64    `gorm:"column:basic_salary;type:decimal(15,2);null"`
	PositionAllowance *float64    `gorm:"column:position_allowance;type:decimal(15,2);null"`
	OtherAllowance    *float64    `gorm:"column:other_allowance;type:decimal(15,2);null"`

	ProfilePictureID  *string `gorm:"column:profile_picture_id;type:char(36);null"`
	ProfilePictureURL *string `gorm:"column:profile_picture_url;type:text;null"`

	JoinedAt  *time.Time `gorm:"column:joined_at;type:timestamp;null;default:null"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp;not null"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
