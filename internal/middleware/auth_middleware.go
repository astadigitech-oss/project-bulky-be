package middleware

import (
	"fmt"
	"strings"

	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak ditemukan",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format authorization tidak valid",
			})
		}

		token := parts[1]

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
			})
		}

		// Set user context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_type", claims.UserType)
		c.Locals("user_email", claims.Email)

		// Set permissions dan role untuk Admin
		if claims.UserType == "ADMIN" {
			c.Locals("user_role_id", claims.RoleID)
			c.Locals("user_role_kode", claims.RoleKode)
			c.Locals("user_permissions", claims.Permissions)
		}

		// Set legacy context for backward compatibility
		if claims.AdminID != "" {
			c.Locals("admin_id", claims.AdminID)
			c.Locals("admin_email", claims.Email)
		} else if claims.UserID != "" && claims.UserType == "ADMIN" {
			c.Locals("admin_id", claims.UserID)
			c.Locals("admin_email", claims.Email)
		}

		return c.Next()
	}
}

// RequireUserType checks if user is of specific type (ADMIN or BUYER)
func RequireUserType(userType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUserType := c.Locals("user_type")
		if currentUserType == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User type tidak ditemukan",
			})
		}

		if currentUserType != userType {
			var message string
			if userType == "ADMIN" {
				message = "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin."
			} else if userType == "BUYER" {
				message = "Akses ditolak. Endpoint ini hanya dapat diakses oleh Buyer."
			} else {
				message = "Akses ditolak"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": message,
			})
		}

		return c.Next()
	}
}

// RequirePermission checks if user has specific permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Cek user type - harus ADMIN
		userType := c.Locals("user_type")
		if userType == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User tidak terautentikasi",
			})
		}

		if userType.(string) != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
		}

		// 2. Get permissions dari context
		permissions := c.Locals("user_permissions")
		if permissions == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
		}

		// 3. Cek apakah punya permission yang dibutuhkan
		perms, ok := permissions.([]string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error parsing permissions",
			})
		}

		// 4. Loop cek permission
		for _, p := range perms {
			if p == permission {
				return c.Next()
			}
		}

		// 5. Permission tidak ditemukan
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Akses ditolak. Anda tidak memiliki permission: %s", permission),
		})
	}
}

// RequireAnyPermission checks if user has at least one of the specified permissions
func RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Cek user type - harus ADMIN
		userType := c.Locals("user_type")
		if userType == nil || userType.(string) != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
		}

		// 2. Get permissions dari context
		userPermissions := c.Locals("user_permissions")
		if userPermissions == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
		}

		perms := userPermissions.([]string)

		// 3. Cek apakah punya salah satu permission
		for _, requiredPerm := range permissions {
			for _, userPerm := range perms {
				if userPerm == requiredPerm {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Akses ditolak. Anda memerlukan salah satu permission: %v", permissions),
		})
	}
}

// AdminOnly is a convenience middleware that requires ADMIN user type
func AdminOnly() fiber.Handler {
	return RequireUserType("ADMIN")
}

// BuyerOnly is a convenience middleware that requires BUYER user type
func BuyerOnly() fiber.Handler {
	return RequireUserType("BUYER")
}

// RequireAllPermissions memastikan user punya SEMUA permissions
func RequireAllPermissions(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Cek user type - harus ADMIN
		userType := c.Locals("user_type")
		if userType == nil || userType.(string) != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
		}

		// 2. Get permissions dari context
		permissions := c.Locals("user_permissions")
		if permissions == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
		}

		perms := permissions.([]string)
		permMap := make(map[string]bool)
		for _, p := range perms {
			permMap[p] = true
		}

		// 3. Cek semua permission harus ada
		var missingPerms []string
		for _, required := range requiredPermissions {
			if !permMap[required] {
				missingPerms = append(missingPerms, required)
			}
		}

		if len(missingPerms) > 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Akses ditolak. Permission yang kurang: %v", missingPerms),
			})
		}

		return c.Next()
	}
}

// SuperAdminOnly memastikan hanya Super Admin yang bisa akses
func SuperAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType == nil || userType.(string) != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
		}

		roleKode := c.Locals("user_role_kode")
		if roleKode == nil || roleKode.(string) != "SUPER_ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Super Admin.",
			})
		}

		return c.Next()
	}
}
