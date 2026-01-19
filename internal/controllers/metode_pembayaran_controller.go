package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MetodePembayaranController struct {
	service        services.MetodePembayaranService
	reorderService *services.ReorderService
}

func NewMetodePembayaranController(service services.MetodePembayaranService, reorderService *services.ReorderService) *MetodePembayaranController {
	return &MetodePembayaranController{
		service:        service,
		reorderService: reorderService,
	}
}

func (c *MetodePembayaranController) GetAll(ctx *gin.Context) {
	// Parse optional filters
	var groupID *string
	if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
		groupID = &groupIDStr
	}

	var isActive *bool
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActiveBool := isActiveStr == "true"
		isActive = &isActiveBool
	}

	items, err := c.service.GetAll(ctx.Request.Context(), groupID, isActive)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", items)
}

func (c *MetodePembayaranController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateMetodePembayaranRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if err.Error() == "Kode sudah digunakan" || err.Error() == "Group tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Metode pembayaran berhasil diupdate", result)
}

func (c *MetodePembayaranController) ReorderByDirection(ctx *gin.Context) {
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
		"metode_pembayaran",
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
