package controllers

import (
	"errors"
	"net/http"
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

type BlogController struct {
	blogService         services.BlogService
	kategoriBlogService services.KategoriBlogService
	labelBlogService    services.LabelBlogService
	cfg                 *config.Config
}

func NewBlogController(
	blogService services.BlogService,
	kategoriBlogService services.KategoriBlogService,
	labelBlogService services.LabelBlogService,
	cfg *config.Config,
) *BlogController {
	return &BlogController{
		blogService:         blogService,
		kategoriBlogService: kategoriBlogService,
		labelBlogService:    labelBlogService,
		cfg:                 cfg,
	}
}

// Admin endpoints
func (c *BlogController) Create(ctx *gin.Context) {
	var req dto.CreateBlogRequest
	var featuredImageURL *string

	contentType := ctx.GetHeader("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.JudulID = ctx.PostForm("judul_id")
		judulEN := ctx.PostForm("judul_en")
		if judulEN != "" {
			req.JudulEN = &judulEN
		}
		req.Slug = ctx.PostForm("slug")
		req.KontenID = ctx.PostForm("konten_id")
		kontenEN := ctx.PostForm("konten_en")
		if kontenEN != "" {
			req.KontenEN = &kontenEN
		}

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

		// Parse label_ids (array)
		labelIDsStr := ctx.PostFormArray("label_ids")
		req.LabelIDs = []uuid.UUID{}
		for _, idStr := range labelIDsStr {
			if idStr != "" {
				id, err := uuid.Parse(idStr)
				if err == nil {
					req.LabelIDs = append(req.LabelIDs, id)
				}
			}
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
		req.IsActive = ctx.PostForm("is_active") == "true"

		// Validate required fields
		if req.JudulID == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_id wajib diisi", "")
			return
		}
		if req.Slug == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "slug wajib diisi", "")
			return
		}
		if req.KontenID == "" {
			utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "konten_id wajib diisi", "")
			return
		}

		// Handle featured_image upload (optional)
		if file, err := ctx.FormFile("featured_image"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file featured_image tidak didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "blog", c.cfg)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), "")
				return
			}
			featuredImageURL = &savedPath
			req.FeaturedImageURL = featuredImageURL
		}

		blog, err := c.blogService.Create(ctx.Request.Context(), &req)
		if err != nil {
			// Rollback: delete uploaded file if creation fails
			if featuredImageURL != nil {
				utils.DeleteFile(*featuredImageURL, c.cfg)
			}
			utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat blog", err.Error())
			return
		}

		utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Blog berhasil dibuat", blog)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	blog, err := c.blogService.Create(ctx.Request.Context(), &req)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat blog", err.Error())
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Blog berhasil dibuat", blog)
}

func (c *BlogController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	var req dto.UpdateBlogRequest
	var featuredImageURL *string

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
		if kontenID := ctx.PostForm("konten_id"); kontenID != "" {
			req.KontenID = &kontenID
		}
		if kontenEN := ctx.PostForm("konten_en"); kontenEN != "" {
			req.KontenEN = &kontenEN
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

		// Parse label_ids (array)
		labelIDsStr := ctx.PostFormArray("label_ids")
		if len(labelIDsStr) > 0 {
			req.LabelIDs = []uuid.UUID{}
			for _, idStr := range labelIDsStr {
				if idStr != "" {
					id, err := uuid.Parse(idStr)
					if err == nil {
						req.LabelIDs = append(req.LabelIDs, id)
					}
				}
			}
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
		if isActiveStr := ctx.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle featured_image upload (optional)
		if file, err := ctx.FormFile("featured_image"); err == nil {
			if !utils.IsValidImageType(file) {
				utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file featured_image tidak didukung", "")
				return
			}
			savedPath, err := utils.SaveUploadedFile(file, "blog", c.cfg)
			if err != nil {
				utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), "")
				return
			}
			featuredImageURL = &savedPath
			req.FeaturedImageURL = featuredImageURL
		}

		blog, err := c.blogService.Update(ctx.Request.Context(), id, &req)
		if err != nil {
			// Rollback: delete uploaded file if update fails
			if featuredImageURL != nil {
				utils.DeleteFile(*featuredImageURL, c.cfg)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
				return
			}
			utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui blog", utils.GetValidationErrorMessage(err))
			return
		}

		utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil diperbarui", blog)
		return
	}

	// Handle JSON request (for backward compatibility)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	blog, err := c.blogService.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui blog", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil diperbarui", blog)
}

func (c *BlogController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	if err := c.blogService.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus blog", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil dihapus", nil)
}

func (c *BlogController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	blog, err := c.blogService.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SuccessResponse(ctx, "Blog berhasil didapatkan", blog)
}

func (c *BlogController) GetAll(ctx *gin.Context) {
	var params dto.BlogFilterRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
		return
	}

	params.SetDefaults()

	blogs, meta, err := c.blogService.GetAll(ctx.Request.Context(), &params)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Blog berhasil didapatkan", blogs, *meta)
}

func (c *BlogController) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Keyword diperlukan", "")
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
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari blog", utils.GetValidationErrorMessage(err))
		return
	}

	meta := models.NewPaginationMeta(page, limit, total)
	utils.PaginatedSuccessResponse(ctx, "Blog berhasil didapatkan", blogs, meta)
}

// Public endpoints
func (c *BlogController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	blog, err := c.blogService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
		return
	}

	// Increment view count
	_ = c.blogService.IncrementView(ctx.Request.Context(), blog.ID)

	utils.SuccessResponse(ctx, "Blog berhasil didapatkan", blog)
}

func (c *BlogController) GetPopular(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	blogs, err := c.blogService.GetPopular(ctx.Request.Context(), limit)
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog populer", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SuccessResponse(ctx, "Blog populer berhasil didapatkan", blogs)
}

func (c *BlogController) GetStatistics(ctx *gin.Context) {
	stats, err := c.blogService.GetStatistics(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan statistik", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SuccessResponse(ctx, "Statistik berhasil didapatkan", stats)
}

func (c *BlogController) ToggleStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
		return
	}

	if err := c.blogService.ToggleStatus(ctx.Request.Context(), id); err != nil {
		// if blog not found, return 404
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			return
		}
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", utils.GetValidationErrorMessage(err))
		return
	}

	utils.SuccessResponse(ctx, "Status blog berhasil diubah", nil)
}

// GetDropdownOptions returns all active kategori and label for dropdown
func (c *BlogController) GetDropdownOptions(ctx *gin.Context) {
	// Get kategori
	kategoriBlog, err := c.kategoriBlogService.GetAllActive(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil kategori", err.Error())
		return
	}

	// Get label
	labelBlog, err := c.labelBlogService.GetAllActive(ctx.Request.Context())
	if err != nil {
		utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil label", err.Error())
		return
	}

	// Build response
	data := map[string]interface{}{
		"kategori": kategoriBlog,
		"label":    labelBlog,
	}

	utils.SuccessResponse(ctx, "Data dropdown berhasil diambil", data)
}
