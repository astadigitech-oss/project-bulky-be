package controllers

import (
	"project-bulky-be/internal/models"

	"github.com/go-playground/validator/v10"
)

func parseValidationErrors(err error) []models.FieldError {
	var errors []models.FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, models.FieldError{
				Field:   e.Field(),
				Message: getValidationMessage(e),
			})
		}
	} else {
		errors = append(errors, models.FieldError{
			Field:   "unknown",
			Message: err.Error(),
		})
	}

	return errors
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " wajib diisi"
	case "min":
		return e.Field() + " minimal " + e.Param() + " karakter"
	case "max":
		return e.Field() + " maksimal " + e.Param() + " karakter"
	case "uuid":
		return e.Field() + " harus berupa UUID yang valid"
	case "oneof":
		return e.Field() + " harus salah satu dari: " + e.Param()
	default:
		return e.Field() + " tidak valid"
	}
}
