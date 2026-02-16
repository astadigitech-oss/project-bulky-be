package controllers

import (
	"errors"
	"net/http"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
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
func (c *KuponController) GetAll(ctx *gin.Context) {
	var params dto.KuponQueryParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
		return
	}

	// Manual validation for required fields (to catch empty string cases)
	if params.Page < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage > 100 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' maksimal 100")
		return
	}

	params.SetDefaults()

	kupons, meta, err := c.kuponService.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kupon", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data kupon berhasil diambil", kupons, *meta)
}

// GetByID retrieves a kupon by ID
func (c *KuponController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	kupon, err := c.kuponService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail kupon", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Detail kupon berhasil diambil", kupon)
}

// Create creates a new kupon
func (c *KuponController) Create(ctx *gin.Context) {
	var req dto.CreateKuponRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Custom validation
	if err := req.Validate(); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	kupon, err := c.kuponService.Create(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "kode kupon sudah digunakan" {
			utils.SimpleErrorResponse(ctx, http.StatusConflict, err.Error(), "")
			return
		}
		if err.Error() == "kategori tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Gagal membuat kupon", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kupon berhasil dibuat",
		"data":    kupon,
	})
}

// Update updates an existing kupon
func (c *KuponController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.UpdateKuponRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Custom validation
	if err := req.Validate(); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	kupon, err := c.kuponService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
			return
		}
		if err.Error() == "kode kupon sudah digunakan" {
			utils.SimpleErrorResponse(ctx, http.StatusConflict, err.Error(), "")
			return
		}
		if err.Error() == "kategori tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Gagal mengupdate kupon", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kupon berhasil diupdate", kupon)
}

// Delete soft deletes a kupon
func (c *KuponController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	err = c.kuponService.Delete(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus kupon", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kupon berhasil dihapus", nil)
}

// ToggleStatus toggles the is_active status of a kupon
func (c *KuponController) ToggleStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	result, err := c.kuponService.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status kupon", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Status kupon berhasil diubah", result)
}

// GenerateKode generates a random kupon code
func (c *KuponController) GenerateKode(ctx *gin.Context) {
	var req dto.GenerateKodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	req.SetDefaults()

	result, err := c.kuponService.GenerateKode(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal generate kode kupon", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kode kupon berhasil digenerate", result)
}

// GetUsages retrieves usage history of a kupon
func (c *KuponController) GetUsages(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var params dto.KuponUsagesQueryParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", parseValidationErrors(err))
		return
	}

	// Manual validation for required fields (to catch empty string cases)
	if params.Page < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage < 1 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' wajib diisi dan minimal 1")
		return
	}
	if params.PerPage > 100 {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", "Parameter 'per_page' maksimal 100")
		return
	}

	result, meta, err := c.kuponService.GetUsages(ctx.Request.Context(), id, params.Page, params.PerPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "kupon tidak ditemukan" {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kupon tidak ditemukan", "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data penggunaan kupon", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data penggunaan kupon berhasil diambil", result, *meta)
}

// GetKategoriDropdown retrieves kategori dropdown for kupon form
func (c *KuponController) GetKategoriDropdown(ctx *gin.Context) {
	kategoris, err := c.kuponService.GetKategoriDropdown(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data kategori", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Data kategori berhasil diambil", kategoris)
}
