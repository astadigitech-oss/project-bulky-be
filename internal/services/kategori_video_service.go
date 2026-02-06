package services

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type KategoriVideoService interface {
	Create(ctx context.Context, req *dto.CreateKategoriVideoRequest) (*models.KategoriVideo, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriVideoRequest) (*models.KategoriVideo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.KategoriVideo, error)
	GetBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error)
	GetAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error)
	Reorder(ctx context.Context, items []dto.ReorderItem) error
	ToggleStatus(ctx context.Context, id uuid.UUID) error
	GetAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error)
}

type kategoriVideoService struct {
	repo repositories.KategoriVideoRepository
}

func NewKategoriVideoService(repo repositories.KategoriVideoRepository) KategoriVideoService {
	return &kategoriVideoService{repo: repo}
}

func (s *kategoriVideoService) Create(ctx context.Context, req *dto.CreateKategoriVideoRequest) (*models.KategoriVideo, error) {
	kategori := &models.KategoriVideo{
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

func (s *kategoriVideoService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateKategoriVideoRequest) (*models.KategoriVideo, error) {
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

func (s *kategoriVideoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *kategoriVideoService) GetByID(ctx context.Context, id uuid.UUID) (*models.KategoriVideo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *kategoriVideoService) GetBySlug(ctx context.Context, slug string) (*models.KategoriVideo, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *kategoriVideoService) GetAll(ctx context.Context, isActive *bool) ([]models.KategoriVideo, error) {
	return s.repo.FindAll(ctx, isActive)
}

func (s *kategoriVideoService) Reorder(ctx context.Context, items []dto.ReorderItem) error {
	for _, item := range items {
		if err := s.repo.UpdateUrutan(ctx, item.ID, item.Urutan); err != nil {
			return err
		}
	}
	return nil
}

func (s *kategoriVideoService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	kategori, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	kategori.IsActive = !kategori.IsActive
	return s.repo.Update(ctx, kategori)
}

func (s *kategoriVideoService) GetAllPublicWithCount(ctx context.Context) ([]models.KategoriVideo, error) {
	return s.repo.FindAllPublicWithCount(ctx)
}
