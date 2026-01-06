// internal/middleware/cors.go

package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Get allowed origins from environment variable
		// Example: ALLOWED_ORIGINS=https://admin.bulky.id,https://panel.bulky.id
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")

		if allowedOriginsEnv != "" {
			// Production: specific origins
			allowedOrigins := strings.Split(allowedOriginsEnv, ",")
			isAllowed := false

			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					isAllowed = true
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if isAllowed {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		} else {
			// Development: allow all
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
