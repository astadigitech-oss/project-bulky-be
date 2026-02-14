package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KuponRepository interface {
	Create(ctx context.Context, kupon *models.Kupon) error
	Update(ctx context.Context, kupon *models.Kupon) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Kupon, error)
	FindByKode(ctx context.Context, kode string) (*models.Kupon, error)
	FindAll(ctx context.Context, jenisDiskon *string, isActive, isExpired *bool, search, sortBy, order string, limit, offset int) ([]models.Kupon, int64, error)
	IsKodeExists(ctx context.Context, kode string, excludeID *uuid.UUID) (bool, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error

	// Kupon Kategori
	AddKategori(ctx context.Context, kuponID uuid.UUID, kategoriIDs []uuid.UUID) error
	RemoveAllKategori(ctx context.Context, kuponID uuid.UUID) error
	IsKategoriAllowed(ctx context.Context, kuponID, kategoriID uuid.UUID) (bool, error)

	// Kupon Usage
	CreateUsage(ctx context.Context, usage *models.KuponUsage) error
	FindUsagesByKuponID(ctx context.Context, kuponID uuid.UUID, limit, offset int) ([]models.KuponUsage, int64, error)
	GetUsageCount(ctx context.Context, kuponID uuid.UUID) (int64, error)
}

type kuponRepository struct {
	db *gorm.DB
}

func NewKuponRepository(db *gorm.DB) KuponRepository {
	return &kuponRepository{db: db}
}

func (r *kuponRepository) Create(ctx context.Context, kupon *models.Kupon) error {
	return r.db.WithContext(ctx).Create(kupon).Error
}

func (r *kuponRepository) Update(ctx context.Context, kupon *models.Kupon) error {
	return r.db.WithContext(ctx).Save(kupon).Error
}

func (r *kuponRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Kupon{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *kuponRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Kupon, error) {
	var kupon models.Kupon
	err := r.db.WithContext(ctx).
		Preload("Kategori.Kategori").
		First(&kupon, id).Error
	if err != nil {
		return nil, err
	}
	return &kupon, nil
}

func (r *kuponRepository) FindByKode(ctx context.Context, kode string) (*models.Kupon, error) {
	var kupon models.Kupon
	err := r.db.WithContext(ctx).
		Where("LOWER(kode) = LOWER(?) AND deleted_at IS NULL", kode).
		First(&kupon).Error
	if err != nil {
		return nil, err
	}
	return &kupon, nil
}

func (r *kuponRepository) FindAll(ctx context.Context, jenisDiskon *string, isActive, isExpired *bool, search, sortBy, order string, limit, offset int) ([]models.Kupon, int64, error) {
	var kupons []models.Kupon
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Kupon{})

	// Filter by jenis_diskon
	if jenisDiskon != nil && *jenisDiskon != "" {
		query = query.Where("jenis_diskon = ?", *jenisDiskon)
	}

	// Filter by is_active
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Filter by is_expired
	if isExpired != nil {
		today := time.Now().Format("2006-01-02")
		if *isExpired {
			query = query.Where("tanggal_kedaluarsa < ?", today)
		} else {
			query = query.Where("tanggal_kedaluarsa >= ?", today)
		}
	}

	// Search by kode or nama
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(kode) LIKE ? OR LOWER(nama) LIKE ?", searchPattern, searchPattern)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	// Special handling for total_usage sorting
	if sortBy == "total_usage" {
		query = query.
			Select("kupon.*, COUNT(kupon_usage.id) as usage_count").
			Joins("LEFT JOIN kupon_usage ON kupon_usage.kupon_id = kupon.id").
			Group("kupon.id").
			Order(fmt.Sprintf("usage_count %s", order))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Execute query
	if err := query.Find(&kupons).Error; err != nil {
		return nil, 0, err
	}

	return kupons, total, nil
}

func (r *kuponRepository) IsKodeExists(ctx context.Context, kode string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Kupon{}).
		Where("LOWER(kode) = LOWER(?) AND deleted_at IS NULL", kode)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kuponRepository) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Kupon{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).
		Error
}

// Kupon Kategori methods

func (r *kuponRepository) AddKategori(ctx context.Context, kuponID uuid.UUID, kategoriIDs []uuid.UUID) error {
	if len(kategoriIDs) == 0 {
		return nil
	}

	var kuponKategoris []models.KuponKategori
	for _, kategoriID := range kategoriIDs {
		kuponKategoris = append(kuponKategoris, models.KuponKategori{
			KuponID:    kuponID,
			KategoriID: kategoriID,
		})
	}

	return r.db.WithContext(ctx).Create(&kuponKategoris).Error
}

func (r *kuponRepository) RemoveAllKategori(ctx context.Context, kuponID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("kupon_id = ?", kuponID).
		Delete(&models.KuponKategori{}).
		Error
}

func (r *kuponRepository) IsKategoriAllowed(ctx context.Context, kuponID, kategoriID uuid.UUID) (bool, error) {
	// First check if kupon applies to all categories
	var kupon models.Kupon
	if err := r.db.WithContext(ctx).Select("is_all_kategori").First(&kupon, kuponID).Error; err != nil {
		return false, err
	}

	if kupon.IsAllKategori {
		return true, nil
	}

	// Check if kategori exists in kupon_kategori
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KuponKategori{}).
		Where("kupon_id = ? AND kategori_id = ?", kuponID, kategoriID).
		Count(&count).Error

	return count > 0, err
}

// Kupon Usage methods

func (r *kuponRepository) CreateUsage(ctx context.Context, usage *models.KuponUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *kuponRepository) FindUsagesByKuponID(ctx context.Context, kuponID uuid.UUID, limit, offset int) ([]models.KuponUsage, int64, error) {
	var usages []models.KuponUsage
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KuponUsage{}).
		Where("kupon_id = ?", kuponID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get usages with buyer and pesanan
	err := query.
		Preload("Buyer").
		Preload("Pesanan").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&usages).Error

	if err != nil {
		return nil, 0, err
	}

	return usages, total, nil
}

func (r *kuponRepository) GetUsageCount(ctx context.Context, kuponID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KuponUsage{}).
		Where("kupon_id = ?", kuponID).
		Count(&count).Error
	return count, err
}
