package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	service     services.AdminService
	activityLog services.ActivityLogService
}

func NewAdminController(service services.AdminService, activityLog services.ActivityLogService) *AdminController {
	return &AdminController{service: service, activityLog: activityLog}
}

func (c *AdminController) Create(ctx *fiber.Ctx) error {
	var req models.CreateAdminRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email sudah terdaftar" {
			status = http.StatusConflict
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "admin", "Admin berhasil dibuat")
	return utils.CreatedResponse(ctx, "Admin berhasil dibuat", result)
}

func (c *AdminController) FindAll(ctx *fiber.Ctx) error {
	var params models.AdminFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data admin berhasil diambil", items, *meta)
}

func (c *AdminController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail admin berhasil diambil", result)
}

func (c *AdminController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateAdminRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "email sudah digunakan oleh admin lain" {
			status = http.StatusConflict
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "admin", "Admin berhasil diupdate")
	return utils.SuccessResponse(ctx, "Admin berhasil diupdate", result)
}

func (c *AdminController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	currentAdminID := localsString(ctx, "admin_id")

	if err := c.service.Delete(ctx.UserContext(), id, currentAdminID); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "admin", "Admin berhasil dihapus")
	return utils.SuccessResponse(ctx, "Admin berhasil dihapus", nil)
}

func (c *AdminController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	currentAdminID := localsString(ctx, "admin_id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id, currentAdminID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "admin", "Status admin berhasil diubah")
	return utils.SuccessResponse(ctx, "Status admin berhasil diubah", result)
}

func (c *AdminController) ResetPassword(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.ResetPasswordRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.ResetPassword(ctx.UserContext(), id, &req); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "admin", "Password admin berhasil direset")
	return utils.SuccessResponse(ctx, "Password admin berhasil direset", nil)
}
