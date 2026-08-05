package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"gorm.io/gorm"
)

// DelivereeWebhookStatus adalah status pengiriman yang dikirim Deliveree via webhook.
type DelivereeWebhookStatus string

const (
	// DelivereeWebhookStatusDeliveryCompleted — paket sudah diterima penerima.
	DelivereeWebhookStatusDeliveryCompleted DelivereeWebhookStatus = "delivery_completed"
	// DelivereeWebhookStatusLocatingDriver — sedang mencari driver.
	DelivereeWebhookStatusLocatingDriver DelivereeWebhookStatus = "locating_driver"
	// DelivereeWebhookStatusDriverAcceptBooking — driver menerima booking.
	DelivereeWebhookStatusDriverAcceptBooking DelivereeWebhookStatus = "driver_accept_booking"
	// DelivereeWebhookStatusDeliveryInProgress — paket sedang dalam perjalanan.
	DelivereeWebhookStatusDeliveryInProgress DelivereeWebhookStatus = "delivery_in_progress"
	// DelivereeWebhookStatusLocatingDriverTimeout — tidak ada driver yang menerima.
	DelivereeWebhookStatusLocatingDriverTimeout DelivereeWebhookStatus = "locating_driver_timeout"
	// DelivereeWebhookStatusCanceled — booking dibatalkan provider.
	DelivereeWebhookStatusCanceled DelivereeWebhookStatus = "canceled"
)

// delivereeWebhookStatuses adalah daftar status yang dikenal & diproses.
var delivereeWebhookStatuses = map[DelivereeWebhookStatus]bool{
	DelivereeWebhookStatusDeliveryCompleted:     true,
	DelivereeWebhookStatusLocatingDriver:        true,
	DelivereeWebhookStatusDriverAcceptBooking:   true,
	DelivereeWebhookStatusDeliveryInProgress:    true,
	DelivereeWebhookStatusLocatingDriverTimeout: true,
	DelivereeWebhookStatusCanceled:              true,
}

// DelivereeWebhookService menangani event webhook dari Deliveree dan
// memperbarui status pesanan yang bersangkutan.
type DelivereeWebhookService interface {
	// Handle memproses payload webhook dan mengembalikan true jika ada pesanan
	// yang cocok & diperbarui.
	Handle(ctx context.Context, req *dto.DelivereeWebhookRequest) (bool, error)
}

type delivereeWebhookService struct {
	pesananRepo repositories.PesananRepository
	db          *gorm.DB
}

func NewDelivereeWebhookService(pesananRepo repositories.PesananRepository, db *gorm.DB) DelivereeWebhookService {
	return &delivereeWebhookService{
		pesananRepo: pesananRepo,
		db:          db,
	}
}

func (s *delivereeWebhookService) Handle(ctx context.Context, req *dto.DelivereeWebhookRequest) (bool, error) {
	status := DelivereeWebhookStatus(req.Status)

	// Status tidak dikenal → abaikan (jangan error, biar provider tidak retry terus)
	if !delivereeWebhookStatuses[status] {
		return false, nil
	}

	// Cari pesanan berdasarkan deliveree_booking_id
	bookingID := string(req.ID)
	if bookingID == "" {
		bookingID = string(req.NoBooking)
	}
	if bookingID == "" {
		return false, errors.New("id/no_booking wajib diisi")
	}

	var pesanan models.Pesanan
	if err := s.db.Where("deliveree_booking_id = ?", bookingID).First(&pesanan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Pesanan belum ditemukan — kemungkinan booking belum tercatat atau
			// webhook datang lebih awal. Biarkan provider mengirim ulang.
			return false, nil
		}
		return false, err
	}

	// Hanya proses jika pesanan memang memakai Deliveree
	if pesanan.DeliveryType != models.DeliveryTypeDeliveree {
		return false, nil
	}

	orderStatus, ok := s.mapToOrderStatus(status)
	if !ok {
		// Status diketahui tapi tidak mengubah order status (locating_driver,
		// driver_accept_booking) — hanya update info booking.
		extra := s.buildExtraUpdates(req)
		if len(extra) > 0 {
			if err := s.pesananRepo.UpdateFromWebhook(pesanan.ID, pesanan.OrderStatus, extra); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	extra := s.buildExtraUpdates(req)
	if err := s.pesananRepo.UpdateFromWebhook(pesanan.ID, orderStatus, extra); err != nil {
		return false, err
	}

	return true, nil
}

// mapToOrderStatus memetakan status webhook Deliveree ke OrderStatus aplikasi.
func (s *delivereeWebhookService) mapToOrderStatus(status DelivereeWebhookStatus) (models.OrderStatus, bool) {
	switch status {
	case DelivereeWebhookStatusDeliveryCompleted:
		return models.OrderStatusCompleted, true
	case DelivereeWebhookStatusDeliveryInProgress:
		return models.OrderStatusShipped, true
	case DelivereeWebhookStatusLocatingDriverTimeout, DelivereeWebhookStatusCanceled:
		return models.OrderStatusProcessing, true
	default:
		return "", false
	}
}

// buildExtraUpdates menyiapkan kolom pendukung dari payload webhook.
func (s *delivereeWebhookService) buildExtraUpdates(req *dto.DelivereeWebhookRequest) map[string]interface{} {
	extra := map[string]interface{}{
		"booking_status": req.Status,
	}
	if req.TrackingURL != "" {
		extra["tracking_url"] = req.TrackingURL
	}
	return extra
}
