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
	Delete(id uuid.UUID) error
	GetStatistics(tanggalDari, tanggalSampai *time.Time) (map[string]interface{}, error)
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
	if paymentStatus, ok := filters["payment_status"].(string); ok && paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}
	if deliveryType, ok := filters["delivery_type"].(string); ok && deliveryType != "" {
		query = query.Where("delivery_type = ?", deliveryType)
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

		if err := tx.Model(&pesanan).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		// Create status history
		statusFrom := string(pesanan.OrderStatus)
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

	query := r.db.Model(&models.Pesanan{})

	// Apply date filters
	if tanggalDari != nil {
		query = query.Where("created_at >= ?", tanggalDari)
	}
	if tanggalSampai != nil {
		query = query.Where("created_at <= ?", tanggalSampai)
	}

	// Total pesanan
	var totalPesanan int64
	if err := query.Count(&totalPesanan).Error; err != nil {
		return nil, err
	}
	stats["total_pesanan"] = totalPesanan

	// Total revenue (only PAID orders)
	var totalRevenue decimal.Decimal
	if err := query.Where("payment_status = ?", "PAID").
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
	if err := query.Select("order_status, COUNT(*) as count").
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
	if err := query.Select("delivery_type, COUNT(*) as count").
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
	if err := query.Select("payment_status, COUNT(*) as count").
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
