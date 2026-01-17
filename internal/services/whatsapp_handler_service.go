package services

import (
	"context"
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WhatsAppHandlerService interface {
	Create(ctx context.Context, req *models.CreateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error)
	FindAll(ctx context.Context, params *models.WhatsAppHandlerFilterRequest) ([]models.WhatsAppHandlerResponse, *models.PaginationMeta, error)
	FindByID(ctx context.Context, id string) (*models.WhatsAppHandlerResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
	GetActive(ctx context.Context) (*models.WhatsAppHandlerPublicResponse, error)
}

type whatsAppHandlerService struct {
	repo repositories.WhatsAppHandlerRepository
}

func NewWhatsAppHandlerService(repo repositories.WhatsAppHandlerRepository) WhatsAppHandlerService {
	return &whatsAppHandlerService{repo: repo}
}

func (s *whatsAppHandlerService) Create(ctx context.Context, req *models.CreateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error) {
	handler := &models.WhatsAppHandler{
		NomorWA:   req.NomorWA,
		PesanAwal: req.PesanAwal,
		IsActive:  req.IsActive,
	}

	if err := s.repo.Create(ctx, handler); err != nil {
		return nil, errors.New("gagal membuat WhatsApp handler")
	}

	return &models.WhatsAppHandlerResponse{
		ID:          handler.ID.String(),
		NomorWA:     handler.NomorWA,
		PesanAwal:   handler.PesanAwal,
		IsActive:    handler.IsActive,
		WhatsAppURL: handler.GetWhatsAppURL(),
		CreatedAt:   handler.CreatedAt,
		UpdatedAt:   handler.UpdatedAt,
	}, nil
}

func (s *whatsAppHandlerService) FindAll(ctx context.Context, params *models.WhatsAppHandlerFilterRequest) ([]models.WhatsAppHandlerResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	items, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, errors.New("gagal mengambil data WhatsApp handler")
	}

	var responses []models.WhatsAppHandlerResponse
	for _, item := range items {
		responses = append(responses, models.WhatsAppHandlerResponse{
			ID:          item.ID.String(),
			NomorWA:     item.NomorWA,
			PesanAwal:   item.PesanAwal,
			IsActive:    item.IsActive,
			WhatsAppURL: item.GetWhatsAppURL(),
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return responses, &meta, nil
}

func (s *whatsAppHandlerService) FindByID(ctx context.Context, id string) (*models.WhatsAppHandlerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	handler, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("WhatsApp handler tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data WhatsApp handler")
	}

	return &models.WhatsAppHandlerResponse{
		ID:          handler.ID.String(),
		NomorWA:     handler.NomorWA,
		PesanAwal:   handler.PesanAwal,
		IsActive:    handler.IsActive,
		WhatsAppURL: handler.GetWhatsAppURL(),
		CreatedAt:   handler.CreatedAt,
		UpdatedAt:   handler.UpdatedAt,
	}, nil
}

func (s *whatsAppHandlerService) Update(ctx context.Context, id string, req *models.UpdateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	handler, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("WhatsApp handler tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data WhatsApp handler")
	}

	if req.NomorWA != nil {
		handler.NomorWA = *req.NomorWA
	}
	if req.PesanAwal != nil {
		handler.PesanAwal = *req.PesanAwal
	}
	if req.IsActive != nil {
		handler.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, handler); err != nil {
		return nil, errors.New("gagal mengupdate WhatsApp handler")
	}

	return &models.WhatsAppHandlerResponse{
		ID:          handler.ID.String(),
		NomorWA:     handler.NomorWA,
		PesanAwal:   handler.PesanAwal,
		IsActive:    handler.IsActive,
		WhatsAppURL: handler.GetWhatsAppURL(),
		CreatedAt:   handler.CreatedAt,
		UpdatedAt:   handler.UpdatedAt,
	}, nil
}

func (s *whatsAppHandlerService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	// Check if exists
	_, err = s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("WhatsApp handler tidak ditemukan")
		}
		return errors.New("gagal mengambil data WhatsApp handler")
	}

	if err := s.repo.Delete(ctx, uid); err != nil {
		return errors.New("gagal menghapus WhatsApp handler")
	}

	return nil
}

func (s *whatsAppHandlerService) SetActive(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	// Check if exists
	_, err = s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("WhatsApp handler tidak ditemukan")
		}
		return errors.New("gagal mengambil data WhatsApp handler")
	}

	// Set active
	if err := s.repo.SetActive(ctx, uid); err != nil {
		return errors.New("gagal mengaktifkan WhatsApp handler")
	}

	return nil
}

func (s *whatsAppHandlerService) GetActive(ctx context.Context) (*models.WhatsAppHandlerPublicResponse, error) {
	handler, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil WhatsApp handler aktif")
	}

	if handler == nil {
		return nil, nil
	}

	return &models.WhatsAppHandlerPublicResponse{
		NomorWA:     handler.NomorWA,
		PesanAwal:   handler.PesanAwal,
		WhatsAppURL: handler.GetWhatsAppURL(),
	}, nil
}
