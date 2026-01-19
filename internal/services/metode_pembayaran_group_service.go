package services

import (
	"context"
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranGroupService interface {
	GetAll(ctx context.Context) ([]models.MetodePembayaranGroupListResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error)
}

type metodePembayaranGroupService struct {
	repo repositories.MetodePembayaranGroupRepository
}

func NewMetodePembayaranGroupService(repo repositories.MetodePembayaranGroupRepository) MetodePembayaranGroupService {
	return &metodePembayaranGroupService{repo: repo}
}

func (s *metodePembayaranGroupService) GetAll(ctx context.Context) ([]models.MetodePembayaranGroupListResponse, error) {
	groups, err := s.repo.FindAllSimple(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.MetodePembayaranGroupListResponse, len(groups))
	for i, group := range groups {
		// Count metode for each group
		jumlahMetode, _ := s.repo.CountActiveMetode(ctx, group.ID)

		responses[i] = models.MetodePembayaranGroupListResponse{
			ID:           group.ID.String(),
			Nama:         group.Nama,
			Urutan:       group.Urutan,
			IsActive:     group.IsActive,
			JumlahMetode: int(jumlahMetode),
			UpdatedAt:    group.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *metodePembayaranGroupService) Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranGroupRequest) (*models.MetodePembayaranGroupResponse, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID group tidak valid")
	}

	group, err := s.repo.FindByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Group metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	// Check if nama already exists (exclude current group)
	if req.Nama != nil {
		exists, err := s.repo.CheckByName(ctx, *req.Nama, &groupID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("Nama group sudah digunakan")
		}
		group.Nama = *req.Nama
	}

	if req.Urutan != nil {
		group.Urutan = *req.Urutan
	}

	if req.IsActive != nil {
		group.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, group); err != nil {
		return nil, err
	}

	return &models.MetodePembayaranGroupResponse{
		ID:        group.ID.String(),
		Nama:      group.Nama,
		Urutan:    group.Urutan,
		IsActive:  group.IsActive,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}, nil
}
