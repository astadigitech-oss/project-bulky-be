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

type LabelBlogController struct {
	service        services.LabelBlogService
	reorderService *services.ReorderService
}

func NewLabelBlogController(service services.LabelBlogService, reorderService *services.ReorderService) *LabelBlogController {
	return &LabelBlogController{service: service, reorderService: reorderService}
}

func (c *LabelBlogController) Create(ctx *gin.Context) {
	var req dto.CreateLabelBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	label, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Label berhasil dibuat", label)
}

func (c *LabelBlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.UpdateLabelBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	label, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label berhasil diperbarui", label)
}

func (c *LabelBlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label berhasil dihapus", nil)
}

func (c *LabelBlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	label, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Label berhasil diambil", label)
}

func (c *LabelBlogController) GetAll(ctx *gin.Context) {
	var params dto.LabelBlogFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	labels, meta, err := c.service.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Label berhasil diambil", labels, meta)
}

func (c *LabelBlogController) Reorder(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	var req dto.ReorderDirectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	result, err := c.reorderService.Reorder(
		ctx.Request.Context(),
		"label_blog",
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

	utils.SuccessResponse(ctx, "Urutan label blog berhasil diubah", gin.H{
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

func (c *LabelBlogController) GetDropdownOptions(ctx *gin.Context) {
	labels, err := c.service.GetAllActive(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Label berhasil diambil", labels)
}

func (c *LabelBlogController) GetAllPublic(ctx *gin.Context) {
	labels, err := c.service.GetAllPublicWithCount(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Label berhasil diambil", labels)
}
