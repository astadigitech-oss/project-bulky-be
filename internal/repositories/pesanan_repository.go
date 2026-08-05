package repositories

import (
	"fmt"
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PesananRepository interface {
	FindByID(id uuid.UUID) (*models.Pesanan, error)
	FindByBuyerID(buyerID uuid.UUID) ([]models.Pesanan, error)

	// Admin methods
	AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Pesanan, int64, error)
	AdminFindByID(id uuid.UUID) (*models.Pesanan, error)
	UpdateStatus(id uuid.UUID, orderStatus models.OrderStatus, note *string, adminID uuid.UUID) error
	UpdateBookingResult(id uuid.UUID, delivereeBookingID *string, forwarderTrackingNo *string, bookingError *string) error
	ClearBookingResult(id uuid.UUID) error
	// UpdateFromWebhook memperbarui status pesanan dari event webhook provider
	// (Deliveree/Forwarder). Tidak memvalidasi transisi status karena event
	// webhook bisa datang kapan saja dan berulang (idempotent).
	UpdateFromWebhook(id uuid.UUID, orderStatus models.OrderStatus, extraUpdates map[string]interface{}) error
	Delete(id uuid.UUID) error
	GetStatistics(tanggalDari, tanggalSampai *time.Time) (map[string]interface{}, error)
	GetChartData(dari, sampai *time.Time, groupBy string) ([]models.ChartRawPoint, error)

	// Cancel order & restore produk
	CancelOrder(id uuid.UUID, reason *string, adminID uuid.UUID) error
}

type pesananRepository struct {
	db *gorm.DB
}

func NewPesananRepository(db *gorm.DB) PesananRepository {
	return &pesananRepository{db: db}
}

func (r *pesananRepository) FindByID(id uuid.UUID) (*models.Pesanan, error) {
	var pesanan models.Pesanan
	if err := r.db.First(&pesanan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pesanan, nil
}

func (r *pesananRepository) FindByBuyerID(buyerID uuid.UUID) ([]models.Pesanan, error) {
	var pesanan []models.Pesanan
	if err := r.db.Where("buyer_id = ?", buyerID).Find(&pesanan).Error; err != nil {
		return nil, err
	}
	return pesanan, nil
}

// ========================================
// Admin Methods
// ========================================

func (r *pesananRepository) AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Pesanan, int64, error) {
	var pesanan []models.Pesanan
	var total int64

	query := r.db.Model(&models.Pesanan{}).
		Preload("Buyer").
		Preload("Items").
		Preload("Items.Produk")

	// Apply filters
	if orderStatus, ok := filters["order_status"].(string); ok && orderStatus != "" {
		query = query.Where("order_status = ?", orderStatus)
	}
	if paymentType, ok := filters["payment_type"].(string); ok && paymentType != "" {
		query = query.Where("payment_type = ?", paymentType)
	}
	if buyer, ok := filters["buyer"].(string); ok && buyer != "" {
		query = query.Joins("JOIN buyer ON buyer.id = pesanan.buyer_id").
			Where("buyer.nama ILIKE ?", "%"+buyer+"%")
	}
	if cari, ok := filters["cari"].(string); ok && cari != "" {
		query = query.Joins("JOIN buyer ON buyer.id = pesanan.buyer_id").
			Where("pesanan.kode ILIKE ? OR buyer.nama ILIKE ?",
				"%"+cari+"%", "%"+cari+"%")
	}
	if tanggalDari, ok := filters["tanggal_dari"].(time.Time); ok {
		query = query.Where("pesanan.created_at >= ?", tanggalDari)
	}
	if tanggalSampai, ok := filters["tanggal_sampai"].(time.Time); ok {
		query = query.Where("pesanan.created_at <= ?", tanggalSampai)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	orderClause := fmt.Sprintf("pesanan.%s %s", sortBy, sortOrder)

	// Pagination
	offset := (page - 1) * perPage
	if err := query.Order(orderClause).Offset(offset).Limit(perPage).Find(&pesanan).Error; err != nil {
		return nil, 0, err
	}

	return pesanan, total, nil
}

func (r *pesananRepository) AdminFindByID(id uuid.UUID) (*models.Pesanan, error) {
	var pesanan models.Pesanan
	err := r.db.
		Preload("Buyer").
		Preload("AlamatBuyer").
		Preload("Items").
		Preload("Items.Produk").
		Preload("Items.Produk.Gambar", "is_primary = true").
		Preload("Pembayaran").
		Preload("Pembayaran.MetodePembayaran").
		Preload("Pembayaran.Buyer").
		First(&pesanan, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &pesanan, nil
}

func (r *pesananRepository) UpdateStatus(id uuid.UUID, orderStatus models.OrderStatus, note *string, adminID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current pesanan
		var pesanan models.Pesanan
		if err := tx.First(&pesanan, "id = ?", id).Error; err != nil {
			return err
		}

		// Validate status transition
		if !isValidStatusTransition(pesanan.OrderStatus, orderStatus) {
			return fmt.Errorf("tidak dapat mengubah status dari %s ke %s", pesanan.OrderStatus, orderStatus)
		}

		// Update pesanan status
		updates := map[string]interface{}{
			"order_status": orderStatus,
		}

		// Set timestamp based on status
		now := time.Now()
		switch orderStatus {
		case models.OrderStatusProcessing:
			updates["processed_at"] = now
		case models.OrderStatusReady:
			updates["ready_at"] = now
		case models.OrderStatusShipped:
			updates["shipped_at"] = now
		case models.OrderStatusCompleted:
			updates["completed_at"] = now
		case models.OrderStatusCancelled:
			updates["cancelled_at"] = now
			if note != nil {
				updates["cancelled_reason"] = *note
			}
		}

		// Capture statusFrom BEFORE update — GORM mutates the struct in memory after Updates()
		statusFrom := string(pesanan.OrderStatus)

		if err := tx.Model(&pesanan).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		// Create status history
		history := models.PesananStatusHistory{
			PesananID:  id,
			StatusFrom: &statusFrom,
			StatusTo:   string(orderStatus),
			StatusType: models.StatusHistoryTypeOrder,
			ChangedBy:  &adminID,
			Note:       note,
		}

		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelOrder sets order status to CANCELLED, records history, and restores is_sold on all produk items.
func (r *pesananRepository) CancelOrder(id uuid.UUID, reason *string, adminID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Load pesanan
		var pesanan models.Pesanan
		if err := tx.Preload("Items").First(&pesanan, "id = ?", id).Error; err != nil {
			return err
		}

		// 2. Validate: only non-terminal statuses can be cancelled
		if pesanan.OrderStatus == models.OrderStatusCompleted {
			return fmt.Errorf("pesanan dengan status COMPLETED tidak dapat dibatalkan")
		}
		if pesanan.OrderStatus == models.OrderStatusCancelled {
			return fmt.Errorf("pesanan sudah berstatus CANCELLED")
		}

		// 3. Update pesanan → CANCELLED
		now := time.Now()
		updates := map[string]interface{}{
			"order_status": models.OrderStatusCancelled,
			"cancelled_at": now,
		}
		if reason != nil && *reason != "" {
			updates["cancelled_reason"] = *reason
		}
		statusFrom := string(pesanan.OrderStatus)
		if err := tx.Model(&pesanan).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		// 4. Record status history
		history := models.PesananStatusHistory{
			PesananID:  id,
			StatusFrom: &statusFrom,
			StatusTo:   string(models.OrderStatusCancelled),
			StatusType: models.StatusHistoryTypeOrder,
			ChangedBy:  &adminID,
			Note:       reason,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		// 5. Restore is_sold = false for all produk in this order
		produkIDs := make([]uuid.UUID, 0, len(pesanan.Items))
		for _, item := range pesanan.Items {
			produkIDs = append(produkIDs, item.ProdukID)
		}
		if len(produkIDs) > 0 {
			if err := tx.Model(&models.Produk{}).Where("id IN ?", produkIDs).Update("is_sold", false).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *pesananRepository) ClearBookingResult(id uuid.UUID) error {
	return r.db.Model(&models.Pesanan{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"deliveree_booking_id":  nil,
		"forwarder_tracking_no": nil,
		"booking_error":         nil,
		"booking_lock_at":       nil,
	}).Error
}

// UpdateFromWebhook memperbarui status pesanan dari event webhook provider.
//
// Berbeda dari UpdateStatus (yang memvalidasi transisi status dan butuh adminID),
// method ini langsung mengubah status ke nilai yang diberikan oleh provider karena:
// 1. Event webhook bisa datang di luar urutan transisi normal (mis. timeout → in_progress).
// 2. Deliveree bisa mengirim ulang event yang sama (retry) — method ini idempotent.
//
// extraUpdates dipakai untuk kolom pendukung seperti booking_status/tracking_url
// yang dikirim provider bersama event.
func (r *pesananRepository) UpdateFromWebhook(id uuid.UUID, orderStatus models.OrderStatus, extraUpdates map[string]interface{}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pesanan models.Pesanan
		if err := tx.First(&pesanan, "id = ?", id).Error; err != nil {
			return err
		}

		// Idempotent: jika status sudah sama, hanya update kolom pendukung.
		if pesanan.OrderStatus == orderStatus {
			if len(extraUpdates) > 0 {
				return tx.Model(&models.Pesanan{}).Where("id = ?", id).UpdateColumns(extraUpdates).Error
			}
			return nil
		}

		updates := map[string]interface{}{
			"order_status": orderStatus,
		}
		for k, v := range extraUpdates {
			updates[k] = v
		}

		// Set timestamp berdasarkan status
		now := time.Now()
		switch orderStatus {
		case models.OrderStatusProcessing:
			updates["processed_at"] = now
		case models.OrderStatusShipped:
			updates["shipped_at"] = now
		case models.OrderStatusCompleted:
			updates["completed_at"] = now
		case models.OrderStatusCancelled:
			updates["cancelled_at"] = now
		}

		statusFrom := string(pesanan.OrderStatus)
		if err := tx.Model(&pesanan).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		// Catat riwayat status — ChangedBy nil karena bukan aksi admin
		history := models.PesananStatusHistory{
			PesananID:  id,
			StatusFrom: &statusFrom,
			StatusTo:   string(orderStatus),
			StatusType: models.StatusHistoryTypeOrder,
			Note:       ptrString("Update otomatis dari webhook provider"),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

func ptrString(s string) *string {
	return &s
}

func (r *pesananRepository) UpdateBookingResult(id uuid.UUID, delivereeBookingID *string, forwarderTrackingNo *string, bookingError *string) error {
	updates := map[string]interface{}{
		"booking_error":   bookingError,
		"booking_lock_at": nil, // release claim setelah booking selesai (sukses/gagal)
	}
	if delivereeBookingID != nil {
		updates["deliveree_booking_id"] = *delivereeBookingID
	}
	if forwarderTrackingNo != nil {
		updates["forwarder_tracking_no"] = *forwarderTrackingNo
	}
	return r.db.Model(&models.Pesanan{}).Where("id = ?", id).UpdateColumns(updates).Error
}

func (r *pesananRepository) Delete(id uuid.UUID) error {
	// Check if status is CANCELLED
	var pesanan models.Pesanan
	if err := r.db.First(&pesanan, "id = ?", id).Error; err != nil {
		return err
	}

	if pesanan.OrderStatus != models.OrderStatusCancelled {
		return fmt.Errorf("hanya pesanan dengan status CANCELLED yang dapat dihapus")
	}

	return r.db.Delete(&models.Pesanan{}, "id = ?", id).Error
}

func (r *pesananRepository) GetStatistics(tanggalDari, tanggalSampai *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// baseQuery builds a fresh query with date filters applied each time
	baseQuery := func() *gorm.DB {
		q := r.db.Model(&models.Pesanan{})
		if tanggalDari != nil {
			q = q.Where("created_at >= ?", tanggalDari)
		}
		if tanggalSampai != nil {
			q = q.Where("created_at <= ?", tanggalSampai)
		}
		return q
	}

	// Total pesanan
	var totalPesanan int64
	if err := baseQuery().Count(&totalPesanan).Error; err != nil {
		return nil, err
	}
	stats["total_pesanan"] = totalPesanan

	// Total revenue (only PAID orders)
	var totalRevenue decimal.Decimal
	if err := baseQuery().Where("payment_status = ?", "PAID").
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	stats["total_revenue"] = totalRevenue

	// Per status
	perStatus := make(map[string]int64)
	var statusCounts []struct {
		OrderStatus string
		Count       int64
	}
	if err := baseQuery().Select("order_status, COUNT(*) as count").
		Group("order_status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	for _, sc := range statusCounts {
		perStatus[sc.OrderStatus] = sc.Count
	}
	stats["per_status"] = perStatus

	// Per delivery type
	perDeliveryType := make(map[string]int64)
	var deliveryCounts []struct {
		DeliveryType string
		Count        int64
	}
	if err := baseQuery().Select("delivery_type, COUNT(*) as count").
		Group("delivery_type").
		Scan(&deliveryCounts).Error; err != nil {
		return nil, err
	}
	for _, dc := range deliveryCounts {
		perDeliveryType[dc.DeliveryType] = dc.Count
	}
	stats["per_delivery_type"] = perDeliveryType

	// Per payment status
	perPaymentStatus := make(map[string]int64)
	var paymentCounts []struct {
		PaymentStatus string
		Count         int64
	}
	if err := baseQuery().Select("payment_status, COUNT(*) as count").
		Group("payment_status").
		Scan(&paymentCounts).Error; err != nil {
		return nil, err
	}
	for _, pc := range paymentCounts {
		perPaymentStatus[pc.PaymentStatus] = pc.Count
	}
	stats["per_payment_status"] = perPaymentStatus

	return stats, nil
}

func (r *pesananRepository) GetChartData(dari, sampai *time.Time, groupBy string) ([]models.ChartRawPoint, error) {
	var format string
	if groupBy == "month" {
		format = "YYYY-MM"
	} else {
		format = "YYYY-MM-DD"
	}

	var results []models.ChartRawPoint
	q := r.db.Model(&models.Pesanan{}).
		Select("TO_CHAR(created_at AT TIME ZONE 'Asia/Jakarta', '" + format + "') as period, COUNT(*) as total_pesanan").
		Group("period").
		Order("period ASC")

	if dari != nil {
		q = q.Where("created_at >= ?", dari)
	}
	if sampai != nil {
		q = q.Where("created_at <= ?", sampai)
	}

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Helper function to validate status transitions
func isValidStatusTransition(from, to models.OrderStatus) bool {
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending: {
			models.OrderStatusProcessing,
			models.OrderStatusCancelled,
		},
		models.OrderStatusProcessing: {
			models.OrderStatusReady,
			models.OrderStatusCancelled,
		},
		models.OrderStatusReady: {
			models.OrderStatusShipped,
			models.OrderStatusCompleted, // PICKUP: langsung READY → COMPLETED tanpa SHIPPED
			models.OrderStatusCancelled,
		},
		models.OrderStatusShipped: {
			models.OrderStatusCompleted,
			models.OrderStatusCancelled,
		},
		models.OrderStatusCompleted: {},
		models.OrderStatusCancelled: {},
	}

	allowedTransitions, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}

	return false
}
