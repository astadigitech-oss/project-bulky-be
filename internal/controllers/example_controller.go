package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExampleController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Example controller response",
		"data":    "Hello from Gin!",
	})
}

func ProtectedController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a protected route",
		"data":    "You are authenticated!",
	})
}
