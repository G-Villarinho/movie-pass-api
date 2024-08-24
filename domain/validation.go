package domain

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors map[string]string

var validationMessages = map[string]string{
	"required":        "This field is required",
	"email":           "Invalid email format",
	"min":             "Value is too short",
	"max":             "Value is too long",
	"eqfield":         "Fields do not match",
	"gt":              "The value must be greater than zero",
	"datetime":        "Invalid birth date",
	StrongPasswordTag: "Password must be at least 8 characters long, contain an uppercase letter, a number, and a special character",
}

func ValidateStruct(s any) ValidationErrors {
	validate := validator.New()

	if err := SetupCustomValidations(validate); err != nil {
		return ValidationErrors{"validation_setup": "Failed to set up custom validations"}
	}

	validationErrors := make(ValidationErrors)
	if err := validate.Struct(s); err != nil {
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
