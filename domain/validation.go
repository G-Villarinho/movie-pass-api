package domain

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validationMessages = map[string]string{
	"required":        "This field is required",
	"email":           "Invalid email format",
	"min":             "Value is too short",
	"max":             "Value is too long",
	"eqfield":         "Fields do not match",
	"gt":              "The value must be greater than zero",
	StrongPasswordTag: "Password must be at least 8 characters long, contain an uppercase letter, a number, and a special character",
	UUIDTag:           "Invalid uuid format",
}

func ValidateStruct(s any) map[string]string {
	validate := validator.New()

	SetupCustomValidations(validate)

	err := validate.Struct(s)
	validationErrors := make(map[string]string)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			validationErrors[fieldName] = getErrorMessage(err)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func getErrorMessage(err validator.FieldError) string {
	if msg, exists := validationMessages[err.Tag()]; exists {
		return msg
	}
	return "Invalid value"
}
