package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// InternalKeyMiddleware validates requests from internal services via X-Internal-Key header.
func InternalKeyMiddleware(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if key == "" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Internal API key tidak dikonfigurasi",
			})
		}
		if c.Get("X-Internal-Key") != key {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}
		return c.Next()
	}
}
