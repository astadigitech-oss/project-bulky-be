package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req models.LoginRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	ipAddress := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	result, err := c.service.Login(ctx.UserContext(), &req, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "akun Anda telah dinonaktifkan. Hubungi administrator" {
			status = http.StatusForbidden
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Login berhasil", result)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.RefreshToken(ctx.UserContext(), req.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Token berhasil diperbarui", result)
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	c.service.Logout(ctx.UserContext(), req.RefreshToken)
	return utils.SuccessResponse(ctx, "Logout berhasil", nil)
}

func (c *AuthController) Me(ctx *fiber.Ctx) error {
	adminID := localsString(ctx, "admin_id")

	result, err := c.service.GetProfile(ctx.UserContext(), adminID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Profil berhasil diambil", result)
}

func (c *AuthController) UpdateProfile(ctx *fiber.Ctx) error {
	adminID := localsString(ctx, "admin_id")

	var req models.UpdateProfileRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateProfile(ctx.UserContext(), adminID, &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email sudah digunakan oleh admin lain" {
			status = http.StatusConflict
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Profil berhasil diupdate", result)
}

func (c *AuthController) ChangePassword(ctx *fiber.Ctx) error {
	adminID := localsString(ctx, "admin_id")

	var req models.ChangePasswordRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.ChangePassword(ctx.UserContext(), adminID, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Password berhasil diubah", nil)
}
