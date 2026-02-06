package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type KategoriBlogService interface {
	Create(ctx context.Context, req *dto.CreateKategoriBlogRequest) (*models.KategoriBlog, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriBlogRequest) (*models.KategoriBlog, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.KategoriBlog, error)
	GetBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error)
	GetAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error)
	Reorder(ctx context.Context, items []dto.ReorderItem) error
	ToggleStatus(ctx context.Context, id uuid.UUID) error
	GetAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error)
}

type kategoriBlogService struct {
	repo repositories.KategoriBlogRepository
}

func NewKategoriBlogService(repo repositories.KategoriBlogRepository) KategoriBlogService {
	return &kategoriBlogService{repo: repo}
}

func (s *kategoriBlogService) Create(ctx context.Context, req *dto.CreateKategoriBlogRequest) (*models.KategoriBlog, error) {
	kategori := &models.KategoriBlog{
		NamaID:   req.NamaID,
		NamaEN:   req.NamaEN,
		Slug:     req.Slug,
		IsActive: req.IsActive,
		Urutan:   req.Urutan,
	}

	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return kategori, nil
}

func (s *kategoriBlogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriBlogRequest) (*models.KategoriBlog, error) {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.NamaID != nil {
		kategori.NamaID = *req.NamaID
	}
	if req.NamaEN != nil {
		kategori.NamaEN = req.NamaEN
	}
	if req.Slug != nil {
		kategori.Slug = *req.Slug
	}
	if req.IsActive != nil {
		kategori.IsActive = *req.IsActive
	}
	if req.Urutan != nil {
		kategori.Urutan = *req.Urutan
	}

	if err := s.repo.Update(ctx, kategori); err != nil {
		return nil, err
	}

	return kategori, nil
}

func (s *kategoriBlogService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if kategori has blogs
	count, err := s.repo.CountBlogByKategori(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("kategori tidak dapat dihapus karena masih memiliki artikel blog")
	}
	return s.repo.Delete(ctx, id)
}

func (s *kategoriBlogService) GetByID(ctx context.Context, id uuid.UUID) (*models.KategoriBlog, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *kategoriBlogService) GetBySlug(ctx context.Context, slug string) (*models.KategoriBlog, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *kategoriBlogService) GetAll(ctx context.Context, isActive *bool) ([]models.KategoriBlog, error) {
	return s.repo.FindAll(ctx, isActive)
}

func (s *kategoriBlogService) Reorder(ctx context.Context, items []dto.ReorderItem) error {
	for _, item := range items {
		if err := s.repo.UpdateUrutan(ctx, item.ID, item.Urutan); err != nil {
			return err
		}
	}
	return nil
}

func (s *kategoriBlogService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	kategori.IsActive = !kategori.IsActive
	return s.repo.Update(ctx, kategori)
}

func (s *kategoriBlogService) GetAllPublicWithCount(ctx context.Context) ([]models.KategoriBlog, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
