package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type MerekProdukService interface {
	Create(ctx context.Context, req *models.CreateMerekProdukRequest) (*models.MerekProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.MerekProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.MerekProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProdukResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateMerekProdukRequest) (*models.MerekProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type merekProdukService struct {
	repo repositories.MerekProdukRepository
}

func NewMerekProdukService(repo repositories.MerekProdukRepository) MerekProdukService {
	return &merekProdukService{repo: repo}
}

func (s *merekProdukService) Create(ctx context.Context, req *models.CreateMerekProdukRequest) (*models.MerekProdukResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("merek dengan nama tersebut sudah ada")
	}

	merek := &models.MerekProduk{
		Nama:     req.Nama,
		Slug:     slug,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindByID(ctx context.Context, id string) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}
	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindBySlug(ctx context.Context, slug string) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}
	return s.toResponse(merek), nil
}

func (s *merekProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	mereks, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.MerekProdukResponse
	for _, m := range mereks {
		items = append(items, *s.toResponse(&m))
	}

	totalHalaman := (total + int64(params.PerHalaman) - 1) / int64(params.PerHalaman)

	meta := &models.PaginationMeta{
		Halaman:      params.Halaman,
		PerHalaman:   params.PerHalaman,
		TotalData:    total,
		TotalHalaman: totalHalaman,
	}

	return items, meta, nil
}

func (s *merekProdukService) Update(ctx context.Context, id string, req *models.UpdateMerekProdukRequest) (*models.MerekProdukResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("merek dengan nama tersebut sudah ada")
		}
		merek.Nama = *req.Nama
		merek.Slug = newSlug
	}
	if req.IsActive != nil {
		merek.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, merek); err != nil {
		return nil, err
	}

	return s.toResponse(merek), nil
}

func (s *merekProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("merek produk tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}

func (s *merekProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	merek, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("merek produk tidak ditemukan")
	}

	merek.IsActive = !merek.IsActive
	if err := s.repo.Update(ctx, merek); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       merek.ID.String(),
		IsActive: merek.IsActive,
	}, nil
}

func (s *merekProdukService) toResponse(m *models.MerekProduk) *models.MerekProdukResponse {
	return &models.MerekProdukResponse{
		ID:           m.ID.String(),
		Nama:         m.Nama,
		Slug:         m.Slug,
		LogoURL:      m.LogoURL,
		IsActive:     m.IsActive,
		JumlahProduk: 0,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
