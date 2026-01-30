package controllers

import (
	"net/http"
	"strconv"

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

// GetAllGrouped - Admin endpoint with grouped response
func (c *MetodePembayaranController) GetAllGrouped(ctx *gin.Context) {
	result, err := c.service.GetAllGrouped(ctx.Request.Context(), true) // isAdmin = true
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data", nil)
		return
	}

	utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", result)
}

// GetAllGroupedPublic - Public endpoint with grouped response (active only)
func (c *MetodePembayaranController) GetAllGroupedPublic(ctx *gin.Context) {
	result, err := c.service.GetAllGrouped(ctx.Request.Context(), false) // isAdmin = false
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data", nil)
		return
	}

	utils.SuccessResponse(ctx, "Data metode pembayaran berhasil diambil", result)
}

// ToggleMethodStatus - Toggle status of individual payment method
func (c *MetodePembayaranController) ToggleMethodStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleMethodStatus(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "Metode pembayaran tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status metode pembayaran berhasil diubah", result)
}

// ToggleGroupStatus - Toggle status of payment group by urutan
func (c *MetodePembayaranController) ToggleGroupStatus(ctx *gin.Context) {
	urutanStr := ctx.Param("urutan")
	urutan, err := strconv.Atoi(urutanStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Urutan tidak valid", nil)
		return
	}

	result, err := c.service.ToggleGroupStatus(ctx.Request.Context(), urutan)
	if err != nil {
		if err.Error() == "Group tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status group berhasil diubah", result)
}

// Legacy methods - kept for backward compatibility if needed
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

	var req models.ReorderByDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
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
		req.Direction,
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
