package response

type UserData struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	User        UserData `json:"user"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

func ToLoginResponse(user UserData, accessToken string) LoginResponse {
	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        user,
	}
}

func ToRegisterResponse(id string) RegisterResponse {
	return RegisterResponse{ID: id}
}
