package controllers

import (
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DokumenKebijakanController struct {
	service services.DokumenKebijakanService
}

func NewDokumenKebijakanController(service services.DokumenKebijakanService) *DokumenKebijakanController {
	return &DokumenKebijakanController{service: service}
}

func (c *DokumenKebijakanController) Create(ctx *gin.Context) {
	var req models.CreateDokumenKebijakanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreatedResponse(ctx, "Dokumen kebijakan berhasil dibuat", result)
}

func (c *DokumenKebijakanController) FindAll(ctx *gin.Context) {
	var params models.PaginationRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
		return
	}

	items, meta, err := c.service.FindAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Data dokumen kebijakan berhasil diambil", items, *meta)
}

func (c *DokumenKebijakanController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Detail dokumen kebijakan berhasil diambil", result)
}

func (c *DokumenKebijakanController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateDokumenKebijakanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Dokumen kebijakan berhasil diupdate", result)
}

func (c *DokumenKebijakanController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "dokumen kebijakan tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(ctx, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Dokumen kebijakan berhasil dihapus", nil)
}

func (c *DokumenKebijakanController) ToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.ToggleStatus(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, "Status dokumen kebijakan berhasil diubah", result)
}

// Public endpoints
func (c *DokumenKebijakanController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	lang := ctx.DefaultQuery("lang", "id")

	result, err := c.service.GetBySlug(ctx.Request.Context(), slug, lang)
	if err != nil {
		msg := err.Error()
		if lang == "en" && msg == "dokumen kebijakan tidak ditemukan" {
			msg = "Policy document not found"
		}
		utils.ErrorResponse(ctx, http.StatusNotFound, msg, nil)
		return
	}

	msg := "Detail dokumen kebijakan berhasil diambil"
	if lang == "en" {
		msg = "Policy document details retrieved successfully"
	}
	utils.SuccessResponse(ctx, msg, result)
}

func (c *DokumenKebijakanController) GetActiveList(ctx *gin.Context) {
	lang := ctx.DefaultQuery("lang", "id")

	items, err := c.service.GetActiveList(ctx.Request.Context(), lang)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	msg := "Daftar dokumen kebijakan aktif berhasil diambil"
	if lang == "en" {
		msg = "Active policy documents list retrieved successfully"
	}
	utils.SuccessResponse(ctx, msg, items)
}
