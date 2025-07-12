package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-z0-9]{6,}$`, username)
	return matched
}
