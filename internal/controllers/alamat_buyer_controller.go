package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AlamatBuyerController struct {
	service services.AlamatBuyerService
}

func NewAlamatBuyerController(service services.AlamatBuyerService) *AlamatBuyerController {
	return &AlamatBuyerController{service: service}
}

func (c *AlamatBuyerController) Create(ctx *gin.Context) {
	var req models.CreateAlamatBuyerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Alamat berhasil ditambahkan", result)
}

func (c *AlamatBuyerController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Detail alamat berhasil diambil", result)
}

func (c *AlamatBuyerController) FindAll(ctx *gin.Context) {
	var params models.AlamatBuyerFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data alamat berhasil diambil", items, *meta)
}

func (c *AlamatBuyerController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateAlamatBuyerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Alamat berhasil diupdate", result)
}

func (c *AlamatBuyerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "alamat tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Alamat berhasil dihapus", nil)
}

func (c *AlamatBuyerController) SetDefault(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.SetDefault(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(ctx, "Alamat berhasil dijadikan default", result)
}
