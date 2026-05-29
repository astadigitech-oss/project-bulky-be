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

	"github.com/gofiber/fiber/v2"
)

type ProdukController struct {
	service        services.ProdukService
	gambarService  services.ProdukGambarService
	dokumenService services.ProdukDokumenService
	activityLog    services.ActivityLogService
}

func NewProdukController(
	service services.ProdukService,
	gambarService services.ProdukGambarService,
	dokumenService services.ProdukDokumenService,
	activityLog services.ActivityLogService,
) *ProdukController {
	return &ProdukController{
		service:        service,
		gambarService:  gambarService,
		dokumenService: dokumenService,
		activityLog:    activityLog,
	}
}

func (c *ProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateProdukRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Format multipart tidak valid", nil)
	}

	// Get gambar files (required, min 1, max 10)
	gambarFiles := form.File["gambar[]"]
	if len(gambarFiles) == 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Minimal 1 gambar produk wajib diupload", nil)
	}
	if len(gambarFiles) > 10 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 10 gambar produk", nil)
	}

	// Validate gambar files
	for i, file := range gambarFiles {
		if err := validateImageFile(file); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("gambar[%d]: %s", i, err.Error()), nil)
		}
	}

	// Get dokumen files (optional, max 5)
	dokumenFiles := form.File["dokumen[]"]
	if len(dokumenFiles) > 5 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 5 dokumen", nil)
	}

	// Validate dokumen files
	for i, file := range dokumenFiles {
		if err := validateDocumentFile(file); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("dokumen[%d]: %s", i, err.Error()), nil)
		}
	}

	// Get dokumen names (parallel array)
	dokumenNama := formValueArray(ctx, "dokumen_nama[]")

	// Handle merek_id array field
	if req.MerekIDs == nil || *req.MerekIDs == "" {
		merekIDArray := formValueArray(ctx, "merek_id")
		if len(merekIDArray) > 0 {
			merekIDStr := strings.Join(merekIDArray, ",")
			req.MerekIDs = &merekIDStr
		}
	}

	// Handle is_active from form - default false (draft)
	isActive := false
	isActiveStr := strings.ToLower(strings.TrimSpace(ctx.FormValue("is_active")))
	if isActiveStr == "true" || isActiveStr == "1" {
		isActive = true
	}

	result, err := c.service.CreateWithFiles(ctx.UserContext(), &req, isActive, gambarFiles, dokumenFiles, dokumenNama)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "produk", "Produk berhasil dibuat")
	return utils.CreatedResponse(ctx, "Produk berhasil dibuat", result)
}

func (c *ProdukController) FindAll(ctx *fiber.Ctx) error {
	var params models.ProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data produk berhasil diambil", items, *meta)
}

func (c *ProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail produk berhasil diambil", result)
}

func (c *ProdukController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail produk berhasil diambil", result)
}

func (c *ProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateProdukRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Handle merek_id array field
	if req.MerekIDs == nil || *req.MerekIDs == "" {
		merekIDArray := formValueArray(ctx, "merek_id")
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
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 5 dokumen", nil)
		}
		for i, file := range dokumenFiles {
			if err := validateDocumentFile(file); err != nil {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("dokumen[%d]: %s", i, err.Error()), nil)
			}
		}
		dokumenNama = formValueArray(ctx, "dokumen_nama[]")
	}

	result, err := c.service.UpdateWithFiles(ctx.UserContext(), id, &req, dokumenFiles, dokumenNama)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Produk berhasil diupdate")
	return utils.SuccessResponse(ctx, "Produk berhasil diupdate", result)
}

func (c *ProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "produk", "Produk berhasil dihapus")
	return utils.SuccessResponse(ctx, "Produk berhasil dihapus", nil)
}

func (c *ProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "produk", "Status produk berhasil diubah")
	return utils.SuccessResponse(ctx, "Status produk berhasil diubah", result)
}

func (c *ProdukController) UpdateStock(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateStockRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.UpdateStock(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Stok produk berhasil diupdate")
	return utils.SuccessResponse(ctx, "Stok produk berhasil diupdate", result)
}

// ========================================
// Produk Gambar Handlers
// ========================================

func (c *ProdukController) AddGambar(ctx *fiber.Ctx) error {
	produkID := ctx.Params("id")

	form, err := ctx.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal memparse form", nil)
	}

	// Collect files: prioritize gambar[] (multiple), fallback to gambar (single)
	files := form.File["gambar[]"]
	if len(files) == 0 {
		if f, ok := form.File["gambar"]; ok {
			files = f
		}
	}

	if len(files) == 0 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "File gambar wajib diupload (field: gambar atau gambar[])", nil)
	}

	if len(files) > 10 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Maksimal 10 gambar per upload", nil)
	}

	// Validate ALL files first before uploading any
	for i, file := range files {
		if err := validateImageFile(file); err != nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("gambar[%d] (%s): %s", i+1, file.Filename, err.Error()), nil)
		}
	}

	results, err := c.gambarService.CreateMultipleWithFiles(ctx.UserContext(), produkID, files)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Gambar produk berhasil ditambahkan")
	return utils.CreatedResponse(ctx, fmt.Sprintf("%d gambar berhasil ditambahkan", len(results)), results)
}

func (c *ProdukController) DeleteGambar(ctx *fiber.Ctx) error {
	produkID := ctx.Params("id")
	gambarID := ctx.Params("gambar_id")

	if err := c.gambarService.Delete(ctx.UserContext(), produkID, gambarID); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Gambar produk berhasil dihapus")
	return utils.SuccessResponse(ctx, "Gambar berhasil dihapus", nil)
}

func (c *ProdukController) ReorderGambar(ctx *fiber.Ctx) error {
	produkID := ctx.Params("id")
	gambarID := ctx.Params("gambar_id")

	var req models.ReorderGambarRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.gambarService.Reorder(ctx.UserContext(), produkID, gambarID, req.Direction)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan gambar berhasil diubah", result)
}

// ========================================
// Produk Dokumen Handlers
// ========================================

func (c *ProdukController) AddDokumen(ctx *fiber.Ctx) error {
	produkID := ctx.Params("id")

	// Get file from form
	file, err := ctx.FormFile("dokumen")
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "File dokumen wajib diupload", nil)
	}

	// Validate file
	if err := validateDocumentFile(file); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	var req models.CreateProdukDokumenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.dokumenService.CreateWithFile(ctx.UserContext(), produkID, file, req.NamaDokumen)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Dokumen produk berhasil ditambahkan")
	return utils.CreatedResponse(ctx, "Dokumen berhasil ditambahkan", result)
}

func (c *ProdukController) DeleteDokumen(ctx *fiber.Ctx) error {
	produkID := ctx.Params("id")
	dokumenID := ctx.Params("dokumen_id")

	if err := c.dokumenService.Delete(ctx.UserContext(), produkID, dokumenID); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "produk", "Dokumen produk berhasil dihapus")
	return utils.SuccessResponse(ctx, "Dokumen berhasil dihapus", nil)
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
