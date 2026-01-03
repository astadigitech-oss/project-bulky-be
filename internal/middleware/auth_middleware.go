package middleware

import (
	"fmt"
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

		// Set permissions dan role untuk Admin
		if claims.UserType == "ADMIN" {
			c.Set("user_role_id", claims.RoleID)
			c.Set("user_role_kode", claims.RoleKode)
			c.Set("user_permissions", claims.Permissions)
		}

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
			var message string
			if userType == "ADMIN" {
				message = "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin."
			} else if userType == "BUYER" {
				message = "Akses ditolak. Endpoint ini hanya dapat diakses oleh Buyer."
			} else {
				message = "Akses ditolak"
			}
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": message,
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
		// 1. Cek user type - harus ADMIN
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User tidak terautentikasi",
			})
			c.Abort()
			return
		}

		if userType.(string) != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
			c.Abort()
			return
		}

		// 2. Get permissions dari context
		permissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
			c.Abort()
			return
		}

		// 3. Cek apakah punya permission yang dibutuhkan
		perms, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error parsing permissions",
			})
			c.Abort()
			return
		}

		// 4. Loop cek permission
		hasPermission := false
		for _, p := range perms {
			if p == permission {
				hasPermission = true
				break
			}
		}

		// 5. Permission tidak ditemukan
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": fmt.Sprintf("Akses ditolak. Anda tidak memiliki permission: %s", permission),
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
		// 1. Cek user type - harus ADMIN
		userType, exists := c.Get("user_type")
		if !exists || userType.(string) != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
			c.Abort()
			return
		}

		// 2. Get permissions dari context
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
			c.Abort()
			return
		}

		perms := userPermissions.([]string)

		// 3. Cek apakah punya salah satu permission
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

		// 4. Tidak punya permission apapun
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": fmt.Sprintf("Akses ditolak. Anda memerlukan salah satu permission: %v", permissions),
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

// RequireAllPermissions memastikan user punya SEMUA permissions
func RequireAllPermissions(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Cek user type - harus ADMIN
		userType, exists := c.Get("user_type")
		if !exists || userType.(string) != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
			c.Abort()
			return
		}

		// 2. Get permissions dari context
		permissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Permission tidak ditemukan.",
			})
			c.Abort()
			return
		}

		perms := permissions.([]string)
		permMap := make(map[string]bool)
		for _, p := range perms {
			permMap[p] = true
		}

		// 3. Cek semua permission harus ada
		missingPerms := []string{}
		for _, required := range requiredPermissions {
			if !permMap[required] {
				missingPerms = append(missingPerms, required)
			}
		}

		if len(missingPerms) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": fmt.Sprintf("Akses ditolak. Permission yang kurang: %v", missingPerms),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SuperAdminOnly memastikan hanya Super Admin yang bisa akses
func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType.(string) != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Admin.",
			})
			c.Abort()
			return
		}

		roleKode, exists := c.Get("user_role_kode")
		if !exists || roleKode.(string) != "SUPER_ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini hanya dapat diakses oleh Super Admin.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
