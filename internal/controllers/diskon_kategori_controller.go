package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DiskonKategoriController struct {
	service services.DiskonKategoriService
}

func NewDiskonKategoriController(service services.DiskonKategoriService) *DiskonKategoriController {
	return &DiskonKategoriController{service: service}
}

func (c *DiskonKategoriController) Create(ctx *gin.Context) {
	var req models.CreateDiskonKategoriRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Diskon kategori berhasil dibuat", result)
}

func (c *DiskonKategoriController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	kategoriID := ctx.Query("kategori_id")
	berlakuHariIni := ctx.Query("berlaku_hari_ini") == "true"

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params, kategoriID, berlakuHariIni)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data diskon kategori berhasil diambil", items, *meta)
}


func (c *DiskonKategoriController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail diskon kategori berhasil diambil", result)
}

func (c *DiskonKategoriController) FindActiveByKategoriID(ctx *gin.Context) {
	kategoriID := ctx.Param("kategori_id")

	result, err := c.service.FindActiveByKategoriID(ctx.Request.Context(), kategoriID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if result == nil {
		utils.SuccessResponse(ctx, "Tidak ada diskon aktif untuk kategori ini", nil)
		return
	}

	utils.SuccessResponse(ctx, "Diskon kategori ditemukan", result)
}

func (c *DiskonKategoriController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateDiskonKategoriRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Diskon kategori berhasil diupdate", result)
}

func (c *DiskonKategoriController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "diskon kategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Diskon kategori berhasil dihapus", nil)
}

func (c *DiskonKategoriController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status diskon kategori berhasil diubah", result)
}
