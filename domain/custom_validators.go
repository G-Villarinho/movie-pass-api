package domain

import (
	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
)

func SetupCustomValidations(validator *validator.Validate) error {
	if err := validator.RegisterValidation(StrongPasswordTag, strongPasswordValidator); err != nil {
		return err
	}

	return nil
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$&*])[A-Za-z\d!@#$&*]{8,}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}
