package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type ProdukGambarService interface {
	CreateWithFile(ctx context.Context, produkID string, file *multipart.FileHeader, isPrimary bool) (*models.ProdukGambarResponse, error)
	Delete(ctx context.Context, produkID, gambarID string) error
	Reorder(ctx context.Context, produkID, gambarID, direction string) (map[string]interface{}, error)
	SetPrimary(ctx context.Context, produkID, gambarID string) (*models.ProdukGambarResponse, error)
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

func (s *produkGambarService) CreateWithFile(ctx context.Context, produkID string, file *multipart.FileHeader, isPrimary bool) (*models.ProdukGambarResponse, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, errors.New("produk_id tidak valid")
	}

	// Upload file
	produkDir := fmt.Sprintf("products/%s", produkID)
	relativePath, err := utils.SaveUploadedFile(file, produkDir, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("gagal upload gambar: %w", err)
	}

	// Auto-increment urutan per produk
	maxUrutan, err := s.repo.GetMaxUrutanByProdukID(ctx, produkID)
	if err != nil {
		utils.DeleteFile(relativePath, s.cfg)
		return nil, err
	}

	gambar := &models.ProdukGambar{
		ProdukID:  produkUUID,
		GambarURL: relativePath,
		Urutan:    maxUrutan + 1,
		IsPrimary: isPrimary,
	}

	if err := s.repo.Create(ctx, gambar); err != nil {
		utils.DeleteFile(relativePath, s.cfg)
		return nil, err
	}

	// If this is primary, update others
	if isPrimary {
		s.repo.SetPrimary(ctx, produkID, gambar.ID.String())
	}

	return &models.ProdukGambarResponse{
		ID:        gambar.ID.String(),
		GambarURL: utils.GetFileURL(gambar.GambarURL, s.cfg),
		Urutan:    gambar.Urutan,
		IsPrimary: gambar.IsPrimary,
	}, nil
}

func (s *produkGambarService) Delete(ctx context.Context, produkID, gambarID string) error {
	gambar, err := s.repo.FindByID(ctx, gambarID)
	if err != nil {
		return errors.New("gambar tidak ditemukan")
	}

	// Verify produk ownership
	if gambar.ProdukID.String() != produkID {
		return errors.New("gambar tidak ditemukan")
	}

	// Check if this is the last image
	count, _ := s.repo.CountByProdukID(ctx, produkID)
	if count <= 1 {
		return errors.New("tidak dapat menghapus gambar terakhir. Produk harus memiliki minimal 1 gambar")
	}

	// Delete file from storage
	if err := utils.DeleteFile(gambar.GambarURL, s.cfg); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to delete file %s: %v\n", gambar.GambarURL, err)
	}

	// Delete from database
	return s.repo.Delete(ctx, gambarID)
}

func (s *produkGambarService) Reorder(ctx context.Context, produkID, gambarID, direction string) (map[string]interface{}, error) {
	gambar, err := s.repo.FindByID(ctx, gambarID)
	if err != nil {
		return nil, errors.New("gambar tidak ditemukan")
	}

	// Verify produk ownership
	if gambar.ProdukID.String() != produkID {
		return nil, errors.New("gambar tidak ditemukan")
	}

	// Get all gambar for this produk
	gambars, err := s.repo.FindByProdukID(ctx, produkID)
	if err != nil {
		return nil, err
	}

	currentUrutan := gambar.Urutan
	var swapUrutan int
	var swappedGambar *models.ProdukGambar

	if direction == "up" {
		// Find previous item
		for i := range gambars {
			if gambars[i].Urutan < currentUrutan {
				if swappedGambar == nil || gambars[i].Urutan > swappedGambar.Urutan {
					swappedGambar = &gambars[i]
					swapUrutan = gambars[i].Urutan
				}
			}
		}
	} else if direction == "down" {
		// Find next item
		for i := range gambars {
			if gambars[i].Urutan > currentUrutan {
				if swappedGambar == nil || gambars[i].Urutan < swappedGambar.Urutan {
					swappedGambar = &gambars[i]
					swapUrutan = gambars[i].Urutan
				}
			}
		}
	} else {
		return nil, errors.New("direction harus 'up' atau 'down'")
	}

	if swappedGambar == nil {
		return nil, errors.New("tidak dapat memindahkan gambar ke arah tersebut")
	}

	// Swap urutan
	gambar.Urutan = swapUrutan
	swappedGambar.Urutan = currentUrutan

	if err := s.repo.Update(ctx, gambar); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, swappedGambar); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"item": map[string]interface{}{
			"id":     gambar.ID.String(),
			"urutan": gambar.Urutan,
		},
		"swapped_with": map[string]interface{}{
			"id":     swappedGambar.ID.String(),
			"urutan": swappedGambar.Urutan,
		},
	}, nil
}

func (s *produkGambarService) SetPrimary(ctx context.Context, produkID, gambarID string) (*models.ProdukGambarResponse, error) {
	gambar, err := s.repo.FindByID(ctx, gambarID)
	if err != nil {
		return nil, errors.New("gambar tidak ditemukan")
	}

	// Verify produk ownership
	if gambar.ProdukID.String() != produkID {
		return nil, errors.New("gambar tidak ditemukan")
	}

	// Set as primary
	if err := s.repo.SetPrimary(ctx, produkID, gambarID); err != nil {
		return nil, err
	}

	gambar.IsPrimary = true

	return &models.ProdukGambarResponse{
		ID:        gambar.ID.String(),
		GambarURL: utils.GetFileURL(gambar.GambarURL, s.cfg),
		Urutan:    gambar.Urutan,
		IsPrimary: gambar.IsPrimary,
	}, nil
}
