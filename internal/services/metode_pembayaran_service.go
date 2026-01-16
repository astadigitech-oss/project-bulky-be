package services

import (
	"context"
	"errors"
	"strings"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranService interface {
	GetAll(ctx context.Context, params *models.PaginationRequest, groupID *string, isActive *bool) ([]models.MetodePembayaranListResponse, *models.PaginationMeta, error)
	GetByID(ctx context.Context, id string) (*models.MetodePembayaranDetailResponse, error)
	Create(ctx context.Context, req *models.CreateMetodePembayaranRequest, logoURL *string) (*models.MetodePembayaranDetailResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranRequest, logoURL *string) (*models.MetodePembayaranDetailResponse, error)
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) (*models.MetodePembayaranResponse, error)
}

type metodePembayaranService struct {
	repo      repositories.MetodePembayaranRepository
	groupRepo repositories.MetodePembayaranGroupRepository
	cfg       *config.Config
}

func NewMetodePembayaranService(repo repositories.MetodePembayaranRepository, groupRepo repositories.MetodePembayaranGroupRepository, cfg *config.Config) MetodePembayaranService {
	return &metodePembayaranService{
		repo:      repo,
		groupRepo: groupRepo,
		cfg:       cfg,
	}
}

func (s *metodePembayaranService) GetAll(ctx context.Context, params *models.PaginationRequest, groupID *string, isActive *bool) ([]models.MetodePembayaranListResponse, *models.PaginationMeta, error) {
	var groupUUID *uuid.UUID
	if groupID != nil && *groupID != "" {
		parsed, err := uuid.Parse(*groupID)
		if err != nil {
			return nil, nil, errors.New("ID group tidak valid")
		}
		groupUUID = &parsed
	}

	metodes, total, err := s.repo.FindAll(ctx, params, groupUUID, isActive)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]models.MetodePembayaranListResponse, len(metodes))
	for i, metode := range metodes {
		responses[i] = models.MetodePembayaranListResponse{
			ID:       metode.ID.String(),
			Nama:     metode.Nama,
			Kode:     metode.Kode,
			Logo:     utils.GetFileURLPtr(metode.Logo, s.cfg),
			Urutan:   metode.Urutan,
			IsActive: metode.IsActive,
			Group: models.MetodePembayaranGroupSimple{
				ID:   metode.Group.ID.String(),
				Nama: metode.Group.Nama,
			},
			UpdatedAt: metode.UpdatedAt,
		}
	}

	meta := &models.PaginationMeta{
		FirstPage:   1,
		LastPage:    int((total + int64(params.PerPage) - 1) / int64(params.PerPage)),
		CurrentPage: params.Page,
		From:        (params.Page-1)*params.PerPage + 1,
		Last:        len(metodes),
		Total:       total,
		PerPage:     params.PerPage,
	}

	return responses, meta, nil
}

func (s *metodePembayaranService) GetByID(ctx context.Context, id string) (*models.MetodePembayaranDetailResponse, error) {
	metodeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID metode pembayaran tidak valid")
	}

	metode, err := s.repo.FindByIDWithGroup(ctx, metodeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	return &models.MetodePembayaranDetailResponse{
		ID:       metode.ID.String(),
		Nama:     metode.Nama,
		Kode:     metode.Kode,
		Logo:     utils.GetFileURLPtr(metode.Logo, s.cfg),
		Urutan:   metode.Urutan,
		IsActive: metode.IsActive,
		Group: models.MetodePembayaranGroupSimple{
			ID:   metode.Group.ID.String(),
			Nama: metode.Group.Nama,
		},
		CreatedAt: metode.CreatedAt,
		UpdatedAt: metode.UpdatedAt,
	}, nil
}

func (s *metodePembayaranService) Create(ctx context.Context, req *models.CreateMetodePembayaranRequest, logoURL *string) (*models.MetodePembayaranDetailResponse, error) {
	// Parse and validate group_id
	groupID, err := uuid.Parse(req.GroupID)
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

	// Check if kode already exists
	exists, err := s.repo.CheckByKode(ctx, strings.ToUpper(req.Kode), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Kode sudah digunakan")
	}

	metode := &models.MetodePembayaran{
		GroupID:  groupID,
		Nama:     req.Nama,
		Kode:     strings.ToUpper(req.Kode),
		Logo:     logoURL,
		Urutan:   req.Urutan,
		IsActive: req.IsActive,
	}

	if err := s.repo.Create(ctx, metode); err != nil {
		return nil, err
	}

	// Reload with group
	return s.GetByID(ctx, metode.ID.String())
}

func (s *metodePembayaranService) Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranRequest, logoURL *string) (*models.MetodePembayaranDetailResponse, error) {
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

	// Update logo if provided
	if logoURL != nil {
		metode.Logo = logoURL
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
	return s.GetByID(ctx, metode.ID.String())
}

func (s *metodePembayaranService) Delete(ctx context.Context, id string) error {
	metodeID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID metode pembayaran tidak valid")
	}

	metode, err := s.repo.FindByID(ctx, metodeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Metode pembayaran tidak ditemukan")
		}
		return err
	}

	// Check if metode is used in transactions
	used, err := s.repo.CheckUsedInTransaction(ctx, metode.ID)
	if err != nil {
		return err
	}
	if used {
		return errors.New("Tidak dapat menghapus metode pembayaran yang sudah digunakan dalam transaksi")
	}

	return s.repo.Delete(ctx, metodeID)
}

func (s *metodePembayaranService) ToggleStatus(ctx context.Context, id string) (*models.MetodePembayaranResponse, error) {
	metodeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID metode pembayaran tidak valid")
	}

	metode, err := s.repo.ToggleStatus(ctx, metodeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Metode pembayaran tidak ditemukan")
		}
		return nil, err
	}

	return &models.MetodePembayaranResponse{
		ID:       metode.ID.String(),
		Nama:     metode.Nama,
		Kode:     metode.Kode,
		IsActive: metode.IsActive,
	}, nil
}
