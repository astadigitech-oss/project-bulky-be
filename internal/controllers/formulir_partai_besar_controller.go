package controllers

import (
	"net/http"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FormulirPartaiBesarController struct {
	service        services.FormulirPartaiBesarService
	reorderService *services.ReorderService
}

func NewFormulirPartaiBesarController(service services.FormulirPartaiBesarService, reorderService *services.ReorderService) *FormulirPartaiBesarController {
	return &FormulirPartaiBesarController{
		service:        service,
		reorderService: reorderService,
	}
}

// ========================================
// Config (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetConfig(ctx *gin.Context) {
	result, err := c.service.GetConfig(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Konfigurasi formulir berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) UpdateConfig(ctx *gin.Context) {
	var req models.UpdateFormulirConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateConfig(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Konfigurasi formulir berhasil diupdate", result)
}

// ========================================
// Anggaran (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetAnggaranList(ctx *gin.Context) {
	items, err := c.service.FindAllAnggaran(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Data anggaran berhasil diambil", items)
}

func (c *FormulirPartaiBesarController) GetAnggaranByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindAnggaranByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail anggaran berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) CreateAnggaran(ctx *gin.Context) {
	var req models.CreateAnggaranRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.CreateAnggaran(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Anggaran berhasil ditambahkan", result)
}

func (c *FormulirPartaiBesarController) UpdateAnggaran(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateAnggaranRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.UpdateAnggaran(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Anggaran berhasil diupdate", result)
}

func (c *FormulirPartaiBesarController) DeleteAnggaran(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteAnggaran(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "anggaran tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Anggaran berhasil dihapus", nil)
}

func (c *FormulirPartaiBesarController) ReorderAnggaran(ctx *gin.Context) {
	var req models.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	if err := c.service.ReorderAnggaran(ctx.Request.Context(), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan anggaran berhasil diubah", nil)
}

// ========================================
// Submission (Admin)
// ========================================

func (c *FormulirPartaiBesarController) GetSubmissionList(ctx *gin.Context) {
	var params models.FormulirSubmissionFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAllSubmission(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data submission berhasil diambil", items, *meta)
}

func (c *FormulirPartaiBesarController) GetSubmissionDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindSubmissionByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail submission berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) ResendEmail(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.ResendEmail(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Email berhasil dikirim ulang", nil)
}

// ========================================
// Buyer Endpoints
// ========================================

func (c *FormulirPartaiBesarController) GetOptions(ctx *gin.Context) {
	result, err := c.service.GetOptions(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Opsi formulir berhasil diambil", result)
}

func (c *FormulirPartaiBesarController) Submit(ctx *gin.Context) {
	var req models.CreateFormulirSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	// Get buyer ID from context (set by auth middleware)
	var buyerID *uuid.UUID
	if buyerIDValue, exists := ctx.Get("buyer_id"); exists {
		if id, ok := buyerIDValue.(uuid.UUID); ok {
			buyerID = &id
		}
	}

	submissionID, err := c.service.SubmitFormulir(ctx.Request.Context(), buyerID, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Formulir berhasil dikirim. Tim kami akan segera menghubungi Anda.", gin.H{
		"id": submissionID,
	})
}

func (c *FormulirPartaiBesarController) ReorderAnggaranByDirection(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.ReorderByDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	result, err := c.reorderService.Reorder(
		ctx.Request.Context(),
		"formulir_partai_besar_anggaran",
		idUUID,
		req.Direction,
		"",
		nil,
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Urutan berhasil diubah", gin.H{
		"item": gin.H{
			"id":     result.ItemID,
			"urutan": result.ItemUrutan,
		},
		"swapped_with": gin.H{
			"id":     result.SwappedID,
			"urutan": result.SwappedUrutan,
		},
	})
}
