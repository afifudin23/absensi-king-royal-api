package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyRegistered      = errors.New("Email is already registered")
	ErrInvalidCredentials          = errors.New("Email or password is invalid, please try again")
	ErrUserNotFound                = errors.New("User not found")
	ErrPayrollSettingAlreadyExists = errors.New("Payroll setting already exists")
)

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "1062")
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found")
}
