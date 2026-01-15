package services

import (
	"context"
	"errors"
	"strings"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PPNService interface {
	GetAll(ctx context.Context, params *models.PaginationRequest) ([]models.PPNResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*models.PPNResponse, error)
	GetActive(ctx context.Context) (*models.PPNResponse, error)
	Create(ctx context.Context, req *models.CreatePPNRequest) (*models.PPNResponse, error)
	Update(ctx context.Context, id string, req *models.UpdatePPNRequest) (*models.PPNResponse, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) (*models.PPNResponse, error)
}

type ppnService struct {
	repo repositories.PPNRepository
}

func NewPPNService(repo repositories.PPNRepository) PPNService {
	return &ppnService{repo: repo}
}

func (s *ppnService) GetAll(ctx context.Context, params *models.PaginationRequest) ([]models.PPNResponse, *models.PaginationMeta, error) {
	ppns, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]models.PPNResponse, len(ppns))
	for i, ppn := range ppns {
		responses[i] = models.PPNResponse{
			ID:         ppn.ID.String(),
			Persentase: ppn.Persentase.InexactFloat64(),
			IsActive:   ppn.IsActive,
			CreatedAt:  ppn.CreatedAt,
			UpdatedAt:  ppn.UpdatedAt,
		}
	}

	meta := &models.PaginationMeta{
		FirstPage:   1,
		LastPage:    int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
		CurrentPage: params.Page,
		From:        (params.Page-1)*params.PerPage + 1,
		Last:        len(ppns),
		Total:       total,
		PerPage:     params.PerPage,
	}

	return responses, meta, nil
}

func (s *ppnService) GetByID(ctx context.Context, id string) (*models.PPNResponse, error) {
	ppnID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID PPN tidak valid")
	}

	ppn, err := s.repo.FindByID(ctx, ppnID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PPN tidak ditemukan")
		}
		return nil, err
	}

	return &models.PPNResponse{
		ID:         ppn.ID.String(),
		Persentase: ppn.Persentase.InexactFloat64(),
		IsActive:   ppn.IsActive,
		CreatedAt:  ppn.CreatedAt,
		UpdatedAt:  ppn.UpdatedAt,
	}, nil
}

func (s *ppnService) GetActive(ctx context.Context) (*models.PPNResponse, error) {
	ppn, err := s.repo.FindActive(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PPN aktif tidak ditemukan")
		}
		return nil, err
	}

	return &models.PPNResponse{
		ID:         ppn.ID.String(),
		Persentase: ppn.Persentase.InexactFloat64(),
		IsActive:   ppn.IsActive,
		CreatedAt:  ppn.CreatedAt,
		UpdatedAt:  ppn.UpdatedAt,
	}, nil
}

func (s *ppnService) Create(ctx context.Context, req *models.CreatePPNRequest) (*models.PPNResponse, error) {
	ppn := &models.PPN{
		Persentase: decimal.NewFromFloat(req.Persentase),
		IsActive:   req.IsActive,
	}

	if err := s.repo.Create(ctx, ppn); err != nil {
		// Check for duplicate persentase error
		if strings.Contains(err.Error(), "unique_ppn_persentase") || strings.Contains(err.Error(), "duplicate key value") {
			return nil, errors.New("Persentase PPN sudah ada, silakan gunakan nilai yang berbeda")
		}
		return nil, err
	}

	return &models.PPNResponse{
		ID:         ppn.ID.String(),
		Persentase: ppn.Persentase.InexactFloat64(),
		IsActive:   ppn.IsActive,
		CreatedAt:  ppn.CreatedAt,
		UpdatedAt:  ppn.UpdatedAt,
	}, nil
}

func (s *ppnService) Update(ctx context.Context, id string, req *models.UpdatePPNRequest) (*models.PPNResponse, error) {
	ppnID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID PPN tidak valid")
	}

	ppn, err := s.repo.FindByID(ctx, ppnID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PPN tidak ditemukan")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Persentase != nil {
		ppn.Persentase = decimal.NewFromFloat(*req.Persentase)
	}
	if req.IsActive != nil {
		ppn.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, ppn); err != nil {
		// Check for duplicate persentase error
		if strings.Contains(err.Error(), "unique_ppn_persentase") || strings.Contains(err.Error(), "duplicate key value") {
			return nil, errors.New("Persentase PPN sudah ada, silakan gunakan nilai yang berbeda")
		}
		return nil, err
	}

	return &models.PPNResponse{
		ID:         ppn.ID.String(),
		Persentase: ppn.Persentase.InexactFloat64(),
		IsActive:   ppn.IsActive,
		CreatedAt:  ppn.CreatedAt,
		UpdatedAt:  ppn.UpdatedAt,
	}, nil
}

func (s *ppnService) Delete(ctx context.Context, id string) error {
	ppnID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID PPN tidak valid")
	}

	ppn, err := s.repo.FindByID(ctx, ppnID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("PPN tidak ditemukan")
		}
		return err
	}

	// Business Rule: Cannot delete active PPN
	if ppn.IsActive {
		return errors.New("Tidak dapat menghapus PPN yang sedang aktif")
	}

	return s.repo.Delete(ctx, ppnID)
}

func (s *ppnService) SetActive(ctx context.Context, id string) (*models.PPNResponse, error) {
	ppnID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID PPN tidak valid")
	}

	// Check if PPN exists
	ppn, err := s.repo.FindByID(ctx, ppnID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PPN tidak ditemukan")
		}
		return nil, err
	}

	// Set as active (trigger will handle deactivating others)
	if err := s.repo.SetActive(ctx, ppnID); err != nil {
		return nil, err
	}

	// Reload to get updated data
	ppn, err = s.repo.FindByID(ctx, ppnID)
	if err != nil {
		return nil, err
	}

	return &models.PPNResponse{
		ID:         ppn.ID.String(),
		Persentase: ppn.Persentase.InexactFloat64(),
		IsActive:   ppn.IsActive,
		CreatedAt:  ppn.CreatedAt,
		UpdatedAt:  ppn.UpdatedAt,
	}, nil
}
