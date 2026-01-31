package services

import (
	"context"
	"errors"
	"strings"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranService interface {
	GetAll(ctx context.Context, groupID *string, isActive *bool) ([]models.MetodePembayaranListResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateMetodePembayaranRequest) (*models.MetodePembayaranDetailResponse, error)
	GetAllGrouped(ctx context.Context, isAdmin bool) ([]dto.PaymentMethodGroupResponse, error)
	ToggleMethodStatus(ctx context.Context, id string) (*dto.ToggleMethodStatusResponse, error)
	ToggleGroupStatus(ctx context.Context, urutan int) (*dto.ToggleGroupStatusResponse, error)
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

// GetAllGrouped returns payment methods grouped by group
func (s *metodePembayaranService) GetAllGrouped(ctx context.Context, isAdmin bool) ([]dto.PaymentMethodGroupResponse, error) {
	// 1. Get all groups ordered by urutan
	groups, err := s.groupRepo.FindAllSimple(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.PaymentMethodGroupResponse

	for _, group := range groups {
		// Skip inactive groups for public
		if !isAdmin && !group.IsActive {
			continue
		}

		// 2. Get methods for each group
		methods, err := s.repo.FindAllSimple(ctx, &group.ID, nil)
		if err != nil {
			return nil, err
		}

		var methodResponses []dto.PaymentMethodResponse
		for _, method := range methods {
			// Skip inactive methods for public
			if !isAdmin && !method.IsActive {
				continue
			}

			resp := dto.PaymentMethodResponse{
				ID:        method.ID.String(),
				Nama:      method.Nama,
				Kode:      method.Kode,
				LogoValue: method.LogoValue,
			}

			// Include admin-only fields
			if isAdmin {
				resp.Urutan = method.Urutan
				resp.IsActive = method.IsActive
			}

			methodResponses = append(methodResponses, resp)
		}

		// Skip group if no active methods (for public)
		if !isAdmin && len(methodResponses) == 0 {
			continue
		}

		groupResp := dto.PaymentMethodGroupResponse{
			Group:   group.Nama,
			Methods: methodResponses,
		}

		// Include admin-only fields
		if isAdmin {
			groupResp.Urutan = group.Urutan
			groupResp.IsActive = group.IsActive
		}

		result = append(result, groupResp)
	}

	return result, nil
}

// ToggleMethodStatus toggles the is_active status of a payment method
func (s *metodePembayaranService) ToggleMethodStatus(ctx context.Context, id string) (*dto.ToggleMethodStatusResponse, error) {
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

	return &dto.ToggleMethodStatusResponse{
		ID:       metode.ID.String(),
		Nama:     metode.Nama,
		Kode:     metode.Kode,
		IsActive: metode.IsActive,
	}, nil
}

// ToggleGroupStatus toggles the is_active status of a payment group by urutan
func (s *metodePembayaranService) ToggleGroupStatus(ctx context.Context, urutan int) (*dto.ToggleGroupStatusResponse, error) {
	// Find group by urutan
	groups, err := s.groupRepo.FindAllSimple(ctx)
	if err != nil {
		return nil, err
	}

	var targetGroup *models.MetodePembayaranGroup
	for i := range groups {
		if groups[i].Urutan == urutan {
			targetGroup = &groups[i]
			break
		}
	}

	if targetGroup == nil {
		return nil, errors.New("Group tidak ditemukan")
	}

	// Toggle status
	group, err := s.groupRepo.ToggleStatus(ctx, targetGroup.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ToggleGroupStatusResponse{
		Group:    group.Nama,
		Urutan:   group.Urutan,
		IsActive: group.IsActive,
	}, nil
}
