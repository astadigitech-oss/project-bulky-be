package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")

	if allowedOriginsEnv != "" {
		// Production: specific origins
		origins := strings.Join(
			func() []string {
				var trimmed []string
				for _, o := range strings.Split(allowedOriginsEnv, ",") {
					trimmed = append(trimmed, strings.TrimSpace(o))
				}
				return trimmed
			}(),
			",",
		)
		return cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Content-Length,Accept,Authorization,X-Requested-With",
			AllowCredentials: true,
			MaxAge:           86400,
		})
	}

	// Development: allow all
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Content-Length,Accept,Authorization,X-Requested-With",
		MaxAge:       86400,
	})
}
