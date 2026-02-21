package common

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CheckTypeError(typeError *json.UnmarshalTypeError, errorMap map[string]string) map[string]string {
	field := strings.ToLower(typeError.Field)
	if field == "" {
		field = "request"
	}

	switch typeError.Type.Kind() {
	case reflect.String:
		errorMap[field] = field + " must be a string"
	case reflect.Int, reflect.Int64:
		errorMap[field] = field + " must be an integer"
	case reflect.Float32, reflect.Float64:
		errorMap[field] = field + " must be a number"
	case reflect.Bool:
		errorMap[field] = field + " must be a boolean"
	default:
		errorMap[field] = field + " has invalid type"
	}

	return errorMap
}

func ErrorValidation(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Body kosong
	if errors.Is(err, io.EOF) {
		errorsMap["request"] = "Request body is required"
		return errorsMap
	}

	// Salah tipe data (json unmarshal)
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return CheckTypeError(typeErr, errorsMap)
	}

	// Validation errors (binding tag)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := strings.ToLower(fieldErr.Field())

			switch fieldErr.Tag() {
			case "uuid":
				errorsMap[field] = field + " must be a valid UUID"
			case "required":
				errorsMap[field] = field + " is required"
			case "min":
				errorsMap[field] = field + " must be at least " + fieldErr.Param() + " characters"
			case "max":
				errorsMap[field] = field + " must be at most " + fieldErr.Param() + " characters"
			case "eqfield":
				errorsMap[field] = field + " does not match " + strings.ToLower(fieldErr.Param())
			case "datetime":
				errorsMap[field] = field + " must be a valid datetime format, must be 'YYYY-MM-DDTHH:MM:SS±HH:MM'"
			default:
				errorsMap[field] = field + " is invalid"
			}
		}
		return errorsMap
	}

	// FALLBACK
	log.Println(err)
	errorsMap["request"] = "Invalid request payload"
	return errorsMap
}
