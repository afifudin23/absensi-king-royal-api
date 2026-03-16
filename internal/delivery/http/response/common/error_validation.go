package common

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"reflect"
	"strings"
	"unicode"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/go-playground/validator/v10"
)

func CheckTypeError(typeError *json.UnmarshalTypeError, errorMap map[string]string) map[string]string {
	field := normalizeFieldName(typeError.Field)
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
			field := normalizeFieldName(fieldErr.Field())

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
			case "email":
				errorsMap[field] = field + " must be a valid email"
			case "oneof":
				errorsMap[field] = field + " must be one of: " + strings.ReplaceAll(fieldErr.Param(), " ", ", ")
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

func normalizeFieldName(field string) string {
	// jika ada dot (struct nested), ambil field terakhir
	parts := strings.Split(field, ".")
	field = parts[len(parts)-1]

	// ambil tag json dari struct
	// default fallback ke snake_case hanya jika json tag tidak ada
	t := reflect.TypeOf(request.AttendanceRequest{}) // contoh untuk struct tertentu
	if f, ok := t.FieldByName(field); ok {
		tag := f.Tag.Get("json")
		if tag != "" {
			tag = strings.Split(tag, ",")[0] // ambil sebelum koma
			return tag
		}
	}

	// fallback: convert CamelCase → snake_case
	var builder strings.Builder
	for i, r := range field {
		if unicode.IsUpper(r) && i > 0 {
			builder.WriteByte('_')
		}
		builder.WriteRune(unicode.ToLower(r))
	}
	return builder.String()
}
