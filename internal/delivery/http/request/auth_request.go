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
	// Email    string `json:"email" binding:"required,email,max=255"` // TODO: re-enable email validation for production
	Email    string `json:"email" binding:"required,max=255"`
	// Password string `json:"password" binding:"required,min=8,max=72"` // TODO: re-enable min=8 for production
	Password string `json:"password" binding:"required,min=3,max=72"`
}

func (r *AuthLoginRequest) Normalize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

func (r *ForgotPasswordRequest) Normalize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email,max=255"`
	OTP         string `json:"otp" binding:"required,len=6"`
	// NewPassword string `json:"new_password" binding:"required,min=8,max=100"` // TODO: re-enable min=8 for production
	NewPassword string `json:"new_password" binding:"required,min=3,max=100"`
}

func (r *ResetPasswordRequest) Normalize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.OTP = strings.TrimSpace(r.OTP)
}
