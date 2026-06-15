package controllers

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/transcoder"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoController struct {
	videoService         services.VideoService
	kategoriVideoService services.KategoriVideoService
	cfg                  *config.Config
	activityLog          services.ActivityLogService
}

func NewVideoController(
	videoService services.VideoService,
	kategoriVideoService services.KategoriVideoService,
	cfg *config.Config,
	activityLog services.ActivityLogService,
) *VideoController {
	return &VideoController{
		videoService:         videoService,
		kategoriVideoService: kategoriVideoService,
		cfg:                  cfg,
		activityLog:          activityLog,
	}
}

// Admin endpoints
func (c *VideoController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateVideoRequest
	var videoURL string
	var thumbnailURL *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.JudulID = ctx.FormValue("judul_id")
		req.JudulEN = ctx.FormValue("judul_en")
		if slugID := ctx.FormValue("slug_id"); slugID != "" {
			req.SlugID = &slugID
		}
		if slugEN := ctx.FormValue("slug_en"); slugEN != "" {
			req.SlugEN = &slugEN
		}
		req.DeskripsiID = ctx.FormValue("deskripsi_id")
		req.DeskripsiEN = ctx.FormValue("deskripsi_en")

		// Parse kategori_id
		kategoriIDStr := ctx.FormValue("kategori_id")
		if kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
			}
			req.KategoriID = kategoriID
		}

		// Parse meta fields
		if metaTitleID := ctx.FormValue("meta_title_id"); metaTitleID != "" {
			req.MetaTitleID = &metaTitleID
		}
		if metaTitleEN := ctx.FormValue("meta_title_en"); metaTitleEN != "" {
			req.MetaTitleEN = &metaTitleEN
		}
		if metaDescID := ctx.FormValue("meta_description_id"); metaDescID != "" {
			req.MetaDescriptionID = &metaDescID
		}
		if metaDescEN := ctx.FormValue("meta_description_en"); metaDescEN != "" {
			req.MetaDescriptionEN = &metaDescEN
		}
		if metaKeywords := ctx.FormValue("meta_keywords"); metaKeywords != "" {
			req.MetaKeywords = &metaKeywords
		}

		// Parse is_active
		req.IsActive = ctx.FormValue("is_active")

		// Validate required fields
		if req.JudulID == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_id wajib diisi", "")
		}
		if req.JudulEN == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_en wajib diisi", "")
		}
		if req.DeskripsiID == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "deskripsi_id wajib diisi", "")
		}
		if req.DeskripsiEN == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "deskripsi_en wajib diisi", "")
		}
		if req.IsActive == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "is_active wajib diisi", "")
		}

		// Handle video file upload
		var uploadedFilePath string
		if file, err := ctx.FormFile("video_file"); err == nil {
			// If video file is uploaded
			if !utils.IsValidVideoType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file video tidak didukung. Hanya MP4, MOV, M4V yang didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "video", c.cfg)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file video: "+err.Error(), "")
			}
			videoURL = savedPath
			uploadedFilePath = savedPath
		} else {
			// If no video file, check for video_url string
			videoURL = ctx.FormValue("video_url")
			if videoURL == "" {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "video_file atau video_url wajib diisi", "")
			}
		}
		req.VideoURL = videoURL

		// Handle thumbnail file upload or auto-generate from video
		if file, err := ctx.FormFile("thumbnail_file"); err == nil {
			// Manual upload thumbnail
			if !utils.IsValidImageType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file thumbnail tidak didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "video/thumbnail", c.cfg)
			if err != nil {
				// Rollback: delete video file if thumbnail upload fails
				if videoURL != "" && strings.HasPrefix(videoURL, "/uploads/") {
					utils.DeleteFile(videoURL, c.cfg)
				}
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file thumbnail: "+err.Error(), "")
			}
			thumbnailURL = &savedPath
			req.ThumbnailURL = thumbnailURL
		} else {
			// Check for thumbnail_url string
			if thumbnailURLStr := ctx.FormValue("thumbnail_url"); thumbnailURLStr != "" {
				req.ThumbnailURL = &thumbnailURLStr
			} else if uploadedFilePath != "" {
				// Auto-generate thumbnail from uploaded video file (requires ffmpeg)
				generatedThumbnail, err := utils.GenerateThumbnailFromVideo(videoURL, "video/thumbnail", c.cfg)
				if err == nil {
					thumbnailURL = &generatedThumbnail
					req.ThumbnailURL = thumbnailURL
				}
				// If ffmpeg not installed or extraction fails, thumbnail_url stays nil
			}
			// If external video URL and no thumbnail provided, thumbnailURL remains nil
		}

		// Async transcode path: video file was uploaded → create draft + goroutine
		if uploadedFilePath != "" {
			video, err := c.videoService.CreateDraft(ctx.UserContext(), &req)
			if err != nil {
				utils.DeleteFile(uploadedFilePath, c.cfg)
				if thumbnailURL != nil && strings.HasPrefix(*thumbnailURL, "/uploads/") {
					utils.DeleteFile(*thumbnailURL, c.cfg)
				}
				if err.Error() == "kategori not found" {
					return utils.SimpleErrorResponse(ctx, http.StatusUnprocessableEntity, "Kategori tidak ditemukan", "")
				}
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat data video", err.Error())
			}

			go c.runTranscode(video.ID, uploadedFilePath)

			c.activityLog.Log(ctx, models.ActionCreate, "video", "Video sedang diproses (transcode)")
			return utils.SimpleSuccessResponse(ctx, http.StatusAccepted, "Video sedang diproses", fiber.Map{
				"id":               video.ID,
				"transcode_status": "processing",
			})
		}

		// Sync path: external video_url provided → create normally
		video, err := c.videoService.Create(ctx.UserContext(), &req)
		if err != nil {
			if err.Error() == "kategori not found" {
				return utils.SimpleErrorResponse(ctx, http.StatusUnprocessableEntity, "Kategori tidak ditemukan", "")
			}
			return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat video", err.Error())
		}

		c.activityLog.Log(ctx, models.ActionCreate, "video", "Video berhasil dibuat")
		return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Video berhasil dibuat", video)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", utils.GetValidationErrorMessage(err))
	}

	video, err := c.videoService.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat video", utils.GetValidationErrorMessage(err))
	}

	c.activityLog.Log(ctx, models.ActionCreate, "video", "Video berhasil dibuat")
	return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Video berhasil dibuat", video)
}

func (c *VideoController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
	}

	var req dto.UpdateVideoRequest
	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if judulID := ctx.FormValue("judul_id"); judulID != "" {
			req.JudulID = &judulID
		}
		if judulEN := ctx.FormValue("judul_en"); judulEN != "" {
			req.JudulEN = &judulEN
		}
		if slugID := ctx.FormValue("slug_id"); slugID != "" {
			req.SlugID = &slugID
		}
		if slugEN := ctx.FormValue("slug_en"); slugEN != "" {
			req.SlugEN = &slugEN
		}
		if deskripsiID := ctx.FormValue("deskripsi_id"); deskripsiID != "" {
			req.DeskripsiID = &deskripsiID
		}
		if deskripsiEN := ctx.FormValue("deskripsi_en"); deskripsiEN != "" {
			req.DeskripsiEN = &deskripsiEN
		}

		// Parse kategori_id
		if kategoriIDStr := ctx.FormValue("kategori_id"); kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
			}
			req.KategoriID = &kategoriID
		}

		// Parse meta fields
		if metaTitleID := ctx.FormValue("meta_title_id"); metaTitleID != "" {
			req.MetaTitleID = &metaTitleID
		}
		if metaTitleEN := ctx.FormValue("meta_title_en"); metaTitleEN != "" {
			req.MetaTitleEN = &metaTitleEN
		}
		if metaDescID := ctx.FormValue("meta_description_id"); metaDescID != "" {
			req.MetaDescriptionID = &metaDescID
		}
		if metaDescEN := ctx.FormValue("meta_description_en"); metaDescEN != "" {
			req.MetaDescriptionEN = &metaDescEN
		}
		if metaKeywords := ctx.FormValue("meta_keywords"); metaKeywords != "" {
			req.MetaKeywords = &metaKeywords
		}

		// Parse is_active
		if isActive := ctx.FormValue("is_active"); isActive != "" {
			req.IsActive = &isActive
		}

		// Handle video file upload
		if file, err := ctx.FormFile("video_file"); err == nil {
			// Get existing video to delete old file later
			existingVideo, err := c.videoService.GetByID(ctx.UserContext(), id)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			}

			if !utils.IsValidVideoType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file video tidak didukung. Hanya MP4, MOV, M4V yang didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "video", c.cfg)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file video: "+err.Error(), "")
			}
			req.VideoURL = &savedPath

			// Delete old video file if exists
			if existingVideo.VideoURL != "" && strings.HasPrefix(existingVideo.VideoURL, "/uploads/") {
				utils.DeleteFile(existingVideo.VideoURL, c.cfg)
			}
		} else {
			// Check for video_url string
			if videoURL := ctx.FormValue("video_url"); videoURL != "" {
				req.VideoURL = &videoURL
			}
		}

		// Handle thumbnail file upload or auto-generate from video
		thumbnailUpdated := false
		var oldThumbnail *string

		if file, err := ctx.FormFile("thumbnail_file"); err == nil {
			// Manual upload thumbnail
			// Get existing video to delete old thumbnail later
			existingVideo, err := c.videoService.GetByID(ctx.UserContext(), id)
			if err != nil && req.VideoURL == nil {
				// Only fail if we haven't already fetched it above
				return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
			}

			if !utils.IsValidImageType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file thumbnail tidak didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "video/thumbnail", c.cfg)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file thumbnail: "+err.Error(), "")
			}
			req.ThumbnailURL = &savedPath
			thumbnailUpdated = true

			// Store old thumbnail for deletion later
			if existingVideo.ThumbnailURL != nil {
				oldThumbnail = existingVideo.ThumbnailURL
			}
		} else if thumbnailURL := ctx.FormValue("thumbnail_url"); thumbnailURL != "" {
			// Thumbnail URL string provided
			req.ThumbnailURL = &thumbnailURL
			thumbnailUpdated = true
		} else if req.VideoURL != nil && strings.HasPrefix(*req.VideoURL, "/uploads/") {
			// New video uploaded but no thumbnail provided → try auto-generate (requires ffmpeg)
			generatedThumbnail, err := utils.GenerateThumbnailFromVideo(*req.VideoURL, "video/thumbnail", c.cfg)
			if err == nil {
				req.ThumbnailURL = &generatedThumbnail
				thumbnailUpdated = true

				// Get existing video to delete old thumbnail
				existingVideo, err := c.videoService.GetByID(ctx.UserContext(), id)
				if err == nil && existingVideo.ThumbnailURL != nil {
					oldThumbnail = existingVideo.ThumbnailURL
				}
			}
			// If ffmpeg not installed or extraction fails, keep existing thumbnail
		}
		// If no thumbnail update and no video update → keep old thumbnail (do nothing)

		video, err := c.videoService.Update(ctx.UserContext(), id, &req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video not found", err.Error())
			}
			return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate video", err.Error())
		}

		// Delete old thumbnail file if it was updated and was an uploaded file
		if thumbnailUpdated && oldThumbnail != nil && *oldThumbnail != "" && strings.HasPrefix(*oldThumbnail, "/uploads/") {
			utils.DeleteFile(*oldThumbnail, c.cfg)
		}

		c.activityLog.Log(ctx, models.ActionUpdate, "video", "Video berhasil diupdate")
		return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil diupdate", video)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", utils.GetValidationErrorMessage(err))
	}

	video, err := c.videoService.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate video", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "video", "Video berhasil diupdate")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil diupdate", video)
}

func (c *VideoController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.videoService.Delete(ctx.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus video", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionDelete, "video", "Video berhasil dihapus")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil dihapus", nil)
}

func (c *VideoController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	video, err := c.videoService.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", err.Error())
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil didapatkan", video)
}

func (c *VideoController) GetAll(ctx *fiber.Ctx) error {
	var params dto.VideoFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	params.SetDefaults()

	videos, meta, err := c.videoService.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", utils.GetValidationErrorMessage(err))
	}

	return utils.PaginatedSuccessResponse(ctx, "Video berhasil didapatkan", videos, *meta)
}

func (c *VideoController) Search(ctx *fiber.Ctx) error {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Keyword wajib diisi", "")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	videos, total, err := c.videoService.Search(ctx.UserContext(), keyword, isActive, page, limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari video", err.Error())
	}

	meta := models.NewPaginationMeta(page, limit, total)
	return utils.PaginatedSuccessResponse(ctx, "Video berhasil ditemukan", videos, meta)
}

func (c *VideoController) GetPopular(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	videos, err := c.videoService.GetPopular(ctx.UserContext(), limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video populer", err.Error())
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video populer berhasil didapatkan", videos)
}

// Public endpoints
func (c *VideoController) GetBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	video, err := c.videoService.GetBySlug(ctx.UserContext(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan video", err.Error())
	}

	// Increment view count
	_ = c.videoService.IncrementView(ctx.UserContext(), video.ID)

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Video berhasil didapatkan", video)
}

func (c *VideoController) GetStatistics(ctx *fiber.Ctx) error {
	stats, err := c.videoService.GetStatistics(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan statistik", err.Error())
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Statistik berhasil didapatkan", stats)
}

func (c *VideoController) ToggleStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", err.Error())
	}

	if err := c.videoService.ToggleStatus(ctx.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", err.Error())
	}

	c.activityLog.Log(ctx, models.ActionToggleStatus, "video", "Status video berhasil diubah")
	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Status video berhasil diubah", nil)
}

// GetDropdownOptions returns all active kategori for dropdown
func (c *VideoController) GetDropdownOptions(ctx *fiber.Ctx) error {
	// Get kategori
	kategoriVideo, err := c.kategoriVideoService.GetAllActive(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
	}

	// Build response
	data := map[string]interface{}{
		"kategori": kategoriVideo,
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Data dropdown berhasil diambil", data)
}

func (c *VideoController) GetTranscodeStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", err.Error())
	}

	result, err := c.videoService.GetTranscodeStatus(ctx.UserContext(), id)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Video tidak ditemukan", err.Error())
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Status transcode berhasil diambil", result)
}

// runTranscode dijalankan di goroutine background setelah Create menerima file upload.
// Ia melakukan transcode, lalu mengupdate DB dengan hasilnya.
func (c *VideoController) runTranscode(videoID uuid.UUID, rawRelativePath string) {
	rawAbsPath := filepath.Join(c.cfg.UploadPath, filepath.FromSlash(rawRelativePath))

	result, err := transcoder.Transcode(rawAbsPath)
	// Hapus file raw setelah proses transcode (sukses maupun gagal)
	transcoder.Cleanup(rawAbsPath)

	if err != nil {
		_ = c.videoService.MarkFailed(context.Background(), videoID, err.Error())
		return
	}

	// Derive relative URL dari rawRelativePath (konsisten dengan SaveUploadedFile yang menyimpan tanpa prefix "uploads/")
	dir := filepath.Dir(rawRelativePath)
	base := strings.TrimSuffix(filepath.Base(rawRelativePath), filepath.Ext(rawRelativePath))
	relativeStreamURL := filepath.ToSlash(filepath.Join(dir, "stream_"+base+".mp4"))

	_ = c.videoService.MarkReady(context.Background(), videoID, relativeStreamURL, result.DurasiDetik)
}
