package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MetodePembayaranGroupController struct {
	service services.MetodePembayaranGroupService
}

func NewMetodePembayaranGroupController(service services.MetodePembayaranGroupService) *MetodePembayaranGroupController {
	return &MetodePembayaranGroupController{service: service}
}

func (c *MetodePembayaranGroupController) GetAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	// Set default values
	params.SetDefaults()

	items, meta, err := c.service.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data group metode pembayaran berhasil diambil", items, *meta)
}

func (c *MetodePembayaranGroupController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail group metode pembayaran berhasil diambil", result)
}

func (c *MetodePembayaranGroupController) Create(ctx *gin.Context) {
	var req models.CreateMetodePembayaranGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "Nama group sudah digunakan" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Group metode pembayaran berhasil dibuat", result)
}

func (c *MetodePembayaranGroupController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateMetodePembayaranGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "Group metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "Nama group sudah digunakan" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Group metode pembayaran berhasil diupdate", result)
}

func (c *MetodePembayaranGroupController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "Group metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "Tidak dapat menghapus group yang masih memiliki metode pembayaran aktif" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Group metode pembayaran berhasil dihapus", nil)
}

func (c *MetodePembayaranGroupController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "Group metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status group metode pembayaran berhasil diubah", result)
}
