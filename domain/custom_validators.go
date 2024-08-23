package domain

import (
	"github.com/dlclark/regexp2"
	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
	UUIDTag           = "uuid"
)

func SetupCustomValidations(validator *validator.Validate) {
	validator.RegisterValidation(StrongPasswordTag, strongPasswordValidator)
	validator.RegisterValidation(UUIDTag, uuidValidator)
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$&*])[A-Za-z\d!@#$&*]{8,}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}

func uuidValidator(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())

	return err != nil
}
