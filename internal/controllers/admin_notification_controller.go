package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
)

type AdminNotificationController struct {
	service *services.AdminNotificationService
}

func NewAdminNotificationController(service *services.AdminNotificationService) *AdminNotificationController {
	return &AdminNotificationController{service: service}
}

func (c *AdminNotificationController) RegisterDevice(ctx *fiber.Ctx) error {
	adminID, err := adminNotificationID(ctx)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}
	var req struct {
		FID string `json:"fid"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Payload perangkat push tidak valid", err.Error())
	}
	if err := c.service.RegisterDevice(ctx.UserContext(), adminID, req.FID); err != nil {
		if errors.Is(err, services.ErrInvalidAdminPushFID) {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "FID tidak valid", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendaftarkan perangkat push", err.Error())
	}
	return utils.SuccessResponse(ctx, "Perangkat push berhasil didaftarkan", nil)
}

func (c *AdminNotificationController) UnregisterDevice(ctx *fiber.Ctx) error {
	adminID, err := adminNotificationID(ctx)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}
	var req struct {
		FID string `json:"fid"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Payload perangkat push tidak valid", err.Error())
	}
	if err := c.service.UnregisterDevice(ctx.UserContext(), adminID, req.FID); err != nil {
		if errors.Is(err, services.ErrInvalidAdminPushFID) {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "FID tidak valid", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal melepas perangkat push", err.Error())
	}
	return utils.SuccessResponse(ctx, "Perangkat push berhasil dilepas", nil)
}

func (c *AdminNotificationController) List(ctx *fiber.Ctx) error {
	adminID, err := adminNotificationID(ctx)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}
	result, err := c.service.List(ctx.UserContext(), adminID)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil notifikasi", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi berhasil diambil", result)
}

func (c *AdminNotificationController) MarkRead(ctx *fiber.Ctx) error {
	adminID, err := adminNotificationID(ctx)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak valid", err.Error())
	}
	notificationID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID notifikasi tidak valid", err.Error())
	}
	if err := c.service.MarkRead(ctx.UserContext(), adminID, notificationID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Notifikasi tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui notifikasi", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi ditandai sudah dibaca", nil)
}

func adminNotificationID(ctx *fiber.Ctx) (uuid.UUID, error) {
	adminID, ok := ctx.Locals("admin_id").(string)
	if !ok || strings.TrimSpace(adminID) == "" {
		return uuid.Nil, errors.New("admin_id tidak tersedia di sesi")
	}
	return uuid.Parse(adminID)
}
