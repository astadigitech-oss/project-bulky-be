package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type ProdukDokumenService interface {
	Create(ctx context.Context, produkID string, req *models.CreateProdukDokumenRequest) (*models.ProdukDokumenResponse, error)
	Delete(ctx context.Context, id string) error
}

type produkDokumenService struct {
	repo repositories.ProdukDokumenRepository
}

func NewProdukDokumenService(repo repositories.ProdukDokumenRepository) ProdukDokumenService {
	return &produkDokumenService{repo: repo}
}

func (s *produkDokumenService) Create(ctx context.Context, produkID string, req *models.CreateProdukDokumenRequest) (*models.ProdukDokumenResponse, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, errors.New("produk_id tidak valid")
	}

	dokumen := &models.ProdukDokumen{
		ProdukID:    produkUUID,
		NamaDokumen: req.NamaDokumen,
		FileURL:     req.FileURL,
		TipeFile:    req.TipeFile,
		UkuranFile:  req.UkuranFile,
	}

	if err := s.repo.Create(ctx, dokumen); err != nil {
		return nil, err
	}

	return &models.ProdukDokumenResponse{
		ID:          dokumen.ID.String(),
		NamaDokumen: dokumen.NamaDokumen,
		FileURL:     dokumen.FileURL,
		TipeFile:    dokumen.TipeFile,
		UkuranFile:  dokumen.UkuranFile,
	}, nil
}

func (s *produkDokumenService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("dokumen tidak ditemukan")
	}

	return s.repo.Delete(ctx, id)
}
