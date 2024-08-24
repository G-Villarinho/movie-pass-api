package domain

import (
	"time"

	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
	NotTooOldTag      = "nottooold"
	NotFutureDateTag  = "notfuturedate"
	MaxAgeInYears     = 200
)

func SetupCustomValidations(validator *validator.Validate) error {
	if err := validator.RegisterValidation(StrongPasswordTag, strongPasswordValidator); err != nil {
		return err
	}

	if err := validator.RegisterValidation(NotTooOldTag, notTooOldValidator); err != nil {
		return err
	}

	if err := validator.RegisterValidation(NotFutureDateTag, notFutureDateValidator); err != nil {
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

func notTooOldValidator(fl validator.FieldLevel) bool {
	birthDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	age := time.Now().Year() - birthDate.Year()
	return age <= MaxAgeInYears
}

func notFutureDateValidator(fl validator.FieldLevel) bool {
	birthDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	if birthDate.After(time.Now()) {
		return false
	}

	return true
}
