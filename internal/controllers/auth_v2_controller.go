package controllers

import (
	"net/http"
	"strconv"

	"project-bulky-be/internal/repositories"
	"project-bulky-be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthV2Controller struct {
	authService services.AuthV2Service
}

func NewAuthV2Controller(authService services.AuthV2Service) *AuthV2Controller {
	return &AuthV2Controller{authService: authService}
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=admin buyer ADMIN BUYER"`
	DeviceInfo string `json:"device_info"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// POST /api/v1/auth/login
func (c *AuthV2Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
		return
	}

	// Normalize type to uppercase
	userType := req.Type
	if userType == "admin" {
		userType = "ADMIN"
	} else if userType == "buyer" {
		userType = "BUYER"
	}

	ipAddress := ctx.ClientIP()
	result, err := c.authService.Login(ctx, req.Email, req.Password, userType, req.DeviceInfo, ipAddress)
	if err != nil {
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

// POST /api/v1/auth/refresh
func (c *AuthV2Controller) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
		})
		return
	}

	result, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token berhasil diperbarui",
		"data":    result,
	})
}

// POST /api/v1/auth/logout
func (c *AuthV2Controller) Logout(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Try to get from header as fallback
		refreshToken := ctx.GetHeader("X-Refresh-Token")
		if refreshToken == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Refresh token diperlukan",
			})
			return
		}
		req.RefreshToken = refreshToken
	}

	err := c.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Logout gagal: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// GET /api/v1/auth/me
func (c *AuthV2Controller) GetMe(ctx *gin.Context) {
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

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID tidak valid",
		})
		return
	}

	user, err := c.authService.GetCurrentUser(ctx, uid, userType.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

type UpdateProfileRequest struct {
	Nama  string `json:"nama" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// PUT /api/v1/auth/profile
func (c *AuthV2Controller) UpdateProfile(ctx *gin.Context) {
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

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid: " + err.Error(),
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

	user, err := c.authService.UpdateProfile(ctx, uid, userType.(string), req.Nama)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diperbarui",
		"data":    user,
	})
}

// PUT /api/v1/auth/change-password
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
			"message": "Password baru dan konfirmasi password tidak sama",
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
