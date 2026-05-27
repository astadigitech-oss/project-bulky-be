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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogController struct {
	blogService services.BlogService
	cfg         *config.Config
}

func NewBlogController(
	blogService services.BlogService,
	cfg *config.Config,
) *BlogController {
	return &BlogController{
		blogService: blogService,
		cfg:         cfg,
	}
}

// Admin endpoints
func (c *BlogController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateBlogRequest
	var featuredImageURL *string

	contentType := ctx.Get("Content-Type")

	// Handle multipart/form-data (with file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.JudulID = ctx.FormValue("judul_id")
		judulEN := ctx.FormValue("judul_en")
		if judulEN != "" {
			req.JudulEN = &judulEN
		}
		if slugID := ctx.FormValue("slug_id"); slugID != "" {
			req.SlugID = &slugID
		}
		if slugEN := ctx.FormValue("slug_en"); slugEN != "" {
			req.SlugEN = &slugEN
		}
		req.KontenID = ctx.FormValue("konten_id")
		kontenEN := ctx.FormValue("konten_en")
		if kontenEN != "" {
			req.KontenEN = &kontenEN
		}

		// Parse kategori_id
		kategoriIDStr := ctx.FormValue("kategori_id")
		if kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
			}
			req.KategoriID = kategoriID
		}

		// Parse label_ids (array)
		labelIDsStr := formValueArray(ctx, "label_ids")
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
		// if metaTitleID := ctx.FormValue("meta_title_id"); metaTitleID != "" {
		// 	req.MetaTitleID = &metaTitleID
		// }
		// if metaTitleEN := ctx.FormValue("meta_title_en"); metaTitleEN != "" {
		// 	req.MetaTitleEN = &metaTitleEN
		// }
		// if metaDescID := ctx.FormValue("meta_description_id"); metaDescID != "" {
		// 	req.MetaDescriptionID = &metaDescID
		// }
		// if metaDescEN := ctx.FormValue("meta_description_en"); metaDescEN != "" {
		// 	req.MetaDescriptionEN = &metaDescEN
		// }
		// if metaKeywords := ctx.FormValue("meta_keywords"); metaKeywords != "" {
		// 	req.MetaKeywords = &metaKeywords
		// }
		if highlightID := ctx.FormValue("highlight_id"); highlightID != "" {
			req.HighlightID = &highlightID
		}
		if highlightEN := ctx.FormValue("highlight_en"); highlightEN != "" {
			req.HighlightEN = &highlightEN
		}

		// Parse is_active
		req.IsActive = ctx.FormValue("is_active") == "true"

		// Validate required fields
		if req.JudulID == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "judul_id wajib diisi", "")
		}
		if req.KontenID == "" {
			return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "konten_id wajib diisi", "")
		}

		// Handle featured_image upload (optional)
		if file, err := ctx.FormFile("featured_image"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file featured_image tidak didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "blog", c.cfg)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), "")
			}
			featuredImageURL = &savedPath
			req.FeaturedImageURL = featuredImageURL
		}

		blog, err := c.blogService.Create(ctx.UserContext(), &req)
		if err != nil {
			// Rollback: delete uploaded file if creation fails
			if featuredImageURL != nil {
				utils.DeleteFile(*featuredImageURL, c.cfg)
			}
			return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat blog", err.Error())
		}

		return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Blog berhasil dibuat", blog)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	blog, err := c.blogService.Create(ctx.UserContext(), &req)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat blog", err.Error())
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusCreated, "Blog berhasil dibuat", blog)
}

func (c *BlogController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
	}

	var req dto.UpdateBlogRequest
	var featuredImageURL *string

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
		if kontenID := ctx.FormValue("konten_id"); kontenID != "" {
			req.KontenID = &kontenID
		}
		if kontenEN := ctx.FormValue("konten_en"); kontenEN != "" {
			req.KontenEN = &kontenEN
		}

		// Parse kategori_id
		if kategoriIDStr := ctx.FormValue("kategori_id"); kategoriIDStr != "" {
			kategoriID, err := uuid.Parse(kategoriIDStr)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "kategori_id tidak valid", err.Error())
			}
			req.KategoriID = &kategoriID
		}

		// Parse label_ids (array)
		labelIDsStr := formValueArray(ctx, "label_ids")
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
		// if metaTitleID := ctx.FormValue("meta_title_id"); metaTitleID != "" {
		// 	req.MetaTitleID = &metaTitleID
		// }
		// if metaTitleEN := ctx.FormValue("meta_title_en"); metaTitleEN != "" {
		// 	req.MetaTitleEN = &metaTitleEN
		// }
		// if metaDescID := ctx.FormValue("meta_description_id"); metaDescID != "" {
		// 	req.MetaDescriptionID = &metaDescID
		// }
		// if metaDescEN := ctx.FormValue("meta_description_en"); metaDescEN != "" {
		// 	req.MetaDescriptionEN = &metaDescEN
		// }
		// if metaKeywords := ctx.FormValue("meta_keywords"); metaKeywords != "" {
		// 	req.MetaKeywords = &metaKeywords
		// }
		if highlightID := ctx.FormValue("highlight_id"); highlightID != "" {
			req.HighlightID = &highlightID
		}
		if highlightEN := ctx.FormValue("highlight_en"); highlightEN != "" {
			req.HighlightEN = &highlightEN
		}

		// Parse is_active
		if isActiveStr := ctx.FormValue("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true"
			req.IsActive = &isActive
		}

		// Handle featured_image upload (optional)
		if file, err := ctx.FormFile("featured_image"); err == nil {
			if !utils.IsValidImageType(file) {
				return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Tipe file featured_image tidak didukung", "")
			}
			savedPath, err := utils.SaveUploadedFile(file, "blog", c.cfg)
			if err != nil {
				return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error(), "")
			}
			featuredImageURL = &savedPath
			req.FeaturedImageURL = featuredImageURL
		}

		blog, err := c.blogService.Update(ctx.UserContext(), id, &req)
		if err != nil {
			// Rollback: delete uploaded file if update fails
			if featuredImageURL != nil {
				utils.DeleteFile(*featuredImageURL, c.cfg)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
			}
			return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui blog", utils.GetValidationErrorMessage(err))
		}

		return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil diperbarui", blog)
	}

	// Handle JSON request (for backward compatibility)
	if err := BindJSON(ctx, &req); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Data request tidak valid", utils.GetValidationErrorMessage(err))
	}

	blog, err := c.blogService.Update(ctx.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui blog", utils.GetValidationErrorMessage(err))
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil diperbarui", blog)
}

func (c *BlogController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
	}

	if err := c.blogService.Delete(ctx.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus blog", utils.GetValidationErrorMessage(err))
	}

	return utils.SimpleSuccessResponse(ctx, http.StatusOK, "Blog berhasil dihapus", nil)
}

func (c *BlogController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
	}

	blog, err := c.blogService.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
	}

	return utils.SuccessResponse(ctx, "Blog berhasil didapatkan", blog)
}

func (c *BlogController) GetAll(ctx *fiber.Ctx) error {
	var params dto.BlogFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", err.Error())
	}

	params.SetDefaults()

	blogs, meta, err := c.blogService.GetAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
	}

	return utils.PaginatedSuccessResponse(ctx, "Blog berhasil didapatkan", blogs, *meta)
}

func (c *BlogController) Search(ctx *fiber.Ctx) error {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "Keyword diperlukan", "")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	blogs, total, err := c.blogService.Search(ctx.UserContext(), keyword, isActive, page, limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari blog", utils.GetValidationErrorMessage(err))
	}

	meta := models.NewPaginationMeta(page, limit, total)
	return utils.PaginatedSuccessResponse(ctx, "Blog berhasil didapatkan", blogs, meta)
}

// Public endpoints
func (c *BlogController) GetBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	blog, err := c.blogService.GetBySlug(ctx.UserContext(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog", utils.GetValidationErrorMessage(err))
	}

	// Increment view count
	_ = c.blogService.IncrementView(ctx.UserContext(), blog.ID)

	return utils.SuccessResponse(ctx, "Blog berhasil didapatkan", blog)
}

func (c *BlogController) GetPopular(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "5"))

	blogs, err := c.blogService.GetPopular(ctx.UserContext(), limit)
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan blog populer", utils.GetValidationErrorMessage(err))
	}

	return utils.SuccessResponse(ctx, "Blog populer berhasil didapatkan", blogs)
}

func (c *BlogController) GetStatistics(ctx *fiber.Ctx) error {
	stats, err := c.blogService.GetStatistics(ctx.UserContext())
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan statistik", utils.GetValidationErrorMessage(err))
	}

	return utils.SuccessResponse(ctx, "Statistik berhasil didapatkan", stats)
}

func (c *BlogController) ToggleStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.SimpleErrorResponse(ctx, http.StatusBadRequest, "ID tidak valid", utils.GetValidationErrorMessage(err))
	}

	if err := c.blogService.ToggleStatus(ctx.UserContext(), id); err != nil {
		// if blog not found, return 404
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SimpleErrorResponse(ctx, http.StatusNotFound, "Blog tidak ditemukan", utils.GetValidationErrorMessage(err))
		}
		return utils.SimpleErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status", utils.GetValidationErrorMessage(err))
	}

	return utils.SuccessResponse(ctx, "Status blog berhasil diubah", nil)
}
