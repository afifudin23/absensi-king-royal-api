package request

import (
	"strings"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type UserCreateRequest struct {
	FullName string         `json:"full_name" binding:"required,min=3,max=255"`
	Email    string         `json:"email" binding:"required,email,max=255"`
	Password string         `json:"password" binding:"required,min=8,max=100"`
	Role     model.UserRole `json:"role" binding:"required,oneof=admin user"`

	EmployeeCode      *string           `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string           `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string           `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time        `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *model.UserGender `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Address           *string           `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string           `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string           `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string           `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string           `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	BasicSalary       *float64          `json:"basic_salary,omitempty" binding:"omitempty"`
	ProfilePictureID  *string           `json:"profile_picture_id,omitempty" binding:"omitempty,uuid"`
}

func (r *UserCreateRequest) Normalize() {
	r.FullName = strings.TrimSpace(r.FullName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
	normalizeOptionalString(&r.ProfilePictureID, false)
}

type UserUpdateRequest struct {
	FullName *string         `json:"full_name" binding:"omitempty,min=3,max=255"`
	Role     *model.UserRole `json:"role" binding:"omitempty,oneof=admin user"`

	EmployeeCode      *string           `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string           `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string           `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time        `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *model.UserGender `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Address           *string           `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string           `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string           `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string           `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string           `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	BasicSalary       *float64          `json:"basic_salary,omitempty" binding:"omitempty"`
	ProfilePictureID  *string           `json:"profile_picture_id,omitempty" binding:"omitempty,uuid"`
}

func (r *UserUpdateRequest) Normalize() {
	normalizeOptionalString(&r.FullName, false)
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
	normalizeOptionalString(&r.ProfilePictureID, false)
}

type UserUpdateProfileRequest struct {
	FullName *string         `json:"full_name" binding:"omitempty,min=3,max=255"`
	Email    *string         `json:"email" binding:"omitempty,email,max=255"`
	Password *string         `json:"password" binding:"omitempty,min=8,max=100"`
	Role     *model.UserRole `json:"role" binding:"omitempty,oneof=admin user"`

	EmployeeCode      *string           `json:"employee_code,omitempty" binding:"omitempty,max=100"`
	EmploymentStatus  *string           `json:"employment_status,omitempty" binding:"omitempty,max=100"`
	BirthPlace        *string           `json:"birth_place,omitempty" binding:"omitempty,max=255"`
	BirthDate         *time.Time        `json:"birth_date,omitempty" binding:"omitempty"`
	Gender            *model.UserGender `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Address           *string           `json:"address,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string           `json:"phone_number,omitempty" binding:"omitempty,max=50"`
	Position          *string           `json:"position,omitempty" binding:"omitempty,max=100"`
	Department        *string           `json:"department,omitempty" binding:"omitempty,max=100"`
	BankAccountNumber *string           `json:"bank_account_number,omitempty" binding:"omitempty,max=100"`
	ProfilePictureID  *string           `json:"profile_picture_id,omitempty" binding:"omitempty,uuid"`
}

func (r *UserUpdateProfileRequest) Normalize() {
	normalizeOptionalString(&r.FullName, false)
	normalizeOptionalString(&r.Email, true)
	normalizeOptionalString(&r.Password, false)
	normalizeOptionalString(&r.EmployeeCode, false)
	normalizeOptionalString(&r.EmploymentStatus, true)
	normalizeOptionalString(&r.BirthPlace, false)
	normalizeOptionalString(&r.Address, false)
	normalizeOptionalString(&r.PhoneNumber, false)
	normalizeOptionalString(&r.Position, false)
	normalizeOptionalString(&r.Department, false)
	normalizeOptionalString(&r.BankAccountNumber, false)
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
