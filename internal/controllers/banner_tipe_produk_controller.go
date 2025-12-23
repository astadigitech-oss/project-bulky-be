package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BannerTipeProdukController struct {
	service services.BannerTipeProdukService
}

func NewBannerTipeProdukController(service services.BannerTipeProdukService) *BannerTipeProdukController {
	return &BannerTipeProdukController{service: service}
}

func (c *BannerTipeProdukController) Create(ctx *gin.Context) {
	var req models.CreateBannerTipeProdukRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Banner berhasil dibuat", result)
}

func (c *BannerTipeProdukController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	tipeProdukID := ctx.Query("tipe_produk_id")

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params, tipeProdukID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data banner berhasil diambil", items, *meta)
}


func (c *BannerTipeProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail banner berhasil diambil", result)
}

func (c *BannerTipeProdukController) FindByTipeProdukID(ctx *gin.Context) {
	tipeProdukID := ctx.Param("tipe_produk_id")

	items, err := c.service.FindByTipeProdukID(ctx.Request.Context(), tipeProdukID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diambil", items)
}

func (c *BannerTipeProdukController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateBannerTipeProdukRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
}

func (c *BannerTipeProdukController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Banner berhasil dihapus", nil)
}

func (c *BannerTipeProdukController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status banner berhasil diubah", result)
}

func (c *BannerTipeProdukController) Reorder(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan banner berhasil diubah", nil)
}
