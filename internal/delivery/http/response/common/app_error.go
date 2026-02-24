package common

import (
	"net/http"
)

type ErrorCode string

const (
	BAD_REQUEST              ErrorCode = "BAD_REQUEST"
	VALIDATION_ERROR         ErrorCode = "VALIDATION_ERROR"
	AUTH_REQUIRED            ErrorCode = "AUTH_REQUIRED"
	AUTH_INVALID_CREDENTIALS ErrorCode = "AUTH_INVALID_CREDENTIALS"
	NOT_FOUND                ErrorCode = "NOT_FOUND"
	FORBIDDEN                ErrorCode = "FORBIDDEN"
	INTERNAL_SERVER_ERROR    ErrorCode = "INTERNAL_SERVER_ERROR"
	DATABASE_ERROR           ErrorCode = "DATABASE_ERROR"
	TOO_MANY_REQUESTS        ErrorCode = "TOO_MANY_REQUESTS"
)

type AppError struct {
	StatusCode int
	Code       ErrorCode
	Message    string
	Details    any
}

func (error *AppError) Error() string {
	return error.Message
}

func NewAppError(statusCode int, code ErrorCode, message string, details any) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, AUTH_REQUIRED, message, nil)
}

func ForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, FORBIDDEN, message, nil)
}

func BadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, BAD_REQUEST, message, nil)
}

func ValidationError(details any) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, VALIDATION_ERROR, "Invalid request body", details)
}

func NotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, NOT_FOUND, message, nil)
}

func InternalServerError() *AppError {
	return NewAppError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR, "An unexpected error occurred", nil)
}
