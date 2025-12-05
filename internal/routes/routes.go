package routes

import (
	"net/http"
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		{
			public.GET("/example", controllers.ExampleController)
		}

		// Protected routes (with authentication middleware)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/protected", controllers.ProtectedController)
		}
	}

	// Endpoint to list all registered routes
	router.GET("/routes", func(c *gin.Context) {
		var endpointList []gin.H
		for _, route := range router.Routes() {
			endpointList = append(endpointList, gin.H{
				"method":  route.Method,
				"path":    route.Path,
				"handler": route.Handler, // This will be the name of the handler function
			})
		}
		c.JSON(http.StatusOK, endpointList)
	})
}
