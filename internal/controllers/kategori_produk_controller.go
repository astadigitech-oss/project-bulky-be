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

type KategoriProdukController struct {
	service services.KategoriProdukService
	cfg     *config.Config
}

func NewKategoriProdukController(service services.KategoriProdukService, cfg *config.Config) *KategoriProdukController {
	return &KategoriProdukController{
		service: service,
		cfg:     cfg,
	}
}

func (c *KategoriProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateKategoriProdukRequest
	var iconURL *string
	var gambarKondisiURL *string

	// Check content type
	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.NamaID = ctx.FormValue("nama_id")
		if namaEN := ctx.FormValue("nama_en"); namaEN != "" {
			req.NamaEN = &namaEN
		}
		if deskripsi := ctx.FormValue("deskripsi"); deskripsi != "" {
			req.Deskripsi = &deskripsi
		}
		// if memilikiKondisi := ctx.FormValue("memiliki_kondisi_tambahan"); memilikiKondisi != "" {
		// 	hasKondisi := memilikiKondisi == "true"
		// 	req.MemilikiKondisiTambahan = hasKondisi
		// }
		if tipeKondisi := ctx.FormValue("tipe_kondisi_tambahan"); tipeKondisi != "" {
			req.TipeKondisiTambahan = &tipeKondisi
		}
		if teksKondisi := ctx.FormValue("teks_kondisi"); teksKondisi != "" {
			req.TeksKondisi = &teksKondisi
		}

		// Validate required field
		if req.NamaID == "" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Nama kategori wajib diisi", nil)
		}

		// Handle icon upload or URL
		if file, err := ctx.FormFile("icon"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file icon tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan icon: "+err.Error(), nil)
			}
			iconURL = &savedPath
		} else if iconStr := ctx.FormValue("icon"); iconStr != "" {
			iconURL = &iconStr
		}

		// Handle gambar kondisi upload or URL
		if file, err := ctx.FormFile("gambar_kondisi"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar kondisi tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories/kondisi", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan gambar kondisi: "+err.Error(), nil)
			}
			gambarKondisiURL = &savedPath
		} else if gambarKondisiStr := ctx.FormValue("gambar_kondisi"); gambarKondisiStr != "" {
			gambarKondisiURL = &gambarKondisiStr
		}

		// Create with icon
		result, err := c.service.CreateWithIcon(ctx.UserContext(), &req, iconURL, gambarKondisiURL)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
		}
		return utils.CreatedResponse(ctx, "Kategori produk berhasil dibuat", result)
	}

	// Handle application/json (no file upload)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
	}

	return utils.CreatedResponse(ctx, "Kategori produk berhasil dibuat", result)
}

func (c *KategoriProdukController) FindAll(ctx *fiber.Ctx) error {
	var params models.KategoriProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data kategori produk berhasil diambil", items, *meta)
}

func (c *KategoriProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kategori produk berhasil diambil", result)
}

func (c *KategoriProdukController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	result, err := c.service.FindBySlug(ctx.UserContext(), slug)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail kategori produk berhasil diambil", result)
}

func (c *KategoriProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateKategoriProdukRequest
	var iconURL *string
	var gambarKondisiURL *string

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
		if deskripsi := ctx.FormValue("deskripsi"); deskripsi != "" {
			req.Deskripsi = &deskripsi
		}
		if isActiveStr := ctx.FormValue("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}
		// if memilikiKondisi := ctx.FormValue("memiliki_kondisi_tambahan"); memilikiKondisi != "" {
		// 	hasKondisi := memilikiKondisi == "true"
		// 	req.MemilikiKondisiTambahan = &hasKondisi
		// }
		if tipeKondisi := ctx.FormValue("tipe_kondisi_tambahan"); tipeKondisi != "" {
			req.TipeKondisiTambahan = &tipeKondisi
		}
		if teksKondisi := ctx.FormValue("teks_kondisi"); teksKondisi != "" {
			req.TeksKondisi = &teksKondisi
		}

		// Handle icon upload or URL
		if file, err := ctx.FormFile("icon"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file icon tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan icon: "+err.Error(), nil)
			}
			iconURL = &savedPath
		} else if iconStr := ctx.FormValue("icon"); iconStr != "" {
			iconURL = &iconStr
		}

		// Handle gambar kondisi upload or URL
		if file, err := ctx.FormFile("gambar_kondisi"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file gambar kondisi tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "product-categories/kondisi", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan gambar kondisi: "+err.Error(), nil)
			}
			gambarKondisiURL = &savedPath
		} else if gambarKondisiStr := ctx.FormValue("gambar_kondisi"); gambarKondisiStr != "" {
			gambarKondisiURL = &gambarKondisiStr
		}

		// Use UpdateWithIcon for multipart
		result, err := c.service.UpdateWithIcon(ctx.UserContext(), id, &req, iconURL, gambarKondisiURL)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.SuccessResponse(ctx, "Kategori produk berhasil diupdate", result)
	}

	// Handle application/json (no file upload)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Kategori produk berhasil diupdate", result)
}

func (c *KategoriProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "kategori produk tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Kategori produk berhasil dihapus", nil)
}

func (c *KategoriProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Status kategori berhasil diubah", result)
}

func (c *KategoriProdukController) Dropdown(ctx *fiber.Ctx) error {
	kategoriList, err := c.service.FindAllActiveForDropdown(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kategori", nil)
	}

	// Convert to dropdown response
	response := make([]map[string]interface{}, len(kategoriList))
	for i, k := range kategoriList {
		response[i] = map[string]interface{}{
			"id":   k.ID.String(),
			"nama": k.NamaID,
		}
	}

	return utils.SuccessResponse(ctx, "Data dropdown kategori produk berhasil diambil", response)
}
