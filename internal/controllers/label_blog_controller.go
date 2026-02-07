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
	service services.LabelBlogService
}

func NewLabelBlogController(service services.LabelBlogService) *LabelBlogController {
	return &LabelBlogController{service: service}
}

func (c *LabelBlogController) Create(ctx *gin.Context) {
	var req dto.CreateLabelBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	label, err := c.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Label created successfully", label)
}

func (c *LabelBlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req dto.UpdateLabelBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	label, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to update label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label updated successfully", label)
}

func (c *LabelBlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete label", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Label deleted successfully", nil)
}

func (c *LabelBlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	label, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Label not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get label", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Label retrieved successfully", label)
}

func (c *LabelBlogController) GetAll(ctx *gin.Context) {
	labels, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get labels", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Labels retrieved successfully", labels)
}

func (c *LabelBlogController) Reorder(ctx *gin.Context) {
	var req dto.ReorderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := c.service.Reorder(ctx.Request.Context(), req.Items); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to reorder", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Labels reordered successfully", nil)
}

func (c *LabelBlogController) GetAllPublic(ctx *gin.Context) {
	labels, err := c.service.GetAllPublicWithCount(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get labels", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Labels retrieved successfully", labels)
}
