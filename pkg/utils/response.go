package utils

import (
	"project-bulky-be/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Response struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    interface{}         `json:"data,omitempty"`
	Errors  []models.FieldError `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    interface{}           `json:"data"`
	Meta    models.PaginationMeta `json:"meta"`
	Summary interface{}           `json:"summary,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors []models.FieldError) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta models.PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func PaginatedSuccessResponseWithSummary(c *fiber.Ctx, message string, data interface{}, meta models.PaginationMeta, summary interface{}) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
		Summary: summary,
	})
}

// SimpleErrorResponse for simple error messages (string)
func SimpleErrorResponse(c *fiber.Ctx, statusCode int, message string, errorDetail string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   errorDetail,
	})
}

// SimpleSuccessResponse for simple success responses with custom status code
func SimpleSuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ParseUUID parses string to UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
