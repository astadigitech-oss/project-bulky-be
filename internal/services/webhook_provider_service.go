package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"gorm.io/gorm"
)

// ProviderWebhookStatus adalah status pengiriman yang dikirim provider
// (Deliveree/Forwarder on-demand) via webhook. Format status sama untuk kedua
// provider karena mereka berbagi platform on-demand yang sama.
type ProviderWebhookStatus string

const (
	// ProviderWebhookStatusDeliveryCompleted — paket sudah diterima penerima.
	ProviderWebhookStatusDeliveryCompleted ProviderWebhookStatus = "delivery_completed"
	// ProviderWebhookStatusLocatingDriver — sedang mencari driver.
	ProviderWebhookStatusLocatingDriver ProviderWebhookStatus = "locating_driver"
	// ProviderWebhookStatusDriverAcceptBooking — driver menerima booking.
	ProviderWebhookStatusDriverAcceptBooking ProviderWebhookStatus = "driver_accept_booking"
	// ProviderWebhookStatusDeliveryInProgress — paket sedang dalam perjalanan.
	ProviderWebhookStatusDeliveryInProgress ProviderWebhookStatus = "delivery_in_progress"
	// ProviderWebhookStatusLocatingDriverTimeout — tidak ada driver yang menerima.
	ProviderWebhookStatusLocatingDriverTimeout ProviderWebhookStatus = "locating_driver_timeout"
	// ProviderWebhookStatusCanceled — booking dibatalkan provider.
	ProviderWebhookStatusCanceled ProviderWebhookStatus = "canceled"
)

// providerWebhookStatuses adalah daftar status yang dikenal & diproses.
var providerWebhookStatuses = map[ProviderWebhookStatus]bool{
	ProviderWebhookStatusDeliveryCompleted:     true,
	ProviderWebhookStatusLocatingDriver:        true,
	ProviderWebhookStatusDriverAcceptBooking:   true,
	ProviderWebhookStatusDeliveryInProgress:    true,
	ProviderWebhookStatusLocatingDriverTimeout: true,
	ProviderWebhookStatusCanceled:              true,
}

// mapProviderWebhookToOrderStatus memetakan status webhook provider ke
// OrderStatus aplikasi.
//
// Catatan pemetaan:
//   - locating_driver_timeout / canceled terjadi SEBELUM pengiriman dimulai
//     (masih tahap mencari driver). Status pesanan kembali ke READY — artinya
//     paket siap dikirim tapi booking gagal/tidak ada driver, bukan kembali ke
//     PROCESSING (yang berarti sedang diproses di gudang). Admin bisa retry
//     booking dari status READY.
func mapProviderWebhookToOrderStatus(status ProviderWebhookStatus) (models.OrderStatus, bool) {
	switch status {
	case ProviderWebhookStatusDeliveryCompleted:
		return models.OrderStatusCompleted, true
	case ProviderWebhookStatusDeliveryInProgress:
		return models.OrderStatusShipped, true
	case ProviderWebhookStatusLocatingDriverTimeout, ProviderWebhookStatusCanceled:
		return models.OrderStatusReady, true
	default:
		return "", false
	}
}

// providerWebhookConfig membedakan perilaku webhook per provider:
// kolom mana yang dipakai untuk mencari pesanan dan delivery type mana yang
// berhak menerima event dari provider ini.
type providerWebhookConfig struct {
	// lookupColumn adalah kolom di tabel pesanan yang diisi dengan identifier
	// booking/tracking provider (deliveree_booking_id / forwarder_tracking_no).
	lookupColumn string
	// deliveryTypes adalah daftar DeliveryType yang sah untuk provider ini.
	deliveryTypes map[models.DeliveryType]bool
}

// ProviderWebhookHandler adalah handler generik untuk event webhook dari
// provider pengiriman on-demand (Deliveree/Forwarder). Kedua provider memakai
// format status yang sama, hanya endpoint dan kolom identifikasi yang beda.
type ProviderWebhookHandler struct {
	pesananRepo repositories.PesananRepository
	db          *gorm.DB
	cfg         providerWebhookConfig
}

// NewProviderWebhookHandler membuat handler generik untuk satu provider.
func NewProviderWebhookHandler(pesananRepo repositories.PesananRepository, db *gorm.DB, cfg providerWebhookConfig) *ProviderWebhookHandler {
	return &ProviderWebhookHandler{
		pesananRepo: pesananRepo,
		db:          db,
		cfg:         cfg,
	}
}

// Handle memproses payload webhook dan mengembalikan true jika ada pesanan
// yang cocok & diperbarui.
func (h *ProviderWebhookHandler) Handle(ctx context.Context, identifier string, status string, trackingURL string) (bool, error) {
	ws := ProviderWebhookStatus(status)

	// Status tidak dikenal → abaikan (jangan error, biar provider tidak retry terus)
	if !providerWebhookStatuses[ws] {
		return false, nil
	}

	// Identifier booking wajib terisi
	if identifier == "" {
		return false, errors.New("identifier booking wajib diisi")
	}

	var pesanan models.Pesanan
	if err := h.db.Where(h.cfg.lookupColumn+" = ?", identifier).First(&pesanan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Pesanan belum ditemukan — kemungkinan booking belum tercatat atau
			// webhook datang lebih awal. Biarkan provider mengirim ulang.
			return false, nil
		}
		return false, err
	}

	// Hanya proses jika pesanan memang memakai provider ini
	if !h.cfg.deliveryTypes[pesanan.DeliveryType] {
		return false, nil
	}

	orderStatus, ok := mapProviderWebhookToOrderStatus(ws)
	if !ok {
		// Status diketahui tapi tidak mengubah order status (locating_driver,
		// driver_accept_booking) — hanya update info booking.
		extra := h.buildExtraUpdates(status, trackingURL)
		if len(extra) > 0 {
			if err := h.pesananRepo.UpdateFromWebhook(pesanan.ID, pesanan.OrderStatus, extra); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	// Anti-lompat-status: event timeout/canceled (→ READY) hanya boleh mengubah
	// order_status jika pesanan memang sudah READY. Jika masih PROCESSING/PENDING
	// (mis. booking via retry terjadi lebih awal dari admin set READY), webhook
	// cukup meng-update kolom pendukung — status gudang tetap kendali admin.
	// Perlindungan SHIPPED/COMPLETED ditangani di UpdateFromWebhook.
	if orderStatus == models.OrderStatusReady && pesanan.OrderStatus != models.OrderStatusReady {
		orderStatus = pesanan.OrderStatus
	}

	extra := h.buildExtraUpdates(status, trackingURL)
	if err := h.pesananRepo.UpdateFromWebhook(pesanan.ID, orderStatus, extra); err != nil {
		return false, err
	}

	return true, nil
}

// buildExtraUpdates menyiapkan kolom pendukung dari payload webhook.
func (h *ProviderWebhookHandler) buildExtraUpdates(status, trackingURL string) map[string]interface{} {
	extra := map[string]interface{}{
		"booking_status": status,
	}
	if trackingURL != "" {
		extra["tracking_url"] = trackingURL
	}
	return extra
}

// ─── Konfigurasi per provider ────────────────────────────────────────────────

// delivereeWebhookConfig: pesanan diidentifikasi lewat deliveree_booking_id.
func delivereeWebhookConfig() providerWebhookConfig {
	return providerWebhookConfig{
		lookupColumn: "deliveree_booking_id",
		deliveryTypes: map[models.DeliveryType]bool{
			models.DeliveryTypeDeliveree: true,
		},
	}
}

// forwarderWebhookConfig: pesanan diidentifikasi lewat forwarder_tracking_no.
func forwarderWebhookConfig() providerWebhookConfig {
	return providerWebhookConfig{
		lookupColumn: "forwarder_tracking_no",
		deliveryTypes: map[models.DeliveryType]bool{
			models.DeliveryTypeForwarder:    true,
			models.DeliveryTypeForwarderLCL: true,
		},
	}
}
