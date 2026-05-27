package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UlasanController struct {
	service services.UlasanService
}

func NewUlasanController(service services.UlasanService) *UlasanController {
	return &UlasanController{service: service}
}

// ========================================
// Admin Endpoints
// ========================================

// AdminFindAll godoc
// @Summary List all ulasan for admin
// @Tags Admin - Ulasan
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param is_approved query string false "Filter by approval status (true/false/all)"
// @Param rating query int false "Filter by rating (1-5)"
// @Param cari query string false "Search by buyer name, order code, or product name"
// @Param tanggal_dari query string false "Filter from date (YYYY-MM-DD)"
// @Param tanggal_sampai query string false "Filter to date (YYYY-MM-DD)"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /admin/ulasan [get]
func (ctrl *UlasanController) AdminFindAll(c *fiber.Ctx) error {
	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	// Parse filters
	filters := make(map[string]interface{})

	if isApprovedStr := c.Query("is_approved"); isApprovedStr != "" && isApprovedStr != "all" {
		if isApprovedStr == "true" {
			filters["is_approved"] = true
		} else if isApprovedStr == "false" {
			filters["is_approved"] = false
		}
	}

	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			filters["rating"] = rating
		}
	}

	if cari := c.Query("cari"); cari != "" {
		filters["cari"] = cari
	}

	if tanggalDari := c.Query("tanggal_dari"); tanggalDari != "" {
		if t, err := time.Parse("2006-01-02", tanggalDari); err == nil {
			filters["tanggal_dari"] = t
		}
	}

	if tanggalSampai := c.Query("tanggal_sampai"); tanggalSampai != "" {
		if t, err := time.Parse("2006-01-02", tanggalSampai); err == nil {
			filters["tanggal_sampai"] = t.Add(24 * time.Hour).Add(-time.Second) // End of day
		}
	}

	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Get data
	data, total, summary, err := ctrl.service.AdminFindAll(filters, page, perPage, sortBy, sortOrder)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
	}

	// Response
	meta := models.NewPaginationMeta(page, perPage, total)
	return utils.PaginatedSuccessResponseWithSummary(c, "Data ulasan berhasil diambil", data, meta, summary)
}

// AdminFindByID godoc
// @Summary Get ulasan detail for admin
// @Tags Admin - Ulasan
// @Security BearerAuth
// @Param id path string true "Ulasan ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/ulasan/{id} [get]
func (ctrl *UlasanController) AdminFindByID(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := ctrl.service.AdminFindByID(id)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusNotFound, "Ulasan tidak ditemukan", err.Error())
	}

	return utils.SuccessResponse(c, "Data ulasan berhasil diambil", data)
}

// Approve godoc
// @Summary Approve or reject ulasan
// @Tags Admin - Ulasan
// @Security BearerAuth
// @Param id path string true "Ulasan ID"
// @Param request body models.ApproveUlasanRequest true "Approval data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/ulasan/{id}/approve [patch]
func (ctrl *UlasanController) Approve(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.ApproveUlasanRequest
	if err := BindJSON(c, &req); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Get admin ID from context
	adminID := c.Locals("admin_id")
	if adminID == nil {
		return utils.SimpleErrorResponse(c, http.StatusUnauthorized, "Admin ID not found", "")
	}

	adminUUID, err := uuid.Parse(adminID.(string))
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Invalid admin ID", err.Error())
	}

	if err := ctrl.service.Approve(id, req.IsApproved, adminUUID); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate status ulasan", err.Error())
	}

	message := "Ulasan berhasil diapprove"
	if !req.IsApproved {
		message = "Ulasan berhasil direject"
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, message, fiber.Map{
		"id":          id,
		"is_approved": req.IsApproved,
		"approved_at": time.Now(),
		"approved_by": adminUUID.String(),
	})
}

// BulkApprove godoc
// @Summary Bulk approve or reject ulasan
// @Tags Admin - Ulasan
// @Security BearerAuth
// @Param request body models.BulkApproveUlasanRequest true "Bulk approval data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/ulasan/bulk-approve [patch]
func (ctrl *UlasanController) BulkApprove(c *fiber.Ctx) error {
	var req models.BulkApproveUlasanRequest
	if err := BindJSON(c, &req); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	// Get admin ID from context
	adminID := c.Locals("admin_id")
	if adminID == nil {
		return utils.SimpleErrorResponse(c, http.StatusUnauthorized, "Admin ID not found", "")
	}

	adminUUID, err := uuid.Parse(adminID.(string))
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Invalid admin ID", err.Error())
	}

	affected, err := ctrl.service.BulkApprove(req.IDs, req.IsApproved, adminUUID)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal bulk update ulasan", err.Error())
	}

	message := "Ulasan berhasil diapprove"
	if !req.IsApproved {
		message = "Ulasan berhasil direject"
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, message, fiber.Map{
		"total_updated": affected,
	})
}

// Delete godoc
// @Summary Delete ulasan (soft delete)
// @Tags Admin - Ulasan
// @Security BearerAuth
// @Param id path string true "Ulasan ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/ulasan/{id} [delete]
func (ctrl *UlasanController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.service.Delete(id); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusNotFound, "Gagal menghapus ulasan", err.Error())
	}

	return utils.SuccessResponse(c, "Ulasan berhasil dihapus", nil)
}

// ========================================
// Buyer Endpoints
// ========================================

// GetPendingReviews godoc
// @Summary Get pending review items for buyer
// @Tags Buyer - Ulasan
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /buyer/ulasan/pending [get]
func (ctrl *UlasanController) GetPendingReviews(c *fiber.Ctx) error {
	// Get buyer ID from context
	buyerID := c.Locals("buyer_id")
	if buyerID == nil {
		return utils.SimpleErrorResponse(c, http.StatusUnauthorized, "Buyer ID not found", "")
	}

	buyerUUID, err := uuid.Parse(buyerID.(string))
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Invalid buyer ID", err.Error())
	}

	data, err := ctrl.service.GetPendingReviews(buyerUUID)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data", err.Error())
	}

	return utils.SuccessResponse(c, "Data item belum di-review", data)
}

// BuyerFindAll godoc
// @Summary Get buyer's ulasan list
// @Tags Buyer - Ulasan
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /buyer/ulasan [get]
func (ctrl *UlasanController) BuyerFindAll(c *fiber.Ctx) error {
	// Get buyer ID from context
	buyerID := c.Locals("buyer_id")
	if buyerID == nil {
		return utils.SimpleErrorResponse(c, http.StatusUnauthorized, "Buyer ID not found", "")
	}

	buyerUUID, err := uuid.Parse(buyerID.(string))
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Invalid buyer ID", err.Error())
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	data, total, err := ctrl.service.BuyerFindAll(buyerUUID, page, perPage)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data ulasan", err.Error())
	}

	meta := models.NewPaginationMeta(page, perPage, total)
	return utils.PaginatedSuccessResponse(c, "Data ulasan berhasil diambil", data, meta)
}

// Create godoc
// @Summary Create new ulasan
// @Tags Buyer - Ulasan
// @Security BearerAuth
// @Accept multipart/form-data
// @Param pesanan_item_id formData string true "Pesanan Item ID"
// @Param rating formData int true "Rating (1-5)"
// @Param komentar formData string false "Review comment (max 1000 chars)"
// @Param gambar formData file false "Review image (jpg, png, max 2MB)"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /buyer/ulasan [post]
func (ctrl *UlasanController) Create(c *fiber.Ctx) error {
	// Get buyer ID from context
	buyerID := c.Locals("buyer_id")
	if buyerID == nil {
		return utils.SimpleErrorResponse(c, http.StatusUnauthorized, "Buyer ID not found", "")
	}

	buyerUUID, err := uuid.Parse(buyerID.(string))
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Invalid buyer ID", err.Error())
	}

	// Parse form data
	var req models.CreateUlasanRequest

	req.PesananItemID = c.FormValue("pesanan_item_id")
	ratingStr := c.FormValue("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Rating harus antara 1-5", "")
	}
	req.Rating = rating

	if komentar := c.FormValue("komentar"); komentar != "" {
		req.Komentar = &komentar
	}

	// Validate
	if req.PesananItemID == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "pesanan_item_id wajib diisi", "")
	}

	// Handle image upload
	file, _ := c.FormFile("gambar")

	// Create ulasan
	ulasan, err := ctrl.service.Create(req, buyerUUID, file)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Gagal membuat ulasan", err.Error())
	}

	return utils.SimpleSuccessResponse(c, http.StatusCreated, "Ulasan berhasil dikirim. Menunggu approval admin.", fiber.Map{
		"id":          ulasan.ID.String(),
		"rating":      ulasan.Rating,
		"komentar":    ulasan.Komentar,
		"is_approved": ulasan.IsApproved,
		"created_at":  ulasan.CreatedAt,
	})
}

// ========================================
// Public Endpoints
// ========================================

// GetProdukUlasan godoc
// @Summary Get product reviews (public)
// @Tags Public - Ulasan
// @Param produk_id path string true "Produk ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(5)
// @Param rating query int false "Filter by rating (1-5)"
// @Param with_photo query bool false "Filter only with photo"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /public/produk/{produk_id}/ulasan [get]
func (ctrl *UlasanController) GetProdukUlasan(c *fiber.Ctx) error {
	produkID := c.Params("produk_id")

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	// Parse filters
	filters := make(map[string]interface{})

	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			filters["rating"] = rating
		}
	}

	if withPhotoStr := c.Query("with_photo"); withPhotoStr == "true" {
		filters["with_photo"] = true
	}

	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Get data
	data, total, err := ctrl.service.GetProdukUlasan(produkID, filters, page, perPage, sortBy, sortOrder)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Gagal mengambil data ulasan", err.Error())
	}

	meta := models.NewPaginationMeta(page, perPage, total)
	return utils.PaginatedSuccessResponse(c, "Data ulasan berhasil diambil", data, meta)
}

// GetProdukRating godoc
// @Summary Get product rating stats (public)
// @Tags Public - Ulasan
// @Param produk_id path string true "Produk ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /public/produk/{produk_id}/rating [get]
func (ctrl *UlasanController) GetProdukRating(c *fiber.Ctx) error {
	produkID := c.Params("produk_id")

	data, err := ctrl.service.GetProdukRatingStats(produkID)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Gagal mengambil rating produk", err.Error())
	}

	return utils.SuccessResponse(c, "Data rating berhasil diambil", data)
}
