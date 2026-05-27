package controllers

import (
	"net/http"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type MerekProdukController struct {
	service services.MerekProdukService
	cfg     *config.Config
}

func NewMerekProdukController(service services.MerekProdukService, cfg *config.Config) *MerekProdukController {
	return &MerekProdukController{
		service: service,
		cfg:     cfg,
	}
}

func (c *MerekProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateMerekProdukRequest
	var logoURL *string

	// Check content type
	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.NamaID = ctx.FormValue("nama_id")
		if namaEN := ctx.FormValue("nama_en"); namaEN != "" {
			req.NamaEN = &namaEN
		}

		// Validate required field
		if req.NamaID == "" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Nama merek (Indonesia) wajib diisi", nil)
		}

		// Handle logo upload or URL
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file logo tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-brands", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
			}
			logoURL = &savedPath
		} else if logoStr := ctx.FormValue("logo"); logoStr != "" {
			logoURL = &logoStr
		}

		// Create with logo
		result, err := c.service.CreateWithLogo(ctx.UserContext(), &req, logoURL)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		}
		return utils.CreatedResponse(ctx, "Merek produk berhasil dibuat", result)
	}

	// Handle application/json (no file upload)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Merek produk berhasil dibuat", result)
}

func (c *MerekProdukController) FindAll(ctx *fiber.Ctx) error {
	var params models.MerekProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data merek produk berhasil diambil", items, *meta)
}

func (c *MerekProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail merek produk berhasil diambil", result)
}

func (c *MerekProdukController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail merek produk berhasil diambil", result)
}

func (c *MerekProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateMerekProdukRequest
	var logoURL *string

	// Check content type
	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if namaID := ctx.FormValue("nama_id"); namaID != "" {
			req.NamaID = &namaID
		}
		if namaEN := ctx.FormValue("nama_en"); namaEN != "" {
			req.NamaEN = &namaEN
		}
		if isActiveStr := ctx.FormValue("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle logo upload or URL
		if file, err := ctx.FormFile("logo"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file logo tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-brands", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan logo: "+err.Error(), nil)
			}
			logoURL = &savedPath
		} else if logoStr := ctx.FormValue("logo"); logoStr != "" {
			logoURL = &logoStr
		}

		// Use UpdateWithLogo for multipart
		result, err := c.service.UpdateWithLogo(ctx.UserContext(), id, &req, logoURL)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.SuccessResponse(ctx, "Merek produk berhasil diupdate", result)
	}

	// Handle application/json (no file upload)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Merek produk berhasil diupdate", result)
}

func (c *MerekProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "merek produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Merek produk berhasil dihapus", nil)
}

func (c *MerekProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status merek berhasil diubah", result)
}

func (c *MerekProdukController) Dropdown(ctx *fiber.Ctx) error {
	merekList, err := c.service.GetAllForDropdown(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data merek", nil)
	}

	return utils.SuccessResponse(ctx, "Data dropdown merek produk berhasil diambil", merekList)
}
