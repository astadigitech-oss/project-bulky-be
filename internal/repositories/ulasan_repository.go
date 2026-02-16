package repositories

import (
	"fmt"
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UlasanRepository interface {
	// Admin
	AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Ulasan, int64, error)
	AdminFindByID(id uuid.UUID) (*models.Ulasan, error)
	Approve(id uuid.UUID, isApproved bool, adminID uuid.UUID) error
	BulkApprove(ids []uuid.UUID, isApproved bool, adminID uuid.UUID) ([]uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetSummary() (map[string]int64, error)

	// Buyer
	GetPendingReviews(buyerID uuid.UUID) ([]models.PesananItem, error)
	BuyerFindAll(buyerID uuid.UUID, page, perPage int) ([]models.Ulasan, int64, error)
	Create(ulasan *models.Ulasan) error
	CheckExistingReview(pesananItemID uuid.UUID) (bool, error)

	// Public
	GetProdukUlasan(produkID uuid.UUID, filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Ulasan, int64, error)
	GetProdukRatingStats(produkID uuid.UUID) (*models.ProdukRatingStats, error)
}

type ulasanRepository struct {
	db *gorm.DB
}

func NewUlasanRepository(db *gorm.DB) UlasanRepository {
	return &ulasanRepository{db: db}
}

// ========================================
// Admin Methods
// ========================================

func (r *ulasanRepository) AdminFindAll(filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Ulasan, int64, error) {
	var ulasan []models.Ulasan
	var total int64

	query := r.db.Model(&models.Ulasan{}).
		Preload("Buyer").
		Preload("Pesanan").
		Preload("Produk").
		Preload("Approver")

	// Apply filters
	if isApproved, ok := filters["is_approved"].(bool); ok {
		query = query.Where("is_approved = ?", isApproved)
	}
	if rating, ok := filters["rating"].(int); ok {
		query = query.Where("rating = ?", rating)
	}
	if cari, ok := filters["cari"].(string); ok && cari != "" {
		query = query.Joins("JOIN buyer ON buyer.id = ulasan.buyer_id").
			Joins("JOIN pesanan ON pesanan.id = ulasan.pesanan_id").
			Joins("JOIN produk ON produk.id = ulasan.produk_id").
			Where("buyer.nama ILIKE ? OR pesanan.kode ILIKE ? OR produk.nama ILIKE ?",
				"%"+cari+"%", "%"+cari+"%", "%"+cari+"%")
	}
	if tanggalDari, ok := filters["tanggal_dari"].(time.Time); ok {
		query = query.Where("ulasan.created_at >= ?", tanggalDari)
	}
	if tanggalSampai, ok := filters["tanggal_sampai"].(time.Time); ok {
		query = query.Where("ulasan.created_at <= ?", tanggalSampai)
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
	orderClause := fmt.Sprintf("ulasan.%s %s", sortBy, sortOrder)

	// Pagination
	offset := (page - 1) * perPage
	if err := query.Order(orderClause).Offset(offset).Limit(perPage).Find(&ulasan).Error; err != nil {
		return nil, 0, err
	}

	return ulasan, total, nil
}

func (r *ulasanRepository) AdminFindByID(id uuid.UUID) (*models.Ulasan, error) {
	var ulasan models.Ulasan
	if err := r.db.Preload("Buyer").
		Preload("Pesanan").
		Preload("Produk").
		Preload("Approver").
		First(&ulasan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ulasan, nil
}

func (r *ulasanRepository) Approve(id uuid.UUID, isApproved bool, adminID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_approved": isApproved,
		"approved_by": adminID,
	}

	if isApproved {
		updates["approved_at"] = time.Now()
	} else {
		updates["approved_at"] = nil
	}

	return r.db.Model(&models.Ulasan{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ulasanRepository) BulkApprove(ids []uuid.UUID, isApproved bool, adminID uuid.UUID) ([]uuid.UUID, error) {
	// First, get valid IDs that exist in database
	var validUlasan []models.Ulasan
	if err := r.db.Model(&models.Ulasan{}).Select("id").Where("id IN ?", ids).Find(&validUlasan).Error; err != nil {
		return nil, err
	}

	// If no valid IDs found, return error
	if len(validUlasan) == 0 {
		return nil, fmt.Errorf("tidak ada ulasan yang valid untuk di-approve")
	}

	// Extract valid IDs
	var validIDs []uuid.UUID
	for _, ulasan := range validUlasan {
		validIDs = append(validIDs, ulasan.ID)
	}

	// Prepare updates
	updates := map[string]interface{}{
		"is_approved": isApproved,
		"approved_by": adminID,
	}

	if isApproved {
		updates["approved_at"] = time.Now()
	} else {
		updates["approved_at"] = nil
	}

	// Update only valid IDs
	if err := r.db.Model(&models.Ulasan{}).Where("id IN ?", validIDs).Updates(updates).Error; err != nil {
		return nil, err
	}

	return validIDs, nil
}

func (r *ulasanRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Ulasan{}, "id = ?", id).Error
}

func (r *ulasanRepository) GetSummary() (map[string]int64, error) {
	summary := make(map[string]int64)

	// Total pending
	var pending int64
	if err := r.db.Model(&models.Ulasan{}).Where("is_approved = false").Count(&pending).Error; err != nil {
		return nil, err
	}
	summary["total_pending"] = pending

	// Total approved
	var approved int64
	if err := r.db.Model(&models.Ulasan{}).Where("is_approved = true").Count(&approved).Error; err != nil {
		return nil, err
	}
	summary["total_approved"] = approved

	// Total rejected (considered deleted or not approved)
	summary["total_rejected"] = 0 // Can be calculated from total - approved - pending

	return summary, nil
}

// ========================================
// Buyer Methods
// ========================================

func (r *ulasanRepository) GetPendingReviews(buyerID uuid.UUID) ([]models.PesananItem, error) {
	var items []models.PesananItem

	// Get pesanan items from COMPLETED orders that haven't been reviewed yet
	err := r.db.
		Joins("JOIN pesanan ON pesanan.id = pesanan_item.pesanan_id").
		Preload("Produk").
		Preload("Produk.Gambar", "is_primary = true").
		Preload("Pesanan"). // Add Pesanan preload
		Where("pesanan.buyer_id = ? AND pesanan.order_status = ?", buyerID, "COMPLETED").
		Where("pesanan_item.id NOT IN (SELECT pesanan_item_id FROM ulasan WHERE deleted_at IS NULL)").
		Find(&items).Error

	return items, err
}

func (r *ulasanRepository) BuyerFindAll(buyerID uuid.UUID, page, perPage int) ([]models.Ulasan, int64, error) {
	var ulasan []models.Ulasan
	var total int64

	query := r.db.Model(&models.Ulasan{}).
		Preload("Produk").
		Preload("Produk.Gambar", "is_primary = true").
		Where("buyer_id = ?", buyerID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&ulasan).Error; err != nil {
		return nil, 0, err
	}

	return ulasan, total, nil
}

func (r *ulasanRepository) Create(ulasan *models.Ulasan) error {
	return r.db.Create(ulasan).Error
}

func (r *ulasanRepository) CheckExistingReview(pesananItemID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Ulasan{}).Where("pesanan_item_id = ?", pesananItemID).Count(&count).Error
	return count > 0, err
}

// ========================================
// Public Methods
// ========================================

func (r *ulasanRepository) GetProdukUlasan(produkID uuid.UUID, filters map[string]interface{}, page, perPage int, sortBy, sortOrder string) ([]models.Ulasan, int64, error) {
	var ulasan []models.Ulasan
	var total int64

	query := r.db.Model(&models.Ulasan{}).
		Preload("Buyer").
		Where("produk_id = ? AND is_approved = true", produkID)

	// Apply filters
	if rating, ok := filters["rating"].(int); ok {
		query = query.Where("rating = ?", rating)
	}
	if withPhoto, ok := filters["with_photo"].(bool); ok && withPhoto {
		query = query.Where("gambar IS NOT NULL")
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
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	// Pagination
	offset := (page - 1) * perPage
	if err := query.Order(orderClause).Offset(offset).Limit(perPage).Find(&ulasan).Error; err != nil {
		return nil, 0, err
	}

	return ulasan, total, nil
}

func (r *ulasanRepository) GetProdukRatingStats(produkID uuid.UUID) (*models.ProdukRatingStats, error) {
	stats := &models.ProdukRatingStats{
		ProdukID: produkID.String(),
	}

	// Total & Average
	var result struct {
		TotalUlasan int
		RataRating  float64
	}

	err := r.db.Model(&models.Ulasan{}).
		Where("produk_id = ? AND is_approved = true", produkID).
		Select("COUNT(*) as total_ulasan, COALESCE(AVG(rating), 0) as rata_rating").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	stats.TotalUlasan = result.TotalUlasan
	stats.RataRating = result.RataRating

	// Count per rating
	type RatingCount struct {
		Rating int
		Count  int
	}
	var counts []RatingCount

	r.db.Model(&models.Ulasan{}).
		Where("produk_id = ? AND is_approved = true", produkID).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&counts)

	for _, c := range counts {
		switch c.Rating {
		case 5:
			stats.Rating5 = c.Count
		case 4:
			stats.Rating4 = c.Count
		case 3:
			stats.Rating3 = c.Count
		case 2:
			stats.Rating2 = c.Count
		case 1:
			stats.Rating1 = c.Count
		}
	}

	return stats, nil
}
