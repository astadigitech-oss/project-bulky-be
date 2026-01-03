package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
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
func (c *AuthV2Controller) AdminLogin(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
		return
	}

	result, err := c.authService.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		// Check error type for proper status code
		if err.Error() == "akun Anda tidak aktif. Silakan hubungi admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})
}

// POST /api/auth/buyer/login
func (c *AuthV2Controller) BuyerLogin(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
		return
	}

	result, err := c.authService.BuyerLogin(ctx, req.Email, req.Password)
	if err != nil {
		// Check error type for proper status code
		if err.Error() == "akun Anda tidak aktif. Silakan hubungi admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})
}

// POST /api/auth/logout
func (c *AuthV2Controller) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// GET /api/auth/me
func (c *AuthV2Controller) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Token tidak valid atau sudah expired",
		})
		return
	}

	userType, exists := ctx.Get("user_type")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Token tidak valid atau sudah expired",
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID tidak valid",
		})
		return
	}

	// Handle based on user type
	switch userType.(string) {
	case "ADMIN":
		admin, err := c.authService.GetAdminWithPermissions(ctx, uid)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Admin tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    admin,
		})

	case "BUYER":
		buyer, err := c.authService.GetBuyer(ctx, uid)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Buyer tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    buyer,
		})

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User type tidak valid",
		})
	}
}

// PUT /api/v1/auth/profile
func (c *AuthV2Controller) UpdateProfile(ctx *gin.Context) {
	// Get user info from JWT context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID tidak ditemukan",
		})
		return
	}

	userType, exists := ctx.Get("user_type")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User type tidak ditemukan",
		})
		return
	}

	// Route to appropriate handler based on user type
	switch userType.(string) {
	case "ADMIN":
		c.updateAdminProfile(ctx, userID.(string))
	case "BUYER":
		c.updateBuyerProfile(ctx, userID.(string))
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User type tidak valid",
		})
	}
}

func (c *AuthV2Controller) updateAdminProfile(ctx *gin.Context, userID string) {
	var req dto.AdminUpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"errors":  err.Error(),
		})
		return
	}

	// Check email unique (exclude current user)
	exists, err := c.adminService.IsEmailExistExcludeID(ctx, req.Email, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memeriksa email",
		})
		return
	}
	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email sudah digunakan oleh user lain",
		})
		return
	}

	// Update profile
	admin, err := c.adminService.UpdateProfile(ctx, userID, req.Nama, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate profile",
		})
		return
	}

	// Simplified response (tanpa permissions & role)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data": gin.H{
			"id":    admin.ID.String(),
			"nama":  admin.Nama,
			"email": admin.Email,
		},
	})
}

func (c *AuthV2Controller) updateBuyerProfile(ctx *gin.Context, userID string) {
	var req dto.BuyerUpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"errors":  err.Error(),
		})
		return
	}

	// Check email unique (exclude current user)
	emailExists, err := c.buyerService.IsEmailExistExcludeID(ctx, req.Email, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memeriksa email",
		})
		return
	}
	if emailExists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email sudah digunakan oleh user lain",
		})
		return
	}

	// Check username unique (exclude current user)
	usernameExists, err := c.buyerService.IsUsernameExistExcludeID(ctx, req.Username, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memeriksa username",
		})
		return
	}
	if usernameExists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username sudah digunakan oleh user lain",
		})
		return
	}

	// Validate phone format
	if !utils.IsValidIndonesianPhone(req.Telepon) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format telepon tidak valid",
		})
		return
	}

	// Update profile
	buyer, err := c.buyerService.UpdateProfile(ctx, userID, req.Nama, req.Username, req.Email, req.Telepon)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate profile",
		})
		return
	}

	// Simplified response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data": gin.H{
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
func (c *AuthV2Controller) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID tidak ditemukan",
		})
		return
	}

	userType, exists := ctx.Get("user_type")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User type tidak ditemukan",
		})
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Konfirmasi password tidak cocok",
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID tidak valid",
		})
		return
	}

	err = c.authService.ChangePassword(ctx, uid, userType.(string), req.CurrentPassword, req.NewPassword)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "password saat ini salah" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Password saat ini salah",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
func (c *ActivityLogController) GetLogs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("halaman", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_halaman", "20"))
	userType := ctx.Query("user_type")
	action := ctx.Query("action")
	modul := ctx.Query("modul")
	tanggalDari := ctx.Query("tanggal_dari")
	tanggalSampai := ctx.Query("tanggal_sampai")

	var userID *uuid.UUID
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}

	filter := repositories.ActivityLogFilter{
		Page:          page,
		PerPage:       perPage,
		UserType:      userType,
		UserID:        userID,
		Action:        action,
		Modul:         modul,
		TanggalDari:   tanggalDari,
		TanggalSampai: tanggalSampai,
	}

	logs, total, err := c.service.GetLogs(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	totalPages := (int(total) + perPage - 1) / perPage

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
		"meta": gin.H{
			"halaman":       page,
			"per_halaman":   perPage,
			"total":         total,
			"total_halaman": totalPages,
		},
	})
}

// GET /api/v1/admin/activity-log/:id
func (c *ActivityLogController) GetLogByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	log, err := c.service.GetLogByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Log tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    log,
	})
}

// GET /api/v1/admin/activity-log/entity/:entity_type/:entity_id
func (c *ActivityLogController) GetLogsByEntity(ctx *gin.Context) {
	entityType := ctx.Param("entity_type")
	entityID, err := uuid.Parse(ctx.Param("entity_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Entity ID tidak valid",
		})
		return
	}

	logs, err := c.service.GetLogsByEntity(entityType, entityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
func (c *RoleController) GetAll(ctx *gin.Context) {
	roles, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// GET /api/v1/admin/role/:id
func (c *RoleController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	role, err := c.service.GetByIDWithPermissions(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Role tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
func (c *PermissionController) GetAll(ctx *gin.Context) {
	permissions, err := c.service.GetByModul()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    permissions,
	})
}
