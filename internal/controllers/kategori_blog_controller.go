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

type KategoriBlogController struct {
	service        services.KategoriBlogService
	reorderService *services.ReorderService
}

func NewKategoriBlogController(service services.KategoriBlogService, reorderService *services.ReorderService) *KategoriBlogController {
	return &KategoriBlogController{
		service:        service,
		reorderService: reorderService,
	}
}

func (c *KategoriBlogController) Create(ctx *gin.Context) {
	var req dto.CreateKategoriBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
		return
	}

	kategori, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Kategori berhasil dibuat", kategori)
}

func (c *KategoriBlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.UpdateKategoriBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
		return
	}

	kategori, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil diupdate", kategori)
}

func (c *KategoriBlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori berhasil dihapus", nil)
}

func (c *KategoriBlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	kategori, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategori)
}

func (c *KategoriBlogController) GetAll(ctx *gin.Context) {
	var params dto.KategoriBlogFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	kategoris, meta, err := c.service.GetAllPaginated(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Kategori berhasil diambil", kategoris, meta)
}

func (c *KategoriBlogController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.service.ToggleStatus(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Status kategori berhasil diubah", nil)
}

func (c *KategoriBlogController) Reorder(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.ReorderDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Request tidak valid", err.Error())
		return
	}

	result, err := c.reorderService.Reorder(
		ctx.Request.Context(),
		"kategori_blog",
		id,
		req.Direction,
		"",
		nil,
	)
	if err != nil {
		if err.Error() == "item sudah berada di urutan paling atas" ||
			err.Error() == "item sudah berada di urutan paling bawah" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, err.Error(), "")
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah urutan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Urutan kategori berhasil diubah", gin.H{
		"item": gin.H{
			"id":     result.ItemID.String(),
			"urutan": result.ItemUrutan,
		},
		"swapped": gin.H{
			"id":     result.SwappedID.String(),
			"urutan": result.SwappedUrutan,
		},
	})
}

func (c *KategoriBlogController) GetDropdownOptions(ctx *gin.Context) {
	kategoris, err := c.service.GetAllActive(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategoris)
}

func (c *KategoriBlogController) GetAllPublic(ctx *gin.Context) {
	kategoris, err := c.service.GetAllPublicWithCount(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori berhasil diambil", kategoris)
}
