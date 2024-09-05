package domain

import (
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
	NotTooOldTag      = "nottooold"
	NotFutureDateTag  = "notfuturedate"
	ValidateImagesTag = "validateImages"
	MaxAgeInYears     = 200
	MaxImagesAllowed  = 5
	MaxImageSize      = 5 * 1024 * 1024
)

var AllowedImagesExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

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

	if err := validator.RegisterValidation(ValidateImagesTag, imageFileValidator); err != nil {
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

func imageFileValidator(fl validator.FieldLevel) bool {
	images := fl.Field().Interface().([]*multipart.FileHeader)

	if len(images) == 0 {
		return true
	}

	if len(images) > MaxImagesAllowed {
		return false
	}

	for _, file := range images {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !AllowedImagesExtensions[ext] {
			return false
		}

		if file.Size > MaxImageSize {
			return false
		}
	}

	return true
}
