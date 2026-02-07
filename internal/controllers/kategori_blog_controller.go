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
	service services.KategoriBlogService
}

func NewKategoriBlogController(service services.KategoriBlogService) *KategoriBlogController {
	return &KategoriBlogController{service: service}
}

func (c *KategoriBlogController) Create(ctx *gin.Context) {
	var req dto.CreateKategoriBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	kategori, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Kategori created successfully", kategori)
}

func (c *KategoriBlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req dto.UpdateKategoriBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	kategori, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to update kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori updated successfully", kategori)
}

func (c *KategoriBlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete kategori", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Kategori deleted successfully", nil)
}

func (c *KategoriBlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	kategori, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Kategori not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get kategori", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori retrieved successfully", kategori)
}

func (c *KategoriBlogController) GetAll(ctx *gin.Context) {
	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	kategoris, err := c.service.GetAll(ctx.Request.Context(), isActive)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get kategoris", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategoris retrieved successfully", kategoris)
}

func (c *KategoriBlogController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.service.ToggleStatus(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to toggle status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori status toggled successfully", nil)
}

func (c *KategoriBlogController) Reorder(ctx *gin.Context) {
	var req dto.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), req.Items); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to reorder", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategori reordered successfully", nil)
}

func (c *KategoriBlogController) GetAllPublic(ctx *gin.Context) {
	kategoris, err := c.service.GetAllPublicWithCount(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get kategoris", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Kategoris retrieved successfully", kategoris)
}
