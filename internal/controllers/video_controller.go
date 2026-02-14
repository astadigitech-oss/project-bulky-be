package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoController struct {
	videoService         services.VideoService
	kategoriVideoService services.KategoriVideoService
	cfg                  *config.Config
}

func NewVideoController(
	videoService services.VideoService,
	kategoriVideoService services.KategoriVideoService,
	cfg *config.Config,
) *VideoController {
	return &VideoController{
		videoService:         videoService,
		kategoriVideoService: kategoriVideoService,
		cfg:                  cfg,
	}
}

// Admin endpoints
func (c *VideoController) Create(ctx *gin.Context) {
	var req dto.CreateVideoRequest
	var videoURL string
	var thumbnailURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.JudulID = ctx.PostForm("judul_id")
		req.JudulEN = ctx.PostForm("judul_en")
		req.Slug = ctx.PostForm("slug")
		req.DeskripsiID = ctx.PostForm("deskripsi_id")
		req.DeskripsiEN = ctx.PostForm("deskripsi_en")

		// Parse kategori_id
		kategoriIDStr := ctx.PostForm("kategori_id")
		if kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
				return
			}
			req.KategoriID = kategoriID
		}

		// Parse durasi_detik
		durasiStr := ctx.PostForm("durasi_detik")
		if durasiStr != "" {
			durasi, err := strconv.Atoi(durasiStr)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "durasi_detik harus berupa angka", err.Error())
				return
			}
			req.DurasiDetik = durasi
		}

		// Parse meta fields
		if metaTitleID := ctx.PostForm("meta_title_id"); metaTitleID != "" {
			req.MetaTitleID = &metaTitleID
		}
		if metaTitleEN := ctx.PostForm("meta_title_en"); metaTitleEN != "" {
			req.MetaTitleEN = &metaTitleEN
		}
		if metaDescID := ctx.PostForm("meta_description_id"); metaDescID != "" {
			req.MetaDescriptionID = &metaDescID
		}
		if metaDescEN := ctx.PostForm("meta_description_en"); metaDescEN != "" {
			req.MetaDescriptionEN = &metaDescEN
		}
		if metaKeywords := ctx.PostForm("meta_keywords"); metaKeywords != "" {
			req.MetaKeywords = &metaKeywords
		}

		// Parse is_active
		req.IsActive = ctx.PostForm("is_active")

		// Validate required fields
		if req.JudulID == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_id wajib diisi", "")
			return
		}
		if req.JudulEN == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_en wajib diisi", "")
			return
		}
		if req.Slug == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "slug wajib diisi", "")
			return
		}
		if req.DeskripsiID == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "deskripsi_id wajib diisi", "")
			return
		}
		if req.DeskripsiEN == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "deskripsi_en wajib diisi", "")
			return
		}
		if req.IsActive == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "is_active wajib diisi", "")
			return
		}

		// Handle video file upload
		var uploadedFilePath string
		if file, err := ctx.FormFile("video_file"); err == nil {
			// If video file is uploaded
			if !utils.IsValidVideoType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file video tidak didukung. Hanya MP4, MOV, M4V yang didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "video", c.cfg)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file video: "+err.Error(), "")
				return
			}
			videoURL = savedPath
			uploadedFilePath = filepath.Join(c.cfg.UploadPath, filepath.FromSlash(savedPath))

			// Auto-detect video duration if not provided
			if req.DurasiDetik == 0 {
				duration, err := utils.GetVideoDurationInSeconds(uploadedFilePath)
				if err != nil {
					// Rollback: delete video file if duration detection fails
					utils.DeleteFile(videoURL, c.cfg)
					utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Gagal mendeteksi durasi video. Error: "+err.Error()+". Path: "+uploadedFilePath, "")
					return
				}
				req.DurasiDetik = duration
			}
		} else {
			// If no video file, check for video_url string
			videoURL = ctx.PostForm("video_url")
			if videoURL == "" {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "video_file atau video_url wajib diisi", "")
				return
			}
			// For external URL, duration must be provided manually
			if req.DurasiDetik == 0 {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "durasi_detik wajib diisi untuk video URL eksternal", "")
				return
			}
		}
		req.VideoURL = videoURL

		// Handle thumbnail file upload or auto-generate from video
		if file, err := ctx.FormFile("thumbnail_file"); err == nil {
			// Manual upload thumbnail
			if !utils.IsValidImageType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file thumbnail tidak didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "video/thumbnail", c.cfg)
			if err != nil {
				// Rollback: delete video file if thumbnail upload fails
				if videoURL != "" && strings.HasPrefix(videoURL, "/uploads/") {
					utils.DeleteFile(videoURL, c.cfg)
				}
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file thumbnail: "+err.Error(), "")
				return
			}
			thumbnailURL = &savedPath
			req.ThumbnailURL = thumbnailURL
		} else {
			// Check for thumbnail_url string
			if thumbnailURLStr := ctx.PostForm("thumbnail_url"); thumbnailURLStr != "" {
				req.ThumbnailURL = &thumbnailURLStr
			} else if uploadedFilePath != "" {
				// Auto-generate thumbnail from uploaded video file
				generatedThumbnail, err := utils.GenerateThumbnailFromVideo(videoURL, "video/thumbnail", c.cfg)
				if err != nil {
					// Rollback: delete video file if thumbnail generation fails
					if videoURL != "" && strings.HasPrefix(videoURL, "/uploads/") {
						utils.DeleteFile(videoURL, c.cfg)
					}
					utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal auto-generate thumbnail: "+err.Error(), "")
					return
				}
				thumbnailURL = &generatedThumbnail
				req.ThumbnailURL = thumbnailURL
			}
			// If external video URL and no thumbnail provided, thumbnailURL remains nil
		}

		video, err := c.videoService.Create(ctx.Request.Context(), &req)
		if err != nil {
			// Rollback: delete uploaded files if creation fails
			if videoURL != "" && strings.HasPrefix(videoURL, "/uploads/") {
				utils.DeleteFile(videoURL, c.cfg)
			}
			if thumbnailURL != nil && strings.HasPrefix(*thumbnailURL, "/uploads/") {
				utils.DeleteFile(*thumbnailURL, c.cfg)
			}
			utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat video", err.Error())
			return
		}

		utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Video berhasil dibuat", video)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	video, err := c.videoService.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat video", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Video berhasil dibuat", video)
}

func (c *VideoController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req dto.UpdateVideoRequest
	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if judulID := ctx.PostForm("judul_id"); judulID != "" {
			req.JudulID = &judulID
		}
		if judulEN := ctx.PostForm("judul_en"); judulEN != "" {
			req.JudulEN = &judulEN
		}
		if slug := ctx.PostForm("slug"); slug != "" {
			req.Slug = &slug
		}
		if deskripsiID := ctx.PostForm("deskripsi_id"); deskripsiID != "" {
			req.DeskripsiID = &deskripsiID
		}
		if deskripsiEN := ctx.PostForm("deskripsi_en"); deskripsiEN != "" {
			req.DeskripsiEN = &deskripsiEN
		}

		// Parse kategori_id
		if kategoriIDStr := ctx.PostForm("kategori_id"); kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
				return
			}
			req.KategoriID = &kategoriID
		}

		// Parse durasi_detik
		if durasiStr := ctx.PostForm("durasi_detik"); durasiStr != "" {
			durasi, err := strconv.Atoi(durasiStr)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "durasi_detik harus berupa angka", err.Error())
				return
			}
			req.DurasiDetik = &durasi
		}

		// Parse meta fields
		if metaTitleID := ctx.PostForm("meta_title_id"); metaTitleID != "" {
			req.MetaTitleID = &metaTitleID
		}
		if metaTitleEN := ctx.PostForm("meta_title_en"); metaTitleEN != "" {
			req.MetaTitleEN = &metaTitleEN
		}
		if metaDescID := ctx.PostForm("meta_description_id"); metaDescID != "" {
			req.MetaDescriptionID = &metaDescID
		}
		if metaDescEN := ctx.PostForm("meta_description_en"); metaDescEN != "" {
			req.MetaDescriptionEN = &metaDescEN
		}
		if metaKeywords := ctx.PostForm("meta_keywords"); metaKeywords != "" {
			req.MetaKeywords = &metaKeywords
		}

		// Parse is_active
		if isActive := ctx.PostForm("is_active"); isActive != "" {
			req.IsActive = &isActive
		}

		// Handle video file upload
		if file, err := ctx.FormFile("video_file"); err == nil {
			// Get existing video to delete old file later
			existingVideo, err := c.videoService.GetByID(ctx.Request.Context(), id)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
				return
			}

			if !utils.IsValidVideoType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file video tidak didukung. Hanya MP4, MOV, M4V yang didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "video", c.cfg)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file video: "+err.Error(), "")
				return
			}
			req.VideoURL = &savedPath

			// Auto-detect video duration if not provided
			if req.DurasiDetik == nil || *req.DurasiDetik == 0 {
				uploadedFilePath := filepath.Join(c.cfg.UploadPath, filepath.FromSlash(savedPath))
				duration, err := utils.GetVideoDurationInSeconds(uploadedFilePath)
				if err == nil {
					req.DurasiDetik = &duration
				}
				// If auto-detect fails, keep existing duration or use provided value
			}

			// Delete old video file if exists
			if existingVideo.VideoURL != "" && strings.HasPrefix(existingVideo.VideoURL, "/uploads/") {
				utils.DeleteFile(existingVideo.VideoURL, c.cfg)
			}
		} else {
			// Check for video_url string
			if videoURL := ctx.PostForm("video_url"); videoURL != "" {
				req.VideoURL = &videoURL
			}
		}

		// Handle thumbnail file upload or auto-generate from video
		thumbnailUpdated := false
		var oldThumbnail *string

		if file, err := ctx.FormFile("thumbnail_file"); err == nil {
			// Manual upload thumbnail
			// Get existing video to delete old thumbnail later
			existingVideo, err := c.videoService.GetByID(ctx.Request.Context(), id)
			if err != nil && req.VideoURL == nil {
				// Only fail if we haven't already fetched it above
				utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
				return
			}

			if !utils.IsValidImageType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file thumbnail tidak didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "video/thumbnail", c.cfg)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file thumbnail: "+err.Error(), "")
				return
			}
			req.ThumbnailURL = &savedPath
			thumbnailUpdated = true

			// Store old thumbnail for deletion later
			if existingVideo.ThumbnailURL != nil {
				oldThumbnail = existingVideo.ThumbnailURL
			}
		} else if thumbnailURL := ctx.PostForm("thumbnail_url"); thumbnailURL != "" {
			// Thumbnail URL string provided
			req.ThumbnailURL = &thumbnailURL
			thumbnailUpdated = true
		} else if req.VideoURL != nil && strings.HasPrefix(*req.VideoURL, "/uploads/") {
			// New video uploaded but no thumbnail provided → auto-generate
			generatedThumbnail, err := utils.GenerateThumbnailFromVideo(*req.VideoURL, "video/thumbnail", c.cfg)
			if err != nil {
				// Log error but don't fail - keep old thumbnail
				// Could also choose to fail here if thumbnail is critical
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal auto-generate thumbnail: "+err.Error(), "")
				return
			}
			req.ThumbnailURL = &generatedThumbnail
			thumbnailUpdated = true

			// Get existing video to delete old thumbnail
			existingVideo, err := c.videoService.GetByID(ctx.Request.Context(), id)
			if err == nil && existingVideo.ThumbnailURL != nil {
				oldThumbnail = existingVideo.ThumbnailURL
			}
		}
		// If no thumbnail update and no video update → keep old thumbnail (do nothing)

		video, err := c.videoService.Update(ctx.Request.Context(), id, &req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
				return
			}
			utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate video", err.Error())
			return
		}

		// Delete old thumbnail file if it was updated and was an uploaded file
		if thumbnailUpdated && oldThumbnail != nil && *oldThumbnail != "" && strings.HasPrefix(*oldThumbnail, "/uploads/") {
			utils.DeleteFile(*oldThumbnail, c.cfg)
		}

		utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil diupdate", video)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	video, err := c.videoService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil diupdate", video)
}

func (c *VideoController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.videoService.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil dihapus", nil)
}

func (c *VideoController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	video, err := c.videoService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil didapatkan", video)
}

func (c *VideoController) GetAll(ctx *gin.Context) {
	var params dto.VideoFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	params.SetDefaults()

	videos, meta, err := c.videoService.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", utils.GetValidationErrorMessage(err))
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Video berhasil didapatkan", videos, *meta)
}

func (c *VideoController) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Keyword wajib diisi", "")
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
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari video", err.Error())
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Video berhasil ditemukan", videos, meta)
}

func (c *VideoController) GetPopular(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	videos, err := c.videoService.GetPopular(ctx.Request.Context(), limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video populer", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video populer berhasil didapatkan", videos)
}

// Public endpoints
func (c *VideoController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	video, err := c.videoService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", err.Error())
		return
	}

	// Increment view count
	_ = c.videoService.IncrementView(ctx.Request.Context(), video.ID)

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil didapatkan", video)
}

func (c *VideoController) GetStatistics(ctx *gin.Context) {
	stats, err := c.videoService.GetStatistics(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan statistik", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Statistik berhasil didapatkan", stats)
}

func (c *VideoController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
		return
	}

	if err := c.videoService.ToggleStatus(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Status video berhasil diubah", nil)
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

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Data dropdown berhasil diambil", data)
}
