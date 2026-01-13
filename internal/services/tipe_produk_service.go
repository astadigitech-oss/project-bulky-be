package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

// TipeProdukService interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukService interface {
	FindByID(ctx context.Context, id string) (*dto.TipeProdukDetailDTO, error)
	// FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]dto.TipeProdukListDTO, *models.PaginationMeta, error)
}

type tipeProdukService struct {
	repo repositories.TipeProdukRepository
}

func NewTipeProdukService(repo repositories.TipeProdukRepository) TipeProdukService {
	return &tipeProdukService{repo: repo}
}

// FindByID retrieves a single tipe produk by ID with complete details
// Returns TipeProdukDetailDTO with all fields
func (s *tipeProdukService) FindByID(ctx context.Context, id string) (*dto.TipeProdukDetailDTO, error) {
	// Parse UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	tipe, err := s.repo.FindByID(ctx, uuid)
	if err != nil {
		return nil, errors.New("tipe produk tidak ditemukan")
	}
	return tipe, nil
}

// FindBySlug retrieves a single tipe produk by slug
// func (s *tipeProdukService) FindBySlug(ctx context.Context, slug string) (*models.TipeProdukResponse, error) {
// 	tipe, err := s.repo.FindBySlug(ctx, slug)
// 	if err != nil {
// 		return nil, errors.New("tipe produk tidak ditemukan")
// 	}
// 	return s.toResponse(tipe), nil
// }

// FindAll retrieves all tipe produk with pagination
// Returns TipeProdukListDTO with simplified fields for list view
func (s *tipeProdukService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]dto.TipeProdukListDTO, *models.PaginationMeta, error) {
	params.SetDefaults()

	tipes, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Ensure empty array instead of null
	if tipes == nil {
		tipes = []dto.TipeProdukListDTO{}
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return tipes, &meta, nil
}
