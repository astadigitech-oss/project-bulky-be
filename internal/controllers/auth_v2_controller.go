package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthV2Controller struct {
	authService  services.AuthV2Service
	adminService services.AdminService
	buyerService services.BuyerService
}

func NewAuthV2Controller(authService services.AuthV2Service, adminService services.AdminService, buyerService services.BuyerService) *AuthV2Controller {
	return &AuthV2Controller{
		authService:  authService,
		adminService: adminService,
		buyerService: buyerService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// POST /api/auth/admin/login
func (c *AuthV2Controller) AdminLogin(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := BindJSON(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
	}

	result, err := c.authService.AdminLogin(ctx.UserContext(), req.Email, req.Password)
	if err != nil {
		// Check error type for proper status code
		if err.Error() == "akun Anda tidak aktif. Silakan hubungi admin" {
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})
}

// POST /api/auth/buyer/login
func (c *AuthV2Controller) BuyerLogin(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := BindJSON(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
	}

	result, err := c.authService.BuyerLogin(ctx.UserContext(), req.Email, req.Password)
	if err != nil {
		// Check error type for proper status code
		if err.Error() == "akun Anda tidak aktif. Silakan hubungi admin" {
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})
}

// POST /api/auth/logout
func (c *AuthV2Controller) Logout(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil",
	})
}

// GET /api/auth/check
func (c *AuthV2Controller) Check(ctx *fiber.Ctx) error {
	// Jika sampai sini, berarti token sudah valid (lolos AuthMiddleware)
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token valid",
	})
}

// GET /api/auth/me
func (c *AuthV2Controller) GetMe(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	if userID == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid atau sudah expired",
		})
	}

	userType := ctx.Locals("user_type")
	if userType == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid atau sudah expired",
		})
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak valid",
		})
	}

	// Handle based on user type
	switch userType.(string) {
	case "ADMIN":
		admin, err := c.authService.GetAdminWithPermissions(ctx.UserContext(), uid)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Admin tidak ditemukan",
			})
		}
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    admin,
		})

	case "BUYER":
		buyer, err := c.authService.GetBuyer(ctx.UserContext(), uid)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Buyer tidak ditemukan",
			})
		}
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    buyer,
		})

	default:
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User type tidak valid",
		})
	}
}

// PUT /api/v1/auth/profile
func (c *AuthV2Controller) UpdateProfile(ctx *fiber.Ctx) error {
	// Get user info from JWT context
	userID := ctx.Locals("user_id")
	if userID == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak ditemukan",
		})
	}

	userType := ctx.Locals("user_type")
	if userType == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User type tidak ditemukan",
		})
	}

	// Route to appropriate handler based on user type
	switch userType.(string) {
	case "ADMIN":
		return c.updateAdminProfile(ctx, userID.(string))
	case "BUYER":
		return c.updateBuyerProfile(ctx, userID.(string))
	default:
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User type tidak valid",
		})
	}
}

func (c *AuthV2Controller) updateAdminProfile(ctx *fiber.Ctx, userID string) error {
	var req dto.AdminUpdateProfileRequest
	if err := BindJSON(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validasi gagal",
			"errors":  err.Error(),
		})
	}

	// Check email unique (exclude current user)
	exists, err := c.adminService.IsEmailExistExcludeID(ctx.UserContext(), req.Email, userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa email",
		})
	}
	if exists {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email sudah digunakan oleh user lain",
		})
	}

	// Update profile
	admin, err := c.adminService.UpdateProfile(ctx.UserContext(), userID, req.Nama, req.Email)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate profile",
		})
	}

	// Simplified response (tanpa permissions & role)
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data": fiber.Map{
			"id":    admin.ID.String(),
			"nama":  admin.Nama,
			"email": admin.Email,
		},
	})
}

func (c *AuthV2Controller) updateBuyerProfile(ctx *fiber.Ctx, userID string) error {
	var req dto.BuyerUpdateProfileRequest
	if err := BindJSON(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validasi gagal",
			"errors":  err.Error(),
		})
	}

	// Check email unique (exclude current user)
	emailExists, err := c.buyerService.IsEmailExistExcludeID(ctx.UserContext(), req.Email, userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa email",
		})
	}
	if emailExists {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email sudah digunakan oleh user lain",
		})
	}

	// Check username unique (exclude current user)
	usernameExists, err := c.buyerService.IsUsernameExistExcludeID(ctx.UserContext(), req.Username, userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa username",
		})
	}
	if usernameExists {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Username sudah digunakan oleh user lain",
		})
	}

	// Validate phone format
	if !utils.IsValidIndonesianPhone(req.Telepon) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format telepon tidak valid",
		})
	}

	// Update profile
	buyer, err := c.buyerService.UpdateProfile(ctx.UserContext(), userID, req.Nama, req.Username, req.Email, req.Telepon)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate profile",
		})
	}

	// Simplified response
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data": fiber.Map{
			"id":       buyer.ID.String(),
			"nama":     buyer.Nama,
			"username": buyer.Username,
			"email":    buyer.Email,
			"telepon":  buyer.Telepon,
		},
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// PUT /api/auth/change-password
func (c *AuthV2Controller) ChangePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	if userID == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak ditemukan",
		})
	}

	userType := ctx.Locals("user_type")
	if userType == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User type tidak ditemukan",
		})
	}

	var req ChangePasswordRequest
	if err := BindJSON(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi password tidak cocok",
		})
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak valid",
		})
	}

	err = c.authService.ChangePassword(ctx.UserContext(), uid, userType.(string), req.CurrentPassword, req.NewPassword)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "password saat ini salah" {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Password saat ini salah",
			})
		}
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Password berhasil diubah",
	})
}

type ActivityLogController struct {
	service services.ActivityLogService
}

func NewActivityLogController(service services.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{service: service}
}

// GET /api/v1/admin/activity-log
func (c *ActivityLogController) GetLogs(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	perPage, _ := strconv.Atoi(ctx.Query("per_page", "20"))
	search := ctx.Query("search")  // ADMIN, BUYER, SYSTEM
	sortBy := ctx.Query("sort_by") // field untuk sorting
	order := ctx.Query("order")    // asc atau desc

	// Validasi sort_by - hanya allow field yang aman
	validSortFields := map[string]bool{
		"created_at": true,
		"user_type":  true,
		"action":     true,
		"modul":      true,
	}
	if sortBy != "" && !validSortFields[sortBy] {
		sortBy = "created_at" // fallback ke default
	}

	// Validasi order - hanya allow asc atau desc
	if order != "asc" && order != "desc" {
		order = "desc" // fallback ke default
	}

	filter := repositories.ActivityLogFilter{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		SortBy:  sortBy,
		Order:   order,
	}

	logs, total, err := c.service.GetLogs(filter)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	meta := models.NewPaginationMeta(page, perPage, total)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    logs,
		"meta":    meta,
	})
}

// GET /api/v1/admin/activity-log/:id
func (c *ActivityLogController) GetLogByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	log, err := c.service.GetLogByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Log tidak ditemukan",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    log,
	})
}

// GET /api/v1/admin/activity-log/entity/:entity_type/:entity_id
func (c *ActivityLogController) GetLogsByEntity(ctx *fiber.Ctx) error {
	entityType := ctx.Params("entity_type")
	entityID, err := uuid.Parse(ctx.Params("entity_id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Entity ID tidak valid",
		})
	}

	logs, err := c.service.GetLogsByEntity(entityType, entityID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    logs,
	})
}

type RoleController struct {
	service services.RoleService
}

func NewRoleController(service services.RoleService) *RoleController {
	return &RoleController{service: service}
}

// GET /api/v1/admin/role
func (c *RoleController) GetAll(ctx *fiber.Ctx) error {
	roles, err := c.service.GetAll()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    roles,
	})
}

// GET /api/v1/admin/role/:id
func (c *RoleController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	role, err := c.service.GetByIDWithPermissions(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Role tidak ditemukan",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    role,
	})
}

type PermissionController struct {
	service services.PermissionService
}

func NewPermissionController(service services.PermissionService) *PermissionController {
	return &PermissionController{service: service}
}

// GET /api/v1/admin/permission
func (c *PermissionController) GetAll(ctx *fiber.Ctx) error {
	permissions, err := c.service.GetByModul()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    permissions,
	})
}
