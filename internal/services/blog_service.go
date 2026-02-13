package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogService interface {
	Create(ctx context.Context, req *dto.CreateBlogRequest) (*dto.BlogResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateBlogRequest) (*dto.BlogResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.BlogResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.BlogResponse, error)
	GetAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, page, limit int) ([]dto.BlogListResponse, int64, error)
	Search(ctx context.Context, keyword string, isActive *bool, page, limit int) ([]dto.BlogListResponse, int64, error)
	IncrementView(ctx context.Context, id uuid.UUID) error
	GetRelated(ctx context.Context, blogID uuid.UUID, limit int) ([]dto.BlogListResponse, error)
	GetPopular(ctx context.Context, limit int) ([]dto.BlogListResponse, error)
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
}

type blogService struct {
	blogRepo     repositories.BlogRepository
	kategoriRepo repositories.KategoriBlogRepository
	labelRepo    repositories.LabelBlogRepository
	sanitizer    *HTMLSanitizer
	cfg          *config.Config
}

func NewBlogService(
	blogRepo repositories.BlogRepository,
	kategoriRepo repositories.KategoriBlogRepository,
	labelRepo repositories.LabelBlogRepository,
	cfg *config.Config,
) BlogService {
	return &blogService{
		blogRepo:     blogRepo,
		kategoriRepo: kategoriRepo,
		labelRepo:    labelRepo,
		sanitizer:    NewHTMLSanitizer(),
		cfg:          cfg,
	}
}

func (s *blogService) Create(ctx context.Context, req *dto.CreateBlogRequest) (*dto.BlogResponse, error) {
	// Validate kategori exists
	_, err := s.kategoriRepo.FindByID(ctx, req.KategoriID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kategori not found")
		}
		return nil, err
	}

	// Validate labels exist
	if len(req.LabelIDs) > 0 {
		labels, err := s.labelRepo.FindByIDs(ctx, req.LabelIDs)
		if err != nil {
			return nil, err
		}
		if len(labels) != len(req.LabelIDs) {
			return nil, errors.New("one or more labels not found")
		}
	}

	// Sanitize HTML content
	sanitizedKontenID := s.sanitizer.Sanitize(req.KontenID)
	var sanitizedKontenEN *string
	if req.KontenEN != nil {
		sanitized := s.sanitizer.Sanitize(*req.KontenEN)
		sanitizedKontenEN = &sanitized
	}

	blog := &models.Blog{
		JudulID:           req.JudulID,
		JudulEN:           req.JudulEN,
		Slug:              req.Slug,
		KontenID:          sanitizedKontenID,
		KontenEN:          sanitizedKontenEN,
		FeaturedImageURL:  req.FeaturedImageURL,
		KategoriID:        req.KategoriID,
		MetaTitleID:       req.MetaTitleID,
		MetaTitleEN:       req.MetaTitleEN,
		MetaDescriptionID: req.MetaDescriptionID,
		MetaDescriptionEN: req.MetaDescriptionEN,
		MetaKeywords:      req.MetaKeywords,
		IsActive:          req.IsActive,
	}

	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return nil, err
	}

	// Add labels
	if len(req.LabelIDs) > 0 {
		if err := s.blogRepo.AddLabels(ctx, blog.ID, req.LabelIDs); err != nil {
			return nil, err
		}
	}

	return s.GetByID(ctx, blog.ID)
}

func (s *blogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateBlogRequest) (*dto.BlogResponse, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.JudulID != nil {
		blog.JudulID = *req.JudulID
	}
	if req.JudulEN != nil {
		blog.JudulEN = req.JudulEN
	}
	if req.Slug != nil {
		blog.Slug = *req.Slug
	}
	if req.KontenID != nil {
		// Sanitize HTML content
		sanitized := s.sanitizer.Sanitize(*req.KontenID)
		blog.KontenID = sanitized
	}
	if req.KontenEN != nil {
		// Sanitize HTML content
		sanitized := s.sanitizer.Sanitize(*req.KontenEN)
		blog.KontenEN = &sanitized
	}
	if req.FeaturedImageURL != nil {
		blog.FeaturedImageURL = req.FeaturedImageURL
	}
	if req.KategoriID != nil {
		blog.KategoriID = *req.KategoriID
	}
	if req.MetaTitleID != nil {
		blog.MetaTitleID = req.MetaTitleID
	}
	if req.MetaTitleEN != nil {
		blog.MetaTitleEN = req.MetaTitleEN
	}
	if req.MetaDescriptionID != nil {
		blog.MetaDescriptionID = req.MetaDescriptionID
	}
	if req.MetaDescriptionEN != nil {
		blog.MetaDescriptionEN = req.MetaDescriptionEN
	}
	if req.MetaKeywords != nil {
		blog.MetaKeywords = req.MetaKeywords
	}
	if req.IsActive != nil {
		blog.IsActive = *req.IsActive
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}

	// Update labels
	if req.LabelIDs != nil {
		if err := s.blogRepo.RemoveAllLabels(ctx, id); err != nil {
			return nil, err
		}
		if len(req.LabelIDs) > 0 {
			if err := s.blogRepo.AddLabels(ctx, id, req.LabelIDs); err != nil {
				return nil, err
			}
		}
	}

	return s.GetByID(ctx, id)
}

func (s *blogService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.blogRepo.Delete(ctx, id)
}

func (s *blogService) GetByID(ctx context.Context, id uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toBlogResponse(blog), nil
}

func (s *blogService) GetBySlug(ctx context.Context, slug string) (*dto.BlogResponse, error) {
	blog, err := s.blogRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return s.toBlogResponse(blog), nil
}

func (s *blogService) GetAll(ctx context.Context, isActive *bool, kategoriID *uuid.UUID, page, limit int) ([]dto.BlogListResponse, int64, error) {
	offset := (page - 1) * limit
	blogs, total, err := s.blogRepo.FindAll(ctx, isActive, kategoriID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.BlogListResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = s.toBlogListResponse(&blog)
	}

	return responses, total, nil
}

func (s *blogService) Search(ctx context.Context, keyword string, isActive *bool, page, limit int) ([]dto.BlogListResponse, int64, error) {
	offset := (page - 1) * limit
	blogs, total, err := s.blogRepo.Search(ctx, keyword, isActive, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.BlogListResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = s.toBlogListResponse(&blog)
	}

	return responses, total, nil
}

func (s *blogService) IncrementView(ctx context.Context, id uuid.UUID) error {
	return s.blogRepo.IncrementViewCount(ctx, id)
}

func (s *blogService) GetRelated(ctx context.Context, blogID uuid.UUID, limit int) ([]dto.BlogListResponse, error) {
	blog, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, err
	}

	blogs, err := s.blogRepo.FindRelated(ctx, blogID, blog.KategoriID, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BlogListResponse, len(blogs))
	for i, b := range blogs {
		responses[i] = s.toBlogListResponse(&b)
	}

	return responses, nil
}

func (s *blogService) GetPopular(ctx context.Context, limit int) ([]dto.BlogListResponse, error) {
	blogs, err := s.blogRepo.FindPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BlogListResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = s.toBlogListResponse(&blog)
	}

	return responses, nil
}

func (s *blogService) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats, err := s.blogRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// Get popular blogs
	popularBlogs, _ := s.blogRepo.FindPopular(ctx, 5)
	stats["blog_populer"] = popularBlogs

	return stats, nil
}

func (s *blogService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	return s.blogRepo.ToggleStatus(ctx, id)
}

func (s *blogService) toBlogResponse(blog *models.Blog) *dto.BlogResponse {
	resp := &dto.BlogResponse{
		ID:                blog.ID,
		JudulID:           blog.JudulID,
		JudulEN:           blog.JudulEN,
		Slug:              blog.Slug,
		KontenID:          blog.KontenID,
		KontenEN:          blog.KontenEN,
		FeaturedImageURL:  utils.GetFileURLPtr(blog.FeaturedImageURL, s.cfg),
		KategoriID:        blog.KategoriID,
		MetaTitleID:       blog.MetaTitleID,
		MetaTitleEN:       blog.MetaTitleEN,
		MetaDescriptionID: blog.MetaDescriptionID,
		MetaDescriptionEN: blog.MetaDescriptionEN,
		MetaKeywords:      blog.MetaKeywords,
		IsActive:          blog.IsActive,
		ViewCount:         blog.ViewCount,
		PublishedAt:       blog.PublishedAt,
		CreatedAt:         blog.CreatedAt,
		UpdatedAt:         blog.UpdatedAt,
	}

	if blog.Kategori != nil {
		resp.Kategori = &dto.KategoriBlogBrief{
			ID:     blog.Kategori.ID,
			NamaID: blog.Kategori.NamaID,
			NamaEN: blog.Kategori.NamaEN,
			Slug:   blog.Kategori.Slug,
		}
	}

	if len(blog.Labels) > 0 {
		resp.Labels = make([]dto.LabelBlogBrief, len(blog.Labels))
		for i, label := range blog.Labels {
			resp.Labels[i] = dto.LabelBlogBrief{
				ID:     label.ID,
				NamaID: label.NamaID,
				NamaEN: label.NamaEN,
				Slug:   label.Slug,
			}
		}
	}

	return resp
}

func (s *blogService) toBlogListResponse(blog *models.Blog) dto.BlogListResponse {
	resp := dto.BlogListResponse{
		ID:               blog.ID,
		JudulID:          blog.JudulID,
		JudulEN:          blog.JudulEN,
		Slug:             blog.Slug,
		FeaturedImageURL: utils.GetFileURLPtr(blog.FeaturedImageURL, s.cfg),
		IsActive:         blog.IsActive,
		ViewCount:        blog.ViewCount,
		PublishedAt:      blog.PublishedAt,
		CreatedAt:        blog.CreatedAt,
	}

	if blog.Kategori != nil {
		resp.Kategori = &dto.KategoriBlogBrief{
			ID:     blog.Kategori.ID,
			NamaID: blog.Kategori.NamaID,
			NamaEN: blog.Kategori.NamaEN,
			Slug:   blog.Kategori.Slug,
		}
	}

	if len(blog.Labels) > 0 {
		resp.Labels = make([]dto.LabelBlogBrief, len(blog.Labels))
		for i, label := range blog.Labels {
			resp.Labels[i] = dto.LabelBlogBrief{
				ID:     label.ID,
				NamaID: label.NamaID,
				NamaEN: label.NamaEN,
				Slug:   label.Slug,
			}
		}
	}

	return resp
}
