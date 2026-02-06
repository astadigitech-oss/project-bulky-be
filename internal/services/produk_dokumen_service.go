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

type ProdukDokumenService interface {
	CreateWithFile(ctx context.Context, produkID string, file *multipart.FileHeader, namaDokumen string) (*models.ProdukDokumenResponse, error)
	Delete(ctx context.Context, produkID, dokumenID string) error
}

type produkDokumenService struct {
	repo repositories.ProdukDokumenRepository
	cfg  *config.Config
}

func NewProdukDokumenService(repo repositories.ProdukDokumenRepository, cfg *config.Config) ProdukDokumenService {
	return &produkDokumenService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *produkDokumenService) CreateWithFile(ctx context.Context, produkID string, file *multipart.FileHeader, namaDokumen string) (*models.ProdukDokumenResponse, error) {
	produkUUID, err := uuid.Parse(produkID)
	if err != nil {
		return nil, errors.New("produk_id tidak valid")
	}

	// Upload file
	dokumenDir := fmt.Sprintf("documents/%s", produkID)
	relativePath, err := utils.SaveUploadedFile(file, dokumenDir, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("gagal upload dokumen: %w", err)
	}

	// Use provided name or fallback to filename
	if namaDokumen == "" {
		namaDokumen = file.Filename
	}

	ukuranFile := int(file.Size)
	dokumen := &models.ProdukDokumen{
		ProdukID:    produkUUID,
		NamaDokumen: namaDokumen,
		FileURL:     relativePath,
		TipeFile:    "pdf",
		UkuranFile:  &ukuranFile,
	}

	if err := s.repo.Create(ctx, dokumen); err != nil {
		utils.DeleteFile(relativePath, s.cfg)
		return nil, err
	}

	return &models.ProdukDokumenResponse{
		ID:          dokumen.ID.String(),
		NamaDokumen: dokumen.NamaDokumen,
		FileURL:     utils.GetFileURL(dokumen.FileURL, s.cfg),
		TipeFile:    dokumen.TipeFile,
		UkuranFile:  dokumen.UkuranFile,
	}, nil
}

func (s *produkDokumenService) Delete(ctx context.Context, produkID, dokumenID string) error {
	dokumen, err := s.repo.FindByID(ctx, dokumenID)
	if err != nil {
		return errors.New("dokumen tidak ditemukan")
	}

	// Verify produk ownership
	if dokumen.ProdukID.String() != produkID {
		return errors.New("dokumen tidak ditemukan")
	}

	// Delete file from storage
	if err := utils.DeleteFile(dokumen.FileURL, s.cfg); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to delete file %s: %v\n", dokumen.FileURL, err)
	}

	// Delete from database
	return s.repo.Delete(ctx, dokumenID)
}
