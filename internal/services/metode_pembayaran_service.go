package services

import (
	"context"
	"errors"
	"strings"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranService interface {
	GetAll(ctx context.Context, groupID *string, isActive *bool) ([]models.MetodePembayaranListResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranRequest) (*models.MetodePembayaranDetailResponse, error)
}

type metodePembayaranService struct {
	repo      repositories.MetodePembayaranRepository
	groupRepo repositories.MetodePembayaranGroupRepository
}

func NewMetodePembayaranService(repo repositories.MetodePembayaranRepository, groupRepo repositories.MetodePembayaranGroupRepository) MetodePembayaranService {
	return &metodePembayaranService{
		repo:      repo,
		groupRepo: groupRepo,
	}
}

func (s *metodePembayaranService) GetAll(ctx context.Context, groupID *string, isActive *bool) ([]models.MetodePembayaranListResponse, error) {
	var groupUUID *uuid.UUID
	if groupID != nil && *groupID != "" {
		parsed, err := uuid.Parse(*groupID)
		if err != nil {
			return nil, errors.New("ID group tidak valid")
		}
		groupUUID = &parsed
	}

	metodes, err := s.repo.FindAllSimple(ctx, groupUUID, isActive)
	if err != nil {
		return nil, err
	}

	responses := make([]models.MetodePembayaranListResponse, len(metodes))
	for i, metode := range metodes {
		responses[i] = models.MetodePembayaranListResponse{
			ID:        metode.ID.String(),
			Nama:      metode.Nama,
			Kode:      metode.Kode,
			LogoValue: metode.LogoValue,
			Urutan:    metode.Urutan,
			IsActive:  metode.IsActive,
			Group: models.MetodePembayaranGroupSimple{
				ID:   metode.Group.ID.String(),
				Nama: metode.Group.Nama,
			},
			UpdatedAt: metode.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *metodePembayaranService) Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranRequest) (*models.MetodePembayaranDetailResponse, error) {
	metodeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID metode pembayaran tidak valid")
	}

	metode, err := s.repo.FindByID(ctx, metodeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	// Update group_id if provided
	if req.GroupID != nil {
		groupID, err := uuid.Parse(*req.GroupID)
		if err != nil {
			return nil, errors.New("ID group tidak valid")
		}

		// Check if group exists
		_, err = s.groupRepo.FindByID(ctx, groupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Group tidak ditemukan")
			}
			return nil, err
		}

		metode.GroupID = groupID
	}

	// Update nama if provided
	if req.Nama != nil {
		metode.Nama = *req.Nama
	}

	// Update kode if provided
	if req.Kode != nil {
		kodeUpper := strings.ToUpper(*req.Kode)
		exists, err := s.repo.CheckByKode(ctx, kodeUpper, &metodeID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("Kode sudah digunakan")
		}
		metode.Kode = kodeUpper
	}

	// Update logo_value if provided
	if req.LogoValue != nil {
		metode.LogoValue = req.LogoValue
	}

	// Update urutan if provided
	if req.Urutan != nil {
		metode.Urutan = *req.Urutan
	}

	// Update is_active if provided
	if req.IsActive != nil {
		metode.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, metode); err != nil {
		return nil, err
	}

	// Reload with group
	metode, err = s.repo.FindByIDWithGroup(ctx, metodeID)
	if err != nil {
		return nil, err
	}

	return &models.MetodePembayaranDetailResponse{
		ID:        metode.ID.String(),
		Nama:      metode.Nama,
		Kode:      metode.Kode,
		LogoValue: metode.LogoValue,
		Urutan:    metode.Urutan,
		IsActive:  metode.IsActive,
		Group: models.MetodePembayaranGroupSimple{
			ID:   metode.Group.ID.String(),
			Nama: metode.Group.Nama,
		},
		CreatedAt: metode.CreatedAt,
		UpdatedAt: metode.UpdatedAt,
	}, nil
}
