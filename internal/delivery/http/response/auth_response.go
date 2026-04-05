package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type UserData struct {
	ID        string         `json:"id"`
	FullName  string         `json:"full_name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	User        UserData `json:"user"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

func ToUserData(user model.User) UserData {
	return UserData{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToLoginResponse(user model.User, accessToken string) LoginResponse {
	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        ToUserData(user),
	}
}

func ToRegisterResponse(id string) RegisterResponse {
	return RegisterResponse{ID: id}
}
