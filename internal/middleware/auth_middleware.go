package middleware

import (
	"net/http"
	"strings"

	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format authorization tidak valid",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
			})
			c.Abort()
			return
		}

		// Set user context (new format)
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_permissions", claims.Permissions)

		// Set legacy context for backward compatibility
		if claims.AdminID != "" {
			c.Set("admin_id", claims.AdminID)
			c.Set("admin_email", claims.Email)
		} else if claims.UserID != "" && claims.UserType == "ADMIN" {
			c.Set("admin_id", claims.UserID)
			c.Set("admin_email", claims.Email)
		}

		c.Next()
	}
}

// RequireUserType checks if user is of specific type (ADMIN or BUYER)
func RequireUserType(userType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User type tidak ditemukan",
			})
			c.Abort()
			return
		}

		if currentUserType != userType {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission checks if user has specific permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Permission tidak ditemukan",
			})
			c.Abort()
			return
		}

		perms, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Format permission tidak valid",
			})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission := false
		for _, p := range perms {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda tidak memiliki akses untuk melakukan aksi ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has at least one of the specified permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Permission tidak ditemukan",
			})
			c.Abort()
			return
		}

		perms, ok := userPermissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Format permission tidak valid",
			})
			c.Abort()
			return
		}

		// Check if user has at least one of the required permissions
		hasPermission := false
		for _, requiredPerm := range permissions {
			for _, userPerm := range perms {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda tidak memiliki akses untuk melakukan aksi ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly is a convenience middleware that requires ADMIN user type
func AdminOnly() gin.HandlerFunc {
	return RequireUserType("ADMIN")
}

// BuyerOnly is a convenience middleware that requires BUYER user type
func BuyerOnly() gin.HandlerFunc {
	return RequireUserType("BUYER")
}
