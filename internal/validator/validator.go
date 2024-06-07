package validator

import (
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("password", ValidatePassword)
	validate.RegisterValidation("username", ValidateUsername)
	validate.RegisterValidation("expiration_date", ValidateExpirationDate)

	return validate
}

func ValidateExpirationDate(field validator.FieldLevel) bool {
	date, ok := field.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

func ValidatePassword(field validator.FieldLevel) bool {
	password := field.Field().String()
	hasUpper := false
	hasLower := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasSpecial
}

func ValidateUsername(field validator.FieldLevel) bool {
	username := field.Field().String()
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}
	return true
}
