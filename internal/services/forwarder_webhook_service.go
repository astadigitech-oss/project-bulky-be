package services

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/repositories"

	"gorm.io/gorm"
)

// ForwarderWebhookService menangani event webhook dari Forwarder dan
// memperbarui status pesanan yang bersangkutan.
//
// Logika pemrosesan dipusatkan di ProviderWebhookHandler (lihat
// webhook_provider_service.go) karena Deliveree & Forwarder berbagi platform
// on-demand yang sama — format status identik, hanya endpoint dan kolom
// identifikasi pesanan yang berbeda.
type ForwarderWebhookService interface {
	// Handle memproses payload webhook dan mengembalikan true jika ada pesanan
	// yang cocok & diperbarui.
	Handle(ctx context.Context, req *dto.ForwarderWebhookRequest) (bool, error)
}

type forwarderWebhookService struct {
	handler *ProviderWebhookHandler
}

func NewForwarderWebhookService(pesananRepo repositories.PesananRepository, db *gorm.DB) ForwarderWebhookService {
	return &forwarderWebhookService{
		handler: NewProviderWebhookHandler(pesananRepo, db, forwarderWebhookConfig()),
	}
}

func (s *forwarderWebhookService) Handle(ctx context.Context, req *dto.ForwarderWebhookRequest) (bool, error) {
	// Cari pesanan berdasarkan forwarder_tracking_no
	identifier := string(req.ID)
	if identifier == "" {
		identifier = string(req.NoBooking)
	}
	return s.handler.Handle(ctx, identifier, req.Status, req.TrackingURL)
}
