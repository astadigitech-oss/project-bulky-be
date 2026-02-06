package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogController struct {
	blogService services.BlogService
}

func NewBlogController(blogService services.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

// Admin endpoints
func (c *BlogController) Create(ctx *gin.Context) {
	var req dto.CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	blog, err := c.blogService.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create blog", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Blog created successfully", blog)
}

func (c *BlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req dto.UpdateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	blog, err := c.blogService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to update blog", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog updated successfully", blog)
}

func (c *BlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.blogService.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete blog", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog deleted successfully", nil)
}

func (c *BlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	blog, err := c.blogService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get blog", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Blog retrieved successfully", blog)
}

func (c *BlogController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	var kategoriID *uuid.UUID
	if ctx.Query("kategori_id") != "" {
		id, err := uuid.Parse(ctx.Query("kategori_id"))
		if err == nil {
			kategoriID = &id
		}
	}

	blogs, total, err := c.blogService.GetAll(ctx.Request.Context(), isActive, kategoriID, page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get blogs", err.Error())
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Blogs retrieved successfully", blogs, meta)
}

func (c *BlogController) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Keyword is required", "")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	blogs, total, err := c.blogService.Search(ctx.Request.Context(), keyword, isActive, page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to search blogs", err.Error())
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Blogs retrieved successfully", blogs, meta)
}

// Public endpoints
func (c *BlogController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	blog, err := c.blogService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get blog", err.Error())
		return
	}

	// Increment view count
	_ = c.blogService.IncrementView(ctx.Request.Context(), blog.ID)

	utils.SuccessResponse(ctx, "Blog retrieved successfully", blog)
}

func (c *BlogController) GetPopular(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	blogs, err := c.blogService.GetPopular(ctx.Request.Context(), limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get popular blogs", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Popular blogs retrieved successfully", blogs)
}

func (c *BlogController) GetStatistics(ctx *gin.Context) {
	stats, err := c.blogService.GetStatistics(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get statistics", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Statistics retrieved successfully", stats)
}

func (c *BlogController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.blogService.ToggleStatus(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to toggle status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Blog status toggled successfully", nil)
}
