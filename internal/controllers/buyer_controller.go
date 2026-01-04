package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BuyerController struct {
	service services.BuyerService
}

func NewBuyerController(service services.BuyerService) *BuyerController {
	return &BuyerController{service: service}
}

func (c *BuyerController) FindAll(ctx *gin.Context) {
	var params models.BuyerFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data buyer berhasil diambil", items, *meta)
}

func (c *BuyerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Detail buyer berhasil diambil", result)
}

func (c *BuyerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "buyer tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Buyer berhasil dihapus", nil)
}

func (c *BuyerController) ResetPassword(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.ResetBuyerPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.ResetPassword(ctx.Request.Context(), id, &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Password buyer berhasil direset", nil)
}

func (c *BuyerController) GetStatistik(ctx *gin.Context) {
	result, err := c.service.GetStatistik(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Statistik buyer berhasil diambil", result)
}

func (c *BuyerController) GetChart(ctx *gin.Context) {
	var params models.ChartParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	result, err := c.service.GetChart(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Data chart berhasil diambil", result)
}
