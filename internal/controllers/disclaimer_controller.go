package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DisclaimerController struct {
	service services.DisclaimerService
}

func NewDisclaimerController(service services.DisclaimerService) *DisclaimerController {
	return &DisclaimerController{service: service}
}

func (c *DisclaimerController) Create(ctx *gin.Context) {
	var req models.CreateDisclaimerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Disclaimer berhasil dibuat", result)
}

func (c *DisclaimerController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data disclaimer berhasil diambil", items, *meta)
}

func (c *DisclaimerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail disclaimer berhasil diambil", result)
}

func (c *DisclaimerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateDisclaimerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Disclaimer berhasil diupdate", result)
}

func (c *DisclaimerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "disclaimer tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Disclaimer berhasil dihapus", nil)
}

func (c *DisclaimerController) SetActive(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.SetActive(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Disclaimer berhasil diaktifkan", result)
}

// Public endpoint
func (c *DisclaimerController) GetActive(ctx *gin.Context) {
	lang := ctx.DefaultQuery("lang", "id")

	result, err := c.service.GetActive(ctx.Request.Context(), lang)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if result == nil {
		msg := "Tidak ada disclaimer aktif"
		if lang == "en" {
			msg = "No active disclaimer"
		}
		utils.SuccessResponse(ctx, msg, nil)
		return
	}

	msg := "Disclaimer aktif berhasil diambil"
	if lang == "en" {
		msg = "Active disclaimer retrieved successfully"
	}
	utils.SuccessResponse(ctx, msg, result)
}
