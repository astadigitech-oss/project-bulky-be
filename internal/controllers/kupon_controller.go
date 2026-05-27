package controllers

import (
	"errors"
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KuponController struct {
	kuponService services.KuponService
}

func NewKuponController(kuponService services.KuponService) *KuponController {
	return &KuponController{
		kuponService: kuponService,
	}
}

// GetAll retrieves all kupon with pagination and filters
func (c *KuponController) GetAll(ctx *fiber.Ctx) error {
	var params dto.KuponQueryParams

	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}

	// Manual validation for required fields (to catch empty string cases)
	if params.Page < 1 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'page' wajib diisi dan minimal 1")
	}
	if params.PerPage < 1 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' wajib diisi dan minimal 1")
	}
	if params.PerPage > 100 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' maksimal 100")
	}

	params.SetDefaults()

	kupons, meta, err := c.kuponService.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kupon", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Data kupon berhasil diambil", kupons, *meta)
}

// GetByID retrieves a kupon by ID
func (c *KuponController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	kupon, err := c.kuponService.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail kupon", err.Error())
	}

	return utils.SuccessResponse(ctx, "Detail kupon berhasil diambil", kupon)
}

// Create creates a new kupon
func (c *KuponController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateKuponRequest

	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Custom validation
	if err := req.Validate(); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", err.Error())
	}

	kupon, err := c.kuponService.Create(ctx.UserContext(), &req)
	if err != nil {
		if err.Error() == "kode kupon sudah digunakan" {
			return utils.SimpleErrorResponse(ctx, http.StatusConflict, err.Error(), "")
		}
		if err.Error() == "kategori tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Gagal membuat kupon", err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Kupon berhasil dibuat",
		"data":    kupon,
	})
}

// Update updates an existing kupon
func (c *KuponController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var req dto.UpdateKuponRequest

	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Custom validation
	if err := req.Validate(); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", err.Error())
	}

	kupon, err := c.kuponService.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
		}
		if err.Error() == "kode kupon sudah digunakan" {
			return utils.SimpleErrorResponse(ctx, http.StatusConflict, err.Error(), "")
		}
		if err.Error() == "kategori tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Gagal mengupdate kupon", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kupon berhasil diupdate", kupon)
}

// Delete soft deletes a kupon
func (c *KuponController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	err = c.kuponService.Delete(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus kupon", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kupon berhasil dihapus", nil)
}

// ToggleStatus toggles the is_active status of a kupon
func (c *KuponController) ToggleStatus(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	result, err := c.kuponService.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status kupon", err.Error())
	}

	return utils.SuccessResponse(ctx, "Status kupon berhasil diubah", result)
}

// GenerateKode generates a random kupon code
func (c *KuponController) GenerateKode(ctx *fiber.Ctx) error {
	var req dto.GenerateKodeRequest

	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	req.SetDefaults()

	result, err := c.kuponService.GenerateKode(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal generate kode kupon", err.Error())
	}

	return utils.SuccessResponse(ctx, "Kode kupon berhasil digenerate", result)
}

// GetUsages retrieves usage history of a kupon
func (c *KuponController) GetUsages(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	var params dto.KuponUsagesQueryParams

	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
	}

	// Manual validation for required fields (to catch empty string cases)
	if params.Page < 1 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'page' wajib diisi dan minimal 1")
	}
	if params.PerPage < 1 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' wajib diisi dan minimal 1")
	}
	if params.PerPage > 100 {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' maksimal 100")
	}

	result, meta, err := c.kuponService.GetUsages(ctx.UserContext(), id, params.Page, params.PerPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data penggunaan kupon", err.Error())
	}

	return utils.PaginatedSuccessResponse(ctx, "Data penggunaan kupon berhasil diambil", result, *meta)
}

// GetKategoriDropdown retrieves kategori dropdown for kupon form
func (c *KuponController) GetKategoriDropdown(ctx *fiber.Ctx) error {
	kategoris, err := c.kuponService.GetKategoriDropdown(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kategori", err.Error())
	}

	return utils.SuccessResponse(ctx, "Data kategori berhasil diambil", kategoris)
}
