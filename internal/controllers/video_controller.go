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

type VideoController struct {
	videoService          services.VideoService
	kategoriVideoService  services.KategoriVideoService
}

func NewVideoController(
	videoService services.VideoService,
	kategoriVideoService services.KategoriVideoService,
) *VideoController {
	return &VideoController{
		videoService:         videoService,
		kategoriVideoService: kategoriVideoService,
	}
}

// Admin endpoints
func (c *VideoController) Create(ctx *gin.Context) {
	var req dto.CreateVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	video, err := c.videoService.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to create video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Video created successfully", video)
}

func (c *VideoController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req dto.UpdateVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	video, err := c.videoService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to update video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video updated successfully", video)
}

func (c *VideoController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.videoService.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video deleted successfully", nil)
}

func (c *VideoController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	video, err := c.videoService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get video", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Video retrieved successfully", video)
}

func (c *VideoController) GetAll(ctx *gin.Context) {
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

	videos, total, err := c.videoService.GetAll(ctx.Request.Context(), isActive, kategoriID, page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get videos", err.Error())
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Videos retrieved successfully", videos, meta)
}

func (c *VideoController) Search(ctx *gin.Context) {
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

	videos, total, err := c.videoService.Search(ctx.Request.Context(), keyword, isActive, page, limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to search videos", err.Error())
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Videos retrieved successfully", videos, meta)
}

func (c *VideoController) GetPopular(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	videos, err := c.videoService.GetPopular(ctx.Request.Context(), limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get popular videos", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Popular videos retrieved successfully", videos)
}

// Public endpoints
func (c *VideoController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	video, err := c.videoService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get video", err.Error())
		return
	}

	// Increment view count
	_ = c.videoService.IncrementView(ctx.Request.Context(), video.ID)

	utils.SuccessResponse(ctx, "Video retrieved successfully", video)
}

func (c *VideoController) GetStatistics(ctx *gin.Context) {
	stats, err := c.videoService.GetStatistics(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to get statistics", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Statistics retrieved successfully", stats)
}

func (c *VideoController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := c.videoService.ToggleStatus(ctx.Request.Context(), id); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Failed to toggle status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Video status toggled successfully", nil)
}

// GetDropdownOptions returns all active kategori for dropdown
func (c *VideoController) GetDropdownOptions(ctx *gin.Context) {
	// Get kategori
	kategoriVideo, err := c.kategoriVideoService.GetAllActive(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	// Build response
	data := map[string]interface{}{
		"kategori": kategoriVideo,
	}

	utils.SuccessResponse(ctx, "Data dropdown berhasil diambil", data)
}
