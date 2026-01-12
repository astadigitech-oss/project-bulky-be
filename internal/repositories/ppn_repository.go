package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PPNRepository interface {
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.PPN, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.PPN, error)
	FindActive(ctx context.Context) (*models.PPN, error)
	Create(ctx context.Context, ppn *models.PPN) error
	Update(ctx context.Context, ppn *models.PPN) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActive(ctx context.Context, id uuid.UUID) error
}

type ppnRepository struct {
	db *gorm.DB
}

func NewPPNRepository(db *gorm.DB) PPNRepository {
	return &ppnRepository{db: db}
}

func (r *ppnRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.PPN, int64, error) {
	var ppns []models.PPN
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PPN{})

	// Filter by active status
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Search by persentase
	if params.Search != "" {
		query = query.Where("CAST(persentase AS TEXT) LIKE ?", "%"+params.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort - validate allowed fields
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	// Validate sort_by: hanya boleh persentase, is_active, created_at, updated_at
	allowedSortFields := map[string]bool{
		"persentase": true,
		"is_active":  true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at" // fallback ke default jika tidak valid
	}

	order := params.Order
	if order == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	// Pagination
	offset := (params.Page - 1) * params.PerPage
	if err := query.Offset(offset).Limit(params.PerPage).Find(&ppns).Error; err != nil {
		return nil, 0, err
	}

	return ppns, total, nil
}

func (r *ppnRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PPN, error) {
	var ppn models.PPN
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ppn).Error
	if err != nil {
		return nil, err
	}
	return &ppn, nil
}

func (r *ppnRepository) FindActive(ctx context.Context) (*models.PPN, error) {
	var ppn models.PPN
	err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&ppn).Error
	if err != nil {
		return nil, err
	}
	return &ppn, nil
}

func (r *ppnRepository) Create(ctx context.Context, ppn *models.PPN) error {
	return r.db.WithContext(ctx).Create(ppn).Error
}

func (r *ppnRepository) Update(ctx context.Context, ppn *models.PPN) error {
	return r.db.WithContext(ctx).Save(ppn).Error
}

func (r *ppnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PPN{}).Error
}

func (r *ppnRepository) SetActive(ctx context.Context, id uuid.UUID) error {
	// Trigger di database akan handle deactivate yang lain
	return r.db.WithContext(ctx).
		Model(&models.PPN{}).
		Where("id = ?", id).
		Update("is_active", true).Error
}
