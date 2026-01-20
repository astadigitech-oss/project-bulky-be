package services

import (
	"context"
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type WhatsAppHandlerService interface {
	Get(ctx context.Context) (*models.WhatsAppHandlerSimpleResponse, error)
	Update(ctx context.Context, req *models.UpdateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error)
	GetActive(ctx context.Context) (*models.WhatsAppHandlerPublicResponse, error)
}

type whatsAppHandlerService struct {
	repo repositories.WhatsAppHandlerRepository
}

func NewWhatsAppHandlerService(repo repositories.WhatsAppHandlerRepository) WhatsAppHandlerService {
	return &whatsAppHandlerService{repo: repo}
}

func (s *whatsAppHandlerService) Get(ctx context.Context) (*models.WhatsAppHandlerSimpleResponse, error) {
	// Get the active handler or first one
	handler, err := s.repo.GetActive(ctx)
	if err != nil {
		// If no active, try to get first one
		items, err := s.repo.FindAllSimple(ctx)
		if err != nil || len(items) == 0 {
			return nil, errors.New("WhatsApp handler tidak ditemukan")
		}
		handler = &items[0]
	}

	return &models.WhatsAppHandlerSimpleResponse{
		ID:        handler.ID.String(),
		NomorWA:   handler.NomorWA,
		PesanAwal: handler.PesanAwal,
		// IsActive:    handler.IsActive,
		WhatsAppURL: handler.GetWhatsAppURL(),
		// CreatedAt:   handler.CreatedAt,
		UpdatedAt: handler.UpdatedAt,
	}, nil
}

func (s *whatsAppHandlerService) Update(ctx context.Context, req *models.UpdateWhatsAppHandlerRequest) (*models.WhatsAppHandlerResponse, error) {
	// Get the active handler or first one
	handler, err := s.repo.GetActive(ctx)
	if err != nil {
		// If no active, try to get first one
		items, err := s.repo.FindAllSimple(ctx)
		if err != nil || len(items) == 0 {
			return nil, errors.New("WhatsApp handler tidak ditemukan")
		}
		handler = &items[0]
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
