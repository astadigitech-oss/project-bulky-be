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

type BannerTipeProdukController struct {
	service        services.BannerTipeProdukService
	reorderService *services.ReorderService
	cfg            *config.Config
	activityLog    services.ActivityLogService
}

func NewBannerTipeProdukController(service services.BannerTipeProdukService, reorderService *services.ReorderService, cfg *config.Config, activityLog services.ActivityLogService) *BannerTipeProdukController {
	return &BannerTipeProdukController{
		service:        service,
		reorderService: reorderService,
		cfg:            cfg,
		activityLog:    activityLog,
	}
}

func (c *BannerTipeProdukController) Create(ctx *fiber.Ctx) error {
	var req models.CreateBannerTipeProdukRequest
	var gambarURL *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.TipeProdukID = ctx.FormValue("tipe_produk_id")
		req.Nama = ctx.FormValue("nama")

		// Validate required fields
		if req.TipeProdukID == "" || req.Nama == "" {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "tipe_produk_id dan nama wajib diisi", nil)
		}

		// Handle file upload
		if file, err := ctx.FormFile("file"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "banners/tipe-produk", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
			}
			gambarURL = &savedPath
		} else {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "File banner wajib diupload", nil)
		}

		req.GambarURL = *gambarURL

		result, err := c.service.Create(ctx.UserContext(), &req)
		if err != nil {
			// Rollback: delete uploaded file if creation fails
			utils.DeleteFile(*gambarURL, c.cfg)
			return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		}

		c.activityLog.Log(ctx, models.ActionCreate, "banner_tipe_produk", "Banner tipe produk berhasil dibuat")
		return utils.CreatedResponse(ctx, "Banner berhasil dibuat", result)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionCreate, "banner_tipe_produk", "Banner tipe produk berhasil dibuat")
	return utils.CreatedResponse(ctx, "Banner berhasil dibuat", result)
}

func (c *BannerTipeProdukController) FindAll(ctx *fiber.Ctx) error {
	// Check if grouped response is requested (no pagination params)
	hasPage := ctx.Query("page") != ""
	hasPerPage := ctx.Query("per_page") != ""
	search := ctx.Query("search")

	// If no pagination params, return grouped response
	if !hasPage && !hasPerPage {
		grouped, meta, err := c.service.FindAllGrouped(ctx.UserContext(), search)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		}

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Data banner tipe produk berhasil diambil",
			"data":    grouped,
			"meta":    meta,
		})
	}

	// Otherwise, return paginated response (old behavior)
	var params models.BannerTipeProdukFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	tipeProdukID := ctx.Query("tipe_produk_id")

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params, tipeProdukID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data banner berhasil diambil", items, *meta)
}

func (c *BannerTipeProdukController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Detail banner berhasil diambil", result)
}

func (c *BannerTipeProdukController) FindByTipeProdukID(ctx *fiber.Ctx) error {
	tipeProdukID := ctx.Params("tipe_produk_id")

	items, err := c.service.FindByTipeProdukID(ctx.UserContext(), tipeProdukID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Banner berhasil diambil", items)
}

func (c *BannerTipeProdukController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateBannerTipeProdukRequest
	var newGambarURL *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with optional file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if tipeProdukID := ctx.FormValue("tipe_produk_id"); tipeProdukID != "" {
			req.TipeProdukID = &tipeProdukID
		}
		if nama := ctx.FormValue("nama"); nama != "" {
			req.Nama = &nama
		}
		if isActiveStr := ctx.FormValue("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle file upload (optional)
		if file, err := ctx.FormFile("file"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe file tidak didukung", nil)
			}
			savedPath, err := utils.SaveUploadedFile(file, "banners/tipe-produk", c.cfg)
			if err != nil {
				return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), nil)
			}
			newGambarURL = &savedPath
			req.GambarURL = newGambarURL
		}

		result, err := c.service.UpdateWithFile(ctx.UserContext(), id, &req, newGambarURL)
		if err != nil {
			// Rollback: delete new uploaded file if update fails
			if newGambarURL != nil {
				utils.DeleteFile(*newGambarURL, c.cfg)
			}
			return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		}

		c.activityLog.Log(ctx, models.ActionUpdate, "banner_tipe_produk", "Banner tipe produk berhasil diupdate")
		return utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	result, err := c.service.Update(ctx.UserContext(), id, &req)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "banner_tipe_produk", "Banner tipe produk berhasil diupdate")
	return utils.SuccessResponse(ctx, "Banner berhasil diupdate", result)
}

func (c *BannerTipeProdukController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.DeleteWithFile(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionDelete, "banner_tipe_produk", "Banner tipe produk berhasil dihapus")
	return utils.SuccessResponse(ctx, "Banner berhasil dihapus", nil)
}

func (c *BannerTipeProdukController) ToggleStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	result, err := c.service.ToggleStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "banner_tipe_produk", "Status banner tipe produk berhasil diubah")
	return utils.SuccessResponse(ctx, "Status banner berhasil diubah", result)
}

func (c *BannerTipeProdukController) Reorder(ctx *fiber.Ctx) error {
	var req models.ReorderRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.Reorder(ctx.UserContext(), &req); err != nil {
		status := http.StatusInternalServerError
		// Return 404 if banner not found
		if err.Error() == "salah satu atau lebih banner tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "items tidak boleh kosong" {
			status = http.StatusBadRequest
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Urutan banner berhasil diubah", nil)
}

func (c *BannerTipeProdukController) ReorderByDirection(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.ReorderByDirectionRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	// Get banner to find its tipe_produk_id
	banner, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Banner tidak ditemukan", nil)
	}

	idUUID, err := utils.ParseUUID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
	}

	tipeProdukUUID, err := utils.ParseUUID(banner.TipeProduk.ID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Tipe produk ID tidak valid", nil)
	}

	// Reorder SCOPED by tipe_produk_id
	result, err := c.reorderService.Reorder(
		ctx.UserContext(),
		"banner_tipe_produk",
		idUUID,
		req.Direction,
		"tipe_produk_id",
		tipeProdukUUID,
	)

	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	return 	utils.SuccessResponse(ctx, "Urutan berhasil diubah", fiber.Map{
		"item": fiber.Map{
			"id":     result.ItemID,
			"urutan": result.ItemUrutan,
		},
		"swapped_with": fiber.Map{
			"id":     result.SwappedID,
			"urutan": result.SwappedUrutan,
		},
	})
}
