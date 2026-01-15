package services

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/repositories"
)

// TipeProdukService interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukService interface {
	FindAll(ctx context.Context) ([]dto.TipeProdukListDTO, error)
	FindAllWithProduk(ctx context.Context) ([]dto.TipeProdukWithProdukDTO, error)
}

type tipeProdukService struct {
	repo repositories.TipeProdukRepository
}

func NewTipeProdukService(repo repositories.TipeProdukRepository) TipeProdukService {
	return &tipeProdukService{repo: repo}
}

// FindAll retrieves all tipe produk without pagination
// Returns TipeProdukListDTO with simplified fields for list view
func (s *tipeProdukService) FindAll(ctx context.Context) ([]dto.TipeProdukListDTO, error) {
	tipes, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure empty array instead of null
	if tipes == nil {
		tipes = []dto.TipeProdukListDTO{}
	}

	return tipes, nil
}

// FindAllWithProduk retrieves all tipe produk with their products
func (s *tipeProdukService) FindAllWithProduk(ctx context.Context) ([]dto.TipeProdukWithProdukDTO, error) {
	result, err := s.repo.FindAllWithProduk(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure empty array instead of null
	if result == nil {
		result = []dto.TipeProdukWithProdukDTO{}
	}

	return result, nil
}
