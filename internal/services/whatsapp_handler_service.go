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
	FindAll(ctx context.Context) ([]models.WhatsAppHandlerResponse, error)
	Update(ctx context.Context, id string, req *models.UpdateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error)
	GetActive(ctx context.Context) (*models.WhatsAppHandlerPublicResponse, error)
}

type whatsAppHandlerService struct {
	repo repositories.WhatsAppHandlerRepository
}

func NewWhatsAppHandlerService(repo repositories.WhatsAppHandlerRepository) WhatsAppHandlerService {
	return &whatsAppHandlerService{repo: repo}
}

func (s *whatsAppHandlerService) FindAll(ctx context.Context) ([]models.WhatsAppHandlerResponse, error) {
	items, err := s.repo.FindAllSimple(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data WhatsApp handler")
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

	return responses, nil
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
		// If setting to active, deactivate all others first
		if *req.IsActive {
			if err := s.repo.DeactivateAll(ctx); err != nil {
				return nil, errors.New("gagal menonaktifkan handler lain")
			}
		}
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
