package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service services.AdminService
}

func NewAdminController(service services.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (c *AdminController) Create(ctx *gin.Context) {
	var req models.CreateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email sudah terdaftar" {
			status = http.StatusConflict
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Admin berhasil dibuat", result)
}

func (c *AdminController) FindAll(ctx *gin.Context) {
	var params models.AdminFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data admin berhasil diambil", items, *meta)
}

func (c *AdminController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail admin berhasil diambil", result)
}

func (c *AdminController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "email sudah digunakan oleh admin lain" {
			status = http.StatusConflict
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Admin berhasil diupdate", result)
}

func (c *AdminController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	currentAdminID := ctx.GetString("admin_id")

	if err := c.service.Delete(ctx.Request.Context(), id, currentAdminID); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Admin berhasil dihapus", nil)
}

func (c *AdminController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	currentAdminID := ctx.GetString("admin_id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id, currentAdminID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status admin berhasil diubah", result)
}

func (c *AdminController) ResetPassword(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.ResetPassword(ctx.Request.Context(), id, &req); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "admin tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Password admin berhasil direset", nil)
}
