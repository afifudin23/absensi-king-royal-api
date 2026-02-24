package request

import "strings"

type AuthRegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (r *AuthRegisterRequest) Normalize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.FullName = strings.TrimSpace(r.FullName)
}

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (r *AuthLoginRequest) Normalize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}
