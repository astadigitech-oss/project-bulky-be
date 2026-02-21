package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

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
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Format multipart tidak valid", nil)
		return
	}

	// Get gambar files (required, min 1, max 10)
	gambarFiles := form.File["gambar[]"]
	if len(gambarFiles) == 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Minimal 1 gambar produk wajib diupload", nil)
		return
	}
	if len(gambarFiles) > 10 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 10 gambar produk", nil)
		return
	}

	// Validate gambar files
	for i, file := range gambarFiles {
		if err := validateImageFile(file); err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("gambar[%d]: %s", i, err.Error()), nil)
			return
		}
	}

	// Get dokumen files (optional, max 5)
	dokumenFiles := form.File["dokumen[]"]
	if len(dokumenFiles) > 5 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 5 dokumen", nil)
		return
	}

	// Validate dokumen files
	for i, file := range dokumenFiles {
		if err := validateDocumentFile(file); err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("dokumen[%d]: %s", i, err.Error()), nil)
			return
		}
	}

	// Get dokumen names (parallel array)
	dokumenNama := ctx.PostFormArray("dokumen_nama[]")

	// Handle merek_id array field
	if req.MerekIDs == nil || *req.MerekIDs == "" {
		merekIDArray := ctx.PostFormArray("merek_id")
		if len(merekIDArray) > 0 {
			merekIDStr := strings.Join(merekIDArray, ",")
			req.MerekIDs = &merekIDStr
		}
	}

	// Handle is_active from form - default false (draft)
	isActive := false
	isActiveStr := strings.ToLower(strings.TrimSpace(ctx.PostForm("is_active")))
	if isActiveStr == "true" || isActiveStr == "1" {
		isActive = true
	}

	result, err := c.service.CreateWithFiles(ctx.Request.Context(), &req, isActive, gambarFiles, dokumenFiles, dokumenNama)
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
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Handle merek_id array field
	if req.MerekIDs == nil || *req.MerekIDs == "" {
		merekIDArray := ctx.PostFormArray("merek_id")
		if len(merekIDArray) > 0 {
			merekIDStr := strings.Join(merekIDArray, ",")
			req.MerekIDs = &merekIDStr
		}
	}

	// Get dokumen files (optional - if provided, replaces all existing dokumen)
	var dokumenFiles []*multipart.FileHeader
	var dokumenNama []string
	if form, err := ctx.MultipartForm(); err == nil {
		dokumenFiles = form.File["dokumen[]"]
		if len(dokumenFiles) > 5 {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 5 dokumen", nil)
			return
		}
		for i, file := range dokumenFiles {
			if err := validateDocumentFile(file); err != nil {
				utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("dokumen[%d]: %s", i, err.Error()), nil)
				return
			}
		}
		dokumenNama = ctx.PostFormArray("dokumen_nama[]")
	}

	result, err := c.service.UpdateWithFiles(ctx.Request.Context(), id, &req, dokumenFiles, dokumenNama)
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
	produkID := ctx.Param("id")

	form, err := ctx.MultipartForm()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal memparse form", nil)
		return
	}

	// Collect files: prioritize gambar[] (multiple), fallback to gambar (single)
	files := form.File["gambar[]"]
	if len(files) == 0 {
		if f, ok := form.File["gambar"]; ok {
			files = f
		}
	}

	if len(files) == 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar wajib diupload (field: gambar atau gambar[])", nil)
		return
	}

	if len(files) > 10 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 10 gambar per upload", nil)
		return
	}

	// Validate ALL files first before uploading any
	for i, file := range files {
		if err := validateImageFile(file); err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("gambar[%d] (%s): %s", i+1, file.Filename, err.Error()), nil)
			return
		}
	}

	results, err := c.gambarService.CreateMultipleWithFiles(ctx.Request.Context(), produkID, files)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, fmt.Sprintf("%d gambar berhasil ditambahkan", len(results)), results)
}

func (c *ProdukController) DeleteGambar(ctx *gin.Context) {
	produkID := ctx.Param("id")
	gambarID := ctx.Param("gambar_id")

	if err := c.gambarService.Delete(ctx.Request.Context(), produkID, gambarID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Gambar berhasil dihapus", nil)
}

func (c *ProdukController) ReorderGambar(ctx *gin.Context) {
	produkID := ctx.Param("id")
	gambarID := ctx.Param("gambar_id")

	var req models.ReorderGambarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.gambarService.Reorder(ctx.Request.Context(), produkID, gambarID, req.Direction)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan gambar berhasil diubah", result)
}

// ========================================
// Produk Dokumen Handlers
// ========================================

func (c *ProdukController) AddDokumen(ctx *gin.Context) {
	produkID := ctx.Param("id")

	// Get file from form
	file, err := ctx.FormFile("dokumen")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "File dokumen wajib diupload", nil)
		return
	}

	// Validate file
	if err := validateDocumentFile(file); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var req models.CreateProdukDokumenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.dokumenService.CreateWithFile(ctx.Request.Context(), produkID, file, req.NamaDokumen)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Dokumen berhasil ditambahkan", result)
}

func (c *ProdukController) DeleteDokumen(ctx *gin.Context) {
	produkID := ctx.Param("id")
	dokumenID := ctx.Param("dokumen_id")

	if err := c.dokumenService.Delete(ctx.Request.Context(), produkID, dokumenID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Dokumen berhasil dihapus", nil)
}

// ========================================
// File Validation Helpers
// ========================================

func validateImageFile(file *multipart.FileHeader) error {
	// Check size (max 5MB)
	if file.Size > 5*1024*1024 {
		return errors.New("ukuran file melebihi 5MB")
	}

	// Check type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return errors.New("format file tidak didukung. Gunakan JPG, PNG, atau WebP")
	}

	return nil
}

func validateDocumentFile(file *multipart.FileHeader) error {
	// Check size (max 10MB)
	if file.Size > 10*1024*1024 {
		return errors.New("ukuran file melebihi 10MB")
	}

	// Check type
	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return errors.New("format file harus PDF")
	}

	return nil
}
