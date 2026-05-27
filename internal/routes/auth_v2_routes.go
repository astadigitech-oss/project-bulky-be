package routes

import (
	"project-bulky-be/internal/controllers"
	"project-bulky-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthV2Routes sets up authentication and authorization routes
func SetupAuthV2Routes(
	router *fiber.App,
	authV2Controller *controllers.AuthV2Controller,
	roleController *controllers.RoleController,
	permissionController *controllers.PermissionController,
	activityLogController *controllers.ActivityLogController,
) {
	api := router.Group("/api")

	// ============================================================
	// PANEL AUTH ROUTES (Admin Panel)
	// ============================================================
	panelAuth := api.Group("/panel/auth")
	// Public - Admin Login
	panelAuth.Post("/login", authV2Controller.AdminLogin)

	// Protected Panel Auth Routes
	panelAuthProtected := api.Group("/panel/auth",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
	)
	panelAuthProtected.Get("/check", authV2Controller.Check)
	panelAuthProtected.Get("/me", authV2Controller.GetMe)
	panelAuthProtected.Post("/logout", authV2Controller.Logout)
	panelAuthProtected.Put("/profile", authV2Controller.UpdateProfile)
	panelAuthProtected.Put("/change-password", authV2Controller.ChangePassword)

	// Role Management Routes (Admin Only)
	roleAdmin := api.Group("/panel/role",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		middleware.RequirePermission("role:manage"),
	)
	roleAdmin.Get("", roleController.GetAll)
	roleAdmin.Get("/:id", roleController.GetByID)

	// Permission Routes (Admin Only)
	permissionAdmin := api.Group("/panel/permission",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		middleware.RequirePermission("role:manage"),
	)
	permissionAdmin.Get("", permissionController.GetAll)

	// Activity Log Routes (Admin Only)
	activityLogAdmin := api.Group("/panel/activity-log",
		middleware.AuthMiddleware(),
		middleware.AdminOnly(),
		middleware.RequirePermission("activity_log:read"),
	)
	activityLogAdmin.Get("", activityLogController.GetLogs)
	activityLogAdmin.Get("/:id", activityLogController.GetLogByID)
	activityLogAdmin.Get("/entity/:entity_type/:entity_id", activityLogController.GetLogsByEntity)
}
