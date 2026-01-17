package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WhatsAppHandlerController struct {
	service services.WhatsAppHandlerService
}

func NewWhatsAppHandlerController(service services.WhatsAppHandlerService) *WhatsAppHandlerController {
	return &WhatsAppHandlerController{service: service}
}

func (c *WhatsAppHandlerController) Create(ctx *gin.Context) {
	var req models.CreateWhatsAppHandlerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "WhatsApp handler berhasil dibuat", result)
}

func (c *WhatsAppHandlerController) FindAll(ctx *gin.Context) {
	var params models.WhatsAppHandlerFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data WhatsApp handler berhasil diambil", items, *meta)
}

func (c *WhatsAppHandlerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail WhatsApp handler berhasil diambil", result)
}

func (c *WhatsAppHandlerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateWhatsAppHandlerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler berhasil diupdate", result)
}

func (c *WhatsAppHandlerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "WhatsApp handler tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler berhasil dihapus", nil)
}

func (c *WhatsAppHandlerController) SetActive(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.SetActive(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler berhasil diaktifkan", nil)
}

// Public endpoint
func (c *WhatsAppHandlerController) GetActive(ctx *gin.Context) {
	result, err := c.service.GetActive(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if result == nil {
		utils.SuccessResponse(ctx, "Tidak ada WhatsApp handler aktif", nil)
		return
	}

	utils.SuccessResponse(ctx, "WhatsApp handler aktif berhasil diambil", result)
}
