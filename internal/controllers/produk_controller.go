package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProdukController struct {
	service        services.ProdukService
	gambarService  services.ProdukGambarService
	dokumenService services.ProdukDokumenService
}

func NewProdukController(
	service services.ProdukService,
	gambarService services.ProdukGambarService,
	dokumenService services.ProdukDokumenService,
) *ProdukController {
	return &ProdukController{
		service:        service,
		gambarService:  gambarService,
		dokumenService: dokumenService,
	}
}

func (c *ProdukController) Create(ctx *gin.Context) {
	var req models.CreateProdukRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Produk berhasil dibuat", result)
}

func (c *ProdukController) FindAll(ctx *gin.Context) {
	var params models.ProdukFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data produk berhasil diambil", items, *meta)
}

func (c *ProdukController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail produk berhasil diambil", result)
}

func (c *ProdukController) FindBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	result, err := c.service.FindBySlug(ctx.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail produk berhasil diambil", result)
}


func (c *ProdukController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateProdukRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Produk berhasil diupdate", result)
}

func (c *ProdukController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Produk berhasil dihapus", nil)
}

func (c *ProdukController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status produk berhasil diubah", result)
}

func (c *ProdukController) UpdateStock(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateStock(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Stok produk berhasil diupdate", result)
}

// ========================================
// Produk Gambar Handlers
// ========================================

func (c *ProdukController) AddGambar(ctx *gin.Context) {
	produkID := ctx.Param("produk_id")

	var req models.CreateProdukGambarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.gambarService.Create(ctx.Request.Context(), produkID, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Gambar berhasil ditambahkan", result)
}

func (c *ProdukController) UpdateGambar(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateProdukGambarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.gambarService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Gambar berhasil diupdate", result)
}

func (c *ProdukController) DeleteGambar(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.gambarService.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Gambar berhasil dihapus", nil)
}

func (c *ProdukController) ReorderGambar(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.gambarService.Reorder(ctx.Request.Context(), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan gambar berhasil diubah", nil)
}

// ========================================
// Produk Dokumen Handlers
// ========================================

func (c *ProdukController) AddDokumen(ctx *gin.Context) {
	produkID := ctx.Param("produk_id")

	var req models.CreateProdukDokumenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.dokumenService.Create(ctx.Request.Context(), produkID, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Dokumen berhasil ditambahkan", result)
}

func (c *ProdukController) DeleteDokumen(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.dokumenService.Delete(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Dokumen berhasil dihapus", nil)
}
