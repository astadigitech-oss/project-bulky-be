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
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(ctx, http.StatusOK, "Data buyer berhasil diambil", items, meta)
}

func (c *BuyerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail buyer berhasil diambil", result)
}

func (c *BuyerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateBuyerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Buyer berhasil diupdate", result)
}

func (c *BuyerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Buyer berhasil dihapus", nil)
}

func (c *BuyerController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Status buyer berhasil diubah", result)
}

func (c *BuyerController) ResetPassword(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.ResetBuyerPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := c.service.ResetPassword(ctx.Request.Context(), id, &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password buyer berhasil direset", nil)
}

func (c *BuyerController) GetStatistik(ctx *gin.Context) {
	result, err := c.service.GetStatistik(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Statistik buyer berhasil diambil", result)
}
