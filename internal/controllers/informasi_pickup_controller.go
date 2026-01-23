package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type InformasiPickupController struct {
	service services.InformasiPickupService
}

func NewInformasiPickupController(service services.InformasiPickupService) *InformasiPickupController {
	return &InformasiPickupController{service: service}
}

// Admin endpoints
func (c *InformasiPickupController) Get(ctx *gin.Context) {
	result, err := c.service.Get(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Informasi pickup berhasil diambil", result)
}

func (c *InformasiPickupController) Update(ctx *gin.Context) {
	var req models.UpdateInformasiPickupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Informasi pickup berhasil diupdate", result)
}

func (c *InformasiPickupController) UpdateJadwal(ctx *gin.Context) {
	var req models.UpdateJadwalGudangRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateJadwal(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Jadwal gudang berhasil diupdate", result)
}

// Public endpoint
func (c *InformasiPickupController) GetPublic(ctx *gin.Context) {
	result, err := c.service.GetPublic(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Informasi pickup berhasil diambil", result)
}
