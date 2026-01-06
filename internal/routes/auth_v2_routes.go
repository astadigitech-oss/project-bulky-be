package routes

import (
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthV2Routes sets up authentication and authorization routes
func SetupAuthV2Routes(
	router *gin.Engine,
	authV2Controller *controllers.AuthV2Controller,
	roleController *controllers.RoleController,
	permissionController *controllers.PermissionController,
	activityLogController *controllers.ActivityLogController,
) {
	api := router.Group("/api")
	{
		// ============================================================
		// PANEL AUTH ROUTES (Admin Panel)
		// ============================================================
		panelAuth := api.Group("/panel/auth")
		{
			// Public - Admin Login
			panelAuth.POST("/login", authV2Controller.AdminLogin)
		}

		// Protected Panel Auth Routes
		panelAuthProtected := api.Group("/panel/auth")
		panelAuthProtected.Use(middleware.AuthMiddleware())
		panelAuthProtected.Use(middleware.AdminOnly())
		{
			panelAuthProtected.GET("/check", authV2Controller.Check)
			panelAuthProtected.GET("/me", authV2Controller.GetMe)
			panelAuthProtected.POST("/logout", authV2Controller.Logout)
			panelAuthProtected.PUT("/profile", authV2Controller.UpdateProfile)
			panelAuthProtected.PUT("/change-password", authV2Controller.ChangePassword)
		}

		// ============================================================
		// BUYER AUTH ROUTES (Future - Buyer App)
		// ============================================================
		// Uncomment when implementing Buyer App
		// buyerAuth := api.Group("/buyer/auth")
		// {
		// 	buyerAuth.POST("/login", authV2Controller.BuyerLogin)
		// }

		// Role Management Routes (Admin Only)
		roleAdmin := api.Group("/panel/role")
		roleAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("role:manage"))
		{
			roleAdmin.GET("", roleController.GetAll)
			roleAdmin.GET("/:id", roleController.GetByID)
			// roleAdmin.POST("", roleController.Create) // TODO: Add create/update/delete
			// roleAdmin.PUT("/:id", roleController.Update)
			// roleAdmin.DELETE("/:id", roleController.Delete)
		}

		// Permission Routes (Admin Only)
		permissionAdmin := api.Group("/panel/permission")
		permissionAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("role:manage"))
		{
			permissionAdmin.GET("", permissionController.GetAll)
		}

		// Activity Log Routes (Admin Only)
		activityLogAdmin := api.Group("/panel/activity-log")
		activityLogAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("activity_log:read"))
		{
			activityLogAdmin.GET("", activityLogController.GetLogs)
			activityLogAdmin.GET("/:id", activityLogController.GetLogByID)
			activityLogAdmin.GET("/entity/:entity_type/:entity_id", activityLogController.GetLogsByEntity)
		}
	}
}
