package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MetodePembayaranGroupController struct {
	service        services.MetodePembayaranGroupService
	reorderService *services.ReorderService
}

func NewMetodePembayaranGroupController(service services.MetodePembayaranGroupService, reorderService *services.ReorderService) *MetodePembayaranGroupController {
	return &MetodePembayaranGroupController{
		service:        service,
		reorderService: reorderService,
	}
}

func (c *MetodePembayaranGroupController) GetAll(ctx *gin.Context) {
	items, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data group metode pembayaran berhasil diambil", items)
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

func (c *MetodePembayaranGroupController) ReorderByDirection(ctx *gin.Context) {
	id := ctx.Param("id")
	direction := ctx.Query("direction")

	if direction == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter 'direction' wajib diisi", nil)
		return
	}

	idUUID, err := utils.ParseUUID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	result, err := c.reorderService.Reorder(
		ctx.Request.Context(),
		"metode_pembayaran_group",
		idUUID,
		direction,
		"",
		nil,
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan berhasil diubah", gin.H{
		"item": gin.H{
			"id":     result.ItemID,
			"urutan": result.ItemUrutan,
		},
		"swapped_with": gin.H{
			"id":     result.SwappedID,
			"urutan": result.SwappedUrutan,
		},
	})
}
