package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PPNController struct {
	service services.PPNService
}

func NewPPNController(service services.PPNService) *PPNController {
	return &PPNController{service: service}
}

func (c *PPNController) GetAll(ctx *gin.Context) {
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

	utils.PaginatedSuccessResponse(ctx, "Data PPN berhasil diambil", items, *meta)
}

func (c *PPNController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail PPN berhasil diambil", result)
}

func (c *PPNController) Create(ctx *gin.Context) {
	var req models.CreatePPNRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "PPN berhasil dibuat", result)
}

func (c *PPNController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdatePPNRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "PPN berhasil diupdate", result)
}

func (c *PPNController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if err.Error() == "Tidak dapat menghapus PPN yang sedang aktif" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if err.Error() == "PPN tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "PPN berhasil dihapus", nil)
}

func (c *PPNController) SetActive(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.SetActive(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "PPN tidak ditemukan" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "PPN berhasil diaktifkan", result)
}
