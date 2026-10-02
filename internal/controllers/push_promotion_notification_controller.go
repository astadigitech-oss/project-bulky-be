package controllers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
)

type PushPromotionNotificationController struct {
	service *services.PushPromotionNotificationService
}

func NewPushPromotionNotificationController(service *services.PushPromotionNotificationService) *PushPromotionNotificationController {
	return &PushPromotionNotificationController{service: service}
}

func (c *PushPromotionNotificationController) Create(ctx *fiber.Ctx) error {
	var req models.CreatePushPromotionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Data notifikasi promo tidak valid", parseValidationErrors(err))
	}
	result, err := c.service.Create(ctx.UserContext(), &req, pushPromotionAdminID(ctx))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Notifikasi promo tidak dapat dibuat", err.Error())
	}
	return utils.CreatedResponse(ctx, "Draft notifikasi promo berhasil dibuat", result)
}

func (c *PushPromotionNotificationController) FindAll(ctx *fiber.Ctx) error {
	var params models.PushPromotionListRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter notifikasi promo tidak valid", err.Error())
	}
	result, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil notifikasi promo", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil diambil", result)
}

func (c *PushPromotionNotificationController) FindByID(ctx *fiber.Ctx) error {
	result, err := c.service.FindByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Notifikasi promo tidak ditemukan", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil diambil", result)
}

func (c *PushPromotionNotificationController) Update(ctx *fiber.Ctx) error {
	var req models.UpdatePushPromotionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Data notifikasi promo tidak valid", parseValidationErrors(err))
	}
	result, err := c.service.Update(ctx.UserContext(), ctx.Params("id"), &req, pushPromotionAdminID(ctx))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Notifikasi promo tidak dapat diubah", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil diubah", result)
}

func (c *PushPromotionNotificationController) Schedule(ctx *fiber.Ctx) error {
	var req models.SchedulePushPromotionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Jadwal notifikasi promo tidak valid", parseValidationErrors(err))
	}
	result, err := c.service.Schedule(ctx.UserContext(), ctx.Params("id"), req.ScheduledAt, pushPromotionAdminID(ctx))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Notifikasi promo tidak dapat dijadwalkan", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil dijadwalkan", result)
}

func (c *PushPromotionNotificationController) CancelSchedule(ctx *fiber.Ctx) error {
	result, err := c.service.CancelSchedule(ctx.UserContext(), ctx.Params("id"), pushPromotionAdminID(ctx))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Jadwal notifikasi promo tidak dapat dibatalkan", err.Error())
	}
	return utils.SuccessResponse(ctx, "Jadwal notifikasi promo dibatalkan", result)
}

func (c *PushPromotionNotificationController) SendNow(ctx *fiber.Ctx) error {
	result, err := c.service.SendNow(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadGateway, "Notifikasi promo gagal dikirim", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil dikirim", result)
}

func (c *PushPromotionNotificationController) Delete(ctx *fiber.Ctx) error {
	if err := c.service.Delete(ctx.UserContext(), ctx.Params("id")); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Notifikasi promo tidak dapat dihapus", err.Error())
	}
	return utils.SuccessResponse(ctx, "Notifikasi promo berhasil dihapus", nil)
}

func pushPromotionAdminID(ctx *fiber.Ctx) *uuid.UUID {
	rawID, ok := ctx.Locals("admin_id").(string)
	if !ok || strings.TrimSpace(rawID) == "" {
		return nil
	}
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil
	}
	return &id
}
