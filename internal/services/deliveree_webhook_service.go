package services

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/repositories"

	"gorm.io/gorm"
)

// DelivereeWebhookService menangani event webhook dari Deliveree dan
// memperbarui status pesanan yang bersangkutan.
//
// Logika pemrosesan dipusatkan di ProviderWebhookHandler (lihat
// webhook_provider_service.go) karena Deliveree & Forwarder berbagi platform
// on-demand yang sama — format status identik, hanya endpoint dan kolom
// identifikasi pesanan yang berbeda.
type DelivereeWebhookService interface {
	// Handle memproses payload webhook dan mengembalikan true jika ada pesanan
	// yang cocok & diperbarui.
	Handle(ctx context.Context, req *dto.DelivereeWebhookRequest) (bool, error)
}

type delivereeWebhookService struct {
	handler *ProviderWebhookHandler
}

func NewDelivereeWebhookService(pesananRepo repositories.PesananRepository, db *gorm.DB) DelivereeWebhookService {
	return &delivereeWebhookService{
		handler: NewProviderWebhookHandler(pesananRepo, db, delivereeWebhookConfig()),
	}
}

func (s *delivereeWebhookService) Handle(ctx context.Context, req *dto.DelivereeWebhookRequest) (bool, error) {
	// Cari pesanan berdasarkan deliveree_booking_id
	identifier := string(req.ID)
	if identifier == "" {
		identifier = string(req.NoBooking)
	}
	return s.handler.Handle(ctx, identifier, req.Status, req.TrackingURL)
}
