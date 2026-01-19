package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type ProdukGambarService interface {
	Create(ctx context.Context, produkID string, req *models.CreateProdukGambarRequest) (*models.ProdukGambarResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateProdukGambarRequest) (*models.ProdukGambarResponse, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, req *models.ReorderRequest) error
}

type produkGambarService struct {
	repo repositories.ProdukGambarRepository
	cfg  *config.Config
}

func NewProdukGambarService(repo repositories.ProdukGambarRepository, cfg *config.Config) ProdukGambarService {
	return &produkGambarService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *produkGambarService) Create(ctx context.Context, produkID string, req *models.CreateProdukGambarRequest) (*models.ProdukGambarResponse, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, errors.New("produk_id tidak valid")
	}

	// Auto-increment urutan per produk
	maxUrutan, err := s.repo.GetMaxUrutanByProdukID(ctx, produkID)
	if err != nil {
		return nil, err
	}

	gambar := &models.ProdukGambar{
		ProdukID:  produkUUID,
		GambarURL: req.GambarURL,
		Urutan:    maxUrutan + 1,
		IsPrimary: req.IsPrimary,
	}

	if err := s.repo.Create(ctx, gambar); err != nil {
		return nil, err
	}

	// If this is primary, update others
	if req.IsPrimary {
		s.repo.SetPrimary(ctx, produkID, gambar.ID.String())
	}

	return &models.ProdukGambarResponse{
		ID:        gambar.ID.String(),
		GambarURL: utils.GetFileURL(gambar.GambarURL, s.cfg),
		Urutan:    gambar.Urutan,
		IsPrimary: gambar.IsPrimary,
	}, nil
}

func (s *produkGambarService) Update(ctx context.Context, id string, req *models.UpdateProdukGambarRequest) (*models.ProdukGambarResponse, error) {
	gambar, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("gambar tidak ditemukan")
	}

	if req.IsPrimary != nil {
		gambar.IsPrimary = *req.IsPrimary
		if *req.IsPrimary {
			s.repo.SetPrimary(ctx, gambar.ProdukID.String(), id)
		}
	}

	if err := s.repo.Update(ctx, gambar); err != nil {
		return nil, err
	}

	return &models.ProdukGambarResponse{
		ID:        gambar.ID.String(),
		GambarURL: utils.GetFileURL(gambar.GambarURL, s.cfg),
		Urutan:    gambar.Urutan,
		IsPrimary: gambar.IsPrimary,
	}, nil
}

func (s *produkGambarService) Delete(ctx context.Context, id string) error {
	gambar, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("gambar tidak ditemukan")
	}

	// Check if this is the last image
	count, _ := s.repo.CountByProdukID(ctx, gambar.ProdukID.String())
	if count <= 1 {
		return errors.New("tidak dapat menghapus gambar terakhir. Produk harus memiliki minimal 1 gambar")
	}

	return s.repo.Delete(ctx, id)
}

func (s *produkGambarService) Reorder(ctx context.Context, req *models.ReorderRequest) error {
	return s.repo.UpdateOrder(ctx, req.Items)
}
