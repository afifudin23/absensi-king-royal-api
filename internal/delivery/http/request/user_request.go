package request

import (
	"strings"
	"time"
)

type UserCreateRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Role     string `json:"role" binding:"required,oneof=admin user"`

	EmployeeCode      *string    `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string    `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string    `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *string    `json:"gender,omitempty" binding:"omitempty,max=20"`
	Address           *string    `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string    `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string    `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string    `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty" binding:"omitempty,max=500"`
	ProfilePictureID  *string    `json:"profile_picture_public_id,omitempty" binding:"omitempty,max=255"`
}

func (r *UserCreateRequest) Normalize() {
	r.FullName = strings.TrimSpace(r.FullName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Role = strings.ToLower(strings.TrimSpace(r.Role))
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Gender, true)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
	normalizeOptionalString(&r.ProfilePictureURL, false)
	normalizeOptionalString(&r.ProfilePictureID, false)
}

type UserUpdateRequest struct {
	FullName *string `json:"full_name" binding:"omitempty,min=3,max=255"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`

	EmployeeCode      *string    `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string    `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string    `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *string    `json:"gender,omitempty" binding:"omitempty,max=20"`
	Address           *string    `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string    `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string    `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string    `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty" binding:"omitempty,max=500"`
	ProfilePictureID  *string    `json:"profile_picture_public_id,omitempty" binding:"omitempty,max=255"`
}

func (r *UserUpdateRequest) Normalize() {
	normalizeOptionalString(&r.FullName, false)
	normalizeOptionalString(&r.Role, true)
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Gender, true)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
	normalizeOptionalString(&r.ProfilePictureURL, false)
	normalizeOptionalString(&r.ProfilePictureID, false)
}

type UserUpdateProfileRequest struct {
	FullName *string `json:"full_name" binding:"omitempty,min=3,max=255"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	Password *string `json:"password" binding:"omitempty,min=8,max=100"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`

	EmployeeCode      *string    `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string    `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string    `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *string    `json:"gender,omitempty" binding:"omitempty,max=20"`
	Address           *string    `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string    `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string    `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string    `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string    `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty" binding:"omitempty,max=500"`
	ProfilePictureID  *string    `json:"profile_picture_public_id,omitempty" binding:"omitempty,max=255"`
}

func (r *UserUpdateProfileRequest) Normalize() {
	normalizeOptionalString(&r.FullName, false)
	normalizeOptionalString(&r.Email, true)
	normalizeOptionalString(&r.Password, false)
	normalizeOptionalString(&r.Role, true)
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Gender, true)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
	normalizeOptionalString(&r.ProfilePictureURL, false)
	normalizeOptionalString(&r.ProfilePictureID, false)
}

func normalizeOptionalString(field **string, toLower bool) {
	if *field == nil {
		return
	}

	value := strings.TrimSpace(**field)
	if value == "" {
		*field = nil
		return
	}

	if toLower {
		value = strings.ToLower(value)
	}
	**field = value
}
