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
	v1 := router.Group("/api/v1")
	{
		// Public Auth Routes (v2)
		authV2 := v1.Group("/auth")
		{
			authV2.POST("/login", authV2Controller.Login)
			authV2.POST("/refresh", authV2Controller.RefreshToken)
		}

		// Protected Auth Routes (v2)
		authV2Protected := v1.Group("/auth")
		authV2Protected.Use(middleware.AuthMiddleware())
		{
			authV2Protected.POST("/logout", authV2Controller.Logout)
			authV2Protected.GET("/me", authV2Controller.GetMe)
			// authV2Protected.PUT("/profile", authV2Controller.UpdateProfile) // TODO
			// authV2Protected.PUT("/change-password", authV2Controller.ChangePassword) // TODO
		}

		// Role Management Routes (Admin Only)
		roleAdmin := v1.Group("/admin/role")
		roleAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("role:manage"))
		{
			roleAdmin.GET("", roleController.GetAll)
			roleAdmin.GET("/:id", roleController.GetByID)
			// roleAdmin.POST("", roleController.Create) // TODO: Add create/update/delete
			// roleAdmin.PUT("/:id", roleController.Update)
			// roleAdmin.DELETE("/:id", roleController.Delete)
		}

		// Permission Routes (Admin Only)
		permissionAdmin := v1.Group("/admin/permission")
		permissionAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("role:manage"))
		{
			permissionAdmin.GET("", permissionController.GetAll)
		}

		// Activity Log Routes (Admin Only)
		activityLogAdmin := v1.Group("/admin/activity-log")
		activityLogAdmin.Use(middleware.AuthMiddleware(), middleware.AdminOnly(), middleware.RequirePermission("activity_log:read"))
		{
			activityLogAdmin.GET("", activityLogController.GetLogs)
			activityLogAdmin.GET("/:id", activityLogController.GetLogByID)
			activityLogAdmin.GET("/entity/:entity_type/:entity_id", activityLogController.GetLogsByEntity)
		}
	}
}
