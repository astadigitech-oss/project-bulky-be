package utils

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitCustomValidators registers custom validators
func InitCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register alphanumund validator (alphanumeric + underscore)
		v.RegisterValidation("alphanumund", validateAlphaNumUnderscore)

		// Register uppercase_snake validator untuk kode role (UPPER_CASE_SNAKE)
		v.RegisterValidation("uppercase_snake", validateUppercaseSnake)
	}
}

// validateAlphaNumUnderscore validates that a string contains only alphanumeric characters and underscores
func validateAlphaNumUnderscore(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(fl.Field().String())
}

// validateUppercaseSnake validates that a string is UPPERCASE with underscores (LIKE_THIS)
func validateUppercaseSnake(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`).MatchString(fl.Field().String())
}

// IsValidIndonesianPhone validates Indonesian phone number format
// Accept formats: 08xx, +628xx, 628xx
func IsValidIndonesianPhone(phone string) bool {
	pattern := `^(\+62|62|0)8[1-9][0-9]{7,10}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}
