package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

// TipeProdukService interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukService interface {
	FindByID(ctx context.Context, id string) (*models.TipeProdukResponse, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProdukResponse, *models.PaginationMeta, error)
}

type tipeProdukService struct {
	repo repositories.TipeProdukRepository
}

func NewTipeProdukService(repo repositories.TipeProdukRepository) TipeProdukService {
	return &tipeProdukService{repo: repo}
}

// FindByID retrieves a single tipe produk by ID
func (s *tipeProdukService) FindByID(ctx context.Context, id string) (*models.TipeProdukResponse, error) {
	tipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}
	return s.toResponse(tipe), nil
}

// FindBySlug retrieves a single tipe produk by slug
func (s *tipeProdukService) FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error) {
	tipe, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}
	return s.toResponse(tipe), nil
}

// FindAll retrieves all tipe produk with pagination
func (s *tipeProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProdukResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	tipes, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.TipeProdukResponse
	for _, t := range tipes {
		items = append(items, *s.toResponse(&t))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *tipeProdukService) toResponse(t *models.TipeProduk) *models.TipeProdukResponse {
	return &models.TipeProdukResponse{
		ID:           t.ID.String(),
		Nama:         t.Nama,
		Slug:         t.Slug,
		Deskripsi:    t.Deskripsi,
		Urutan:       t.Urutan,
		IsActive:     t.IsActive,
		JumlahProduk: 0, // TODO: Count from produk table
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
