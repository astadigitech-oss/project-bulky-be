package controllers

import (
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// BindJSON parses JSON body and validates the struct
func BindJSON(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return err
	}
	return utils.Validator.Struct(req)
}

// localsString reads a string value from Fiber's Locals context
func localsString(c *fiber.Ctx, key string) string {
	val := c.Locals(key)
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// formValueArray returns all values for a multipart/form-data field key
func formValueArray(c *fiber.Ctx, key string) []string {
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil
	}
	return form.Value[key]
}

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
