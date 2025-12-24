package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	result, err := c.service.Login(ctx.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "akun Anda telah dinonaktifkan. Hubungi administrator" {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Login berhasil", result)
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Token berhasil diperbarui", result)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	c.service.Logout(ctx.Request.Context(), req.RefreshToken)
	utils.SuccessResponse(ctx, "Logout berhasil", nil)
}


func (c *AuthController) Me(ctx *gin.Context) {
	adminID := ctx.GetString("admin_id")

	result, err := c.service.GetProfile(ctx.Request.Context(), adminID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Profil berhasil diambil", result)
}

func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	adminID := ctx.GetString("admin_id")

	var req models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateProfile(ctx.Request.Context(), adminID, &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email sudah digunakan oleh admin lain" {
			status = http.StatusConflict
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Profil berhasil diupdate", result)
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	adminID := ctx.GetString("admin_id")

	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.ChangePassword(ctx.Request.Context(), adminID, &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Password berhasil diubah", nil)
}
