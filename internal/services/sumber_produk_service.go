package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type SumberProdukService interface {
	Create(ctx context.Context, req *models.CreateSumberProdukRequest) (*models.SumberProdukResponse, error)
	FindByID(ctx context.Context, id string) (*models.SumberProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.SumberProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.SumberProdukResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateSumberProdukRequest) (*models.SumberProdukResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error)
}

type sumberProdukService struct {
	repo repositories.SumberProdukRepository
}

func NewSumberProdukService(repo repositories.SumberProdukRepository) SumberProdukService {
	return &sumberProdukService{repo: repo}
}

func (s *sumberProdukService) Create(ctx context.Context, req *models.CreateSumberProdukRequest) (*models.SumberProdukResponse, error) {
	slug := utils.GenerateSlug(req.Nama)

	exists, _ := s.repo.ExistsBySlug(ctx, slug, nil)
	if exists {
		return nil, errors.New("sumber produk dengan nama tersebut sudah ada")
	}

	sumber := &models.SumberProduk{
		Nama:      req.Nama,
		Slug:      slug,
		Deskripsi: req.Deskripsi,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, sumber); err != nil {
		return nil, err
	}

	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindByID(ctx context.Context, id string) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}
	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindBySlug(ctx context.Context, slug string) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}
	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.SumberProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	sumbers, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.SumberProdukResponse
	for _, sb := range sumbers {
		items = append(items, *s.toResponse(&sb))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *sumberProdukService) Update(ctx context.Context, id string, req *models.UpdateSumberProdukRequest) (*models.SumberProdukResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}

	if req.Nama != nil {
		newSlug := utils.GenerateSlug(*req.Nama)
		exists, _ := s.repo.ExistsBySlug(ctx, newSlug, &id)
		if exists {
			return nil, errors.New("sumber produk dengan nama tersebut sudah ada")
		}
		sumber.Nama = *req.Nama
		sumber.Slug = newSlug
	}
	if req.Deskripsi != nil {
		sumber.Deskripsi = req.Deskripsi
	}
	if req.IsActive != nil {
		sumber.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, sumber); err != nil {
		return nil, err
	}

	return s.toResponse(sumber), nil
}

func (s *sumberProdukService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("sumber produk tidak ditemukan")
	}

	// TODO: Check if sumber has products

	return s.repo.Delete(ctx, id)
}

func (s *sumberProdukService) ToggleStatus(ctx context.Context, id string) (*models.ToggleStatusResponse, error) {
	sumber, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("sumber produk tidak ditemukan")
	}

	sumber.IsActive = !sumber.IsActive
	if err := s.repo.Update(ctx, sumber); err != nil {
		return nil, err
	}

	return &models.ToggleStatusResponse{
		ID:       sumber.ID.String(),
		IsActive: sumber.IsActive,
	}, nil
}

func (s *sumberProdukService) toResponse(sb *models.SumberProduk) *models.SumberProdukResponse {
	return &models.SumberProdukResponse{
		ID:           sb.ID.String(),
		Nama:         sb.Nama,
		Slug:         sb.Slug,
		Deskripsi:    sb.Deskripsi,
		IsActive:     sb.IsActive,
		JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt:    sb.CreatedAt,
		UpdatedAt:    sb.UpdatedAt,
	}
}
