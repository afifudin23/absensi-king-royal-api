package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type UserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`

	EmployeeCode      *string    `json:"employee_code"`
	EmploymentStatus  *string    `json:"employment_status"`
	BirthPlace        *string    `json:"birth_place"`
	BirthDate         *time.Time `json:"birth_date"`
	Gender            *string    `json:"gender"`
	Address           *string    `json:"address"`
	PhoneNumber       *string    `json:"phone_number"`
	Position          *string    `json:"position"`
	Department        *string    `json:"department"`
	BankAccountNumber *string    `json:"bank_account_number"`
	ProfilePictureURL *string    `json:"profile_picture_url"`
	ProfilePictureID  *string    `json:"profile_picture_public_id"`
	DeletedAt         *time.Time `json:"deleted_at"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
}

func ToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:                user.ID,
		FullName:          user.FullName,
		Email:             user.Email,
		Role:              user.Role,
		EmployeeCode:      user.EmployeeCode,
		EmploymentStatus:  user.EmploymentStatus,
		BirthPlace:        user.BirthPlace,
		BirthDate:         user.BirthDate,
		Gender:            user.Gender,
		Address:           user.Address,
		PhoneNumber:       user.PhoneNumber,
		Position:          user.Position,
		Department:        user.Department,
		BankAccountNumber: user.BankAccountNumber,
		ProfilePictureURL: user.ProfilePictureURL,
		ProfilePictureID:  user.ProfilePictureID,
		DeletedAt:         user.DeletedAt,
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         user.UpdatedAt.Format(time.RFC3339),
	}
}

func ToUserListResponse(users []model.User) []UserResponse {
	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, ToUserResponse(user))
	}
	return response
}
