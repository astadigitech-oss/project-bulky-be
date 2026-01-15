package repositories

import (
	"context"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranRepository interface {
	FindAll(ctx context.Context, params *models.PaginationRequest, groupID *uuid.UUID, isActive *bool) ([]models.MetodePembayaran, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error)
	FindByIDWithGroup(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error)
	Create(ctx context.Context, metode *models.MetodePembayaran) error
	Update(ctx context.Context, metode *models.MetodePembayaran) error
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleStatus(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error)
	CheckByKode(ctx context.Context, kode string, excludeID *uuid.UUID) (bool, error)
	CheckUsedInTransaction(ctx context.Context, id uuid.UUID) (bool, error)
}

type metodePembayaranRepository struct {
	db *gorm.DB
}

func NewMetodePembayaranRepository(db *gorm.DB) MetodePembayaranRepository {
	return &metodePembayaranRepository{db: db}
}

func (r *metodePembayaranRepository) FindAll(ctx context.Context, params *models.PaginationRequest, groupID *uuid.UUID, isActive *bool) ([]models.MetodePembayaran, int64, error) {
	var metodes []models.MetodePembayaran
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MetodePembayaran{}).Preload("Group")

	// Filter by group_id
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	// Filter by is_active
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Search by nama or kode
	if params.Search != "" {
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort - validate allowed fields
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "urutan"
	}
	allowedSortFields := map[string]bool{
		"nama":       true,
		"kode":       true,
		"urutan":     true,
		"is_active":  true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "urutan"
	}

	order := params.Order
	if order == "" {
		order = "asc"
	}
	query = query.Order(sortBy + " " + order)

	// Pagination
	offset := (params.Page - 1) * params.PerPage
	if err := query.Offset(offset).Limit(params.PerPage).Find(&metodes).Error; err != nil {
		return nil, 0, err
	}

	return metodes, total, nil
}

func (r *metodePembayaranRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error) {
	var metode models.MetodePembayaran
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&metode).Error
	if err != nil {
		return nil, err
	}
	return &metode, nil
}

func (r *metodePembayaranRepository) FindByIDWithGroup(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error) {
	var metode models.MetodePembayaran
	err := r.db.WithContext(ctx).
		Preload("Group").
		Where("id = ?", id).
		First(&metode).Error
	if err != nil {
		return nil, err
	}
	return &metode, nil
}

func (r *metodePembayaranRepository) Create(ctx context.Context, metode *models.MetodePembayaran) error {
	return r.db.WithContext(ctx).Create(metode).Error
}

func (r *metodePembayaranRepository) Update(ctx context.Context, metode *models.MetodePembayaran) error {
	return r.db.WithContext(ctx).Save(metode).Error
}

func (r *metodePembayaranRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MetodePembayaran{}).Error
}

func (r *metodePembayaranRepository) ToggleStatus(ctx context.Context, id uuid.UUID) (*models.MetodePembayaran, error) {
	var metode models.MetodePembayaran

	// Find the metode
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&metode).Error; err != nil {
		return nil, err
	}

	// Toggle is_active
	metode.IsActive = !metode.IsActive

	// Update
	if err := r.db.WithContext(ctx).Save(&metode).Error; err != nil {
		return nil, err
	}

	return &metode, nil
}

func (r *metodePembayaranRepository) CheckByKode(ctx context.Context, kode string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.MetodePembayaran{}).Where("kode = ?", kode)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *metodePembayaranRepository) CheckUsedInTransaction(ctx context.Context, id uuid.UUID) (bool, error) {
	// Check if metode pembayaran is used in pesanan_pembayaran table
	var count int64
	err := r.db.WithContext(ctx).
		Table("pesanan_pembayaran").
		Where("metode_pembayaran_id = ?", id).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
