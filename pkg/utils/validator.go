package utils

import (
	"fmt"
	"regexp"
	"strings"

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

// GetValidationErrorMessage converts validation errors to Indonesian messages
func GetValidationErrorMessage(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("Field '%s' wajib diisi", field)
			case "email":
				message = fmt.Sprintf("Field '%s' harus berupa email yang valid", field)
			case "min":
				message = fmt.Sprintf("Field '%s' minimal %s karakter", field, param)
			case "max":
				message = fmt.Sprintf("Field '%s' maksimal %s karakter", field, param)
			case "oneof":
				message = fmt.Sprintf("Field '%s' harus salah satu dari: %s", field, param)
			case "alphanum":
				message = fmt.Sprintf("Field '%s' hanya boleh berisi huruf dan angka", field)
			case "alphanumund":
				message = fmt.Sprintf("Field '%s' hanya boleh berisi huruf, angka, dan underscore", field)
			case "uppercase_snake":
				message = fmt.Sprintf("Field '%s' harus dalam format UPPERCASE_SNAKE_CASE", field)
			case "numeric":
				message = fmt.Sprintf("Field '%s' harus berupa angka", field)
			case "url":
				message = fmt.Sprintf("Field '%s' harus berupa URL yang valid", field)
			case "uuid":
				message = fmt.Sprintf("Field '%s' harus berupa UUID yang valid", field)
			case "gte":
				message = fmt.Sprintf("Field '%s' harus lebih besar atau sama dengan %s", field, param)
			case "lte":
				message = fmt.Sprintf("Field '%s' harus lebih kecil atau sama dengan %s", field, param)
			case "gt":
				message = fmt.Sprintf("Field '%s' harus lebih besar dari %s", field, param)
			case "lt":
				message = fmt.Sprintf("Field '%s' harus lebih kecil dari %s", field, param)
			default:
				message = fmt.Sprintf("Field '%s' tidak valid", field)
			}
			messages = append(messages, message)
		}
		return strings.Join(messages, "; ")
	}
	return err.Error()
}
