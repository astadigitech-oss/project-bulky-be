package utils

import (
	"net/http"

	"project-bulky-be/internal/models"

	"github.com/gin-gonic/gin"
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

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errors []models.FieldError) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, meta models.PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func PaginatedSuccessResponseWithSummary(c *gin.Context, message string, data interface{}, meta models.PaginationMeta, summary interface{}) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
		Summary: summary,
	})
}

// SimpleErrorResponse for simple error messages (string)
func SimpleErrorResponse(c *gin.Context, statusCode int, message string, errorDetail string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error":   errorDetail,
	})
}

// SimpleSuccessResponse for simple success responses with custom status code
func SimpleSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}
