package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type LabelBlogService interface {
	Create(ctx context.Context, req *dto.CreateLabelBlogRequest) (*models.LabelBlog, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateLabelBlogRequest) (*models.LabelBlog, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LabelBlog, error)
	GetBySlug(ctx context.Context, slug string) (*models.LabelBlog, error)
	GetAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]models.LabelBlog, models.PaginationMeta, error)
	GetAllActive(ctx context.Context) ([]dto.LabelBlogDropdownResponse, error)
	GetAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error)
}

type labelBlogService struct {
	repo repositories.LabelBlogRepository
}

func NewLabelBlogService(repo repositories.LabelBlogRepository) LabelBlogService {
	return &labelBlogService{repo: repo}
}

func (s *labelBlogService) Create(ctx context.Context, req *dto.CreateLabelBlogRequest) (*models.LabelBlog, error) {
	label := &models.LabelBlog{
		NamaID: req.NamaID,
		NamaEN: req.NamaEN,
		Slug:   req.Slug,
		Urutan: req.Urutan,
	}

	if err := s.repo.Create(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelBlogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateLabelBlogRequest) (*models.LabelBlog, error) {
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NamaID != nil {
		label.NamaID = *req.NamaID
	}
	if req.NamaEN != nil {
		label.NamaEN = *req.NamaEN
	}
	if req.Slug != nil {
		label.Slug = *req.Slug
	}
	if req.Urutan != nil {
		label.Urutan = *req.Urutan
	}

	if err := s.repo.Update(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelBlogService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if label exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if label is used in any blog posts
	count, err := s.repo.CountBlogByLabel(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("label tidak dapat dihapus karena masih memiliki artikel blog")
	}

	return s.repo.Delete(ctx, id)
}

func (s *labelBlogService) GetByID(ctx context.Context, id uuid.UUID) (*models.LabelBlog, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *labelBlogService) GetBySlug(ctx context.Context, slug string) (*models.LabelBlog, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *labelBlogService) GetAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]models.LabelBlog, models.PaginationMeta, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *labelBlogService) GetAllActive(ctx context.Context) ([]dto.LabelBlogDropdownResponse, error) {
	labelList, _, err := s.repo.FindAll(ctx, &dto.LabelBlogFilterRequest{})
	if err != nil {
		return nil, err
	}

	var result []dto.LabelBlogDropdownResponse
	for _, l := range labelList {
		nama := map[string]interface{}{
			"id": l.NamaID,
		}
		if l.NamaEN != "" {
			nama["en"] = l.NamaEN
		} else {
			nama["en"] = l.NamaID
		}

		result = append(result, dto.LabelBlogDropdownResponse{
			ID:   l.ID,
			Nama: nama,
			Slug: l.Slug,
		})
	}

	return result, nil
}

func (s *labelBlogService) GetAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
