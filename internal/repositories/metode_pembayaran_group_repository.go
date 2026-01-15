package repositories

import (
	"context"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetodePembayaranGroupRepository interface {
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MetodePembayaranGroup, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error)
	FindByIDWithMetode(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error)
	Create(ctx context.Context, group *models.MetodePembayaranGroup) error
	Update(ctx context.Context, group *models.MetodePembayaranGroup) error
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleStatus(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error)
	CheckByName(ctx context.Context, nama string, excludeID *uuid.UUID) (bool, error)
	CountActiveMetode(ctx context.Context, groupID uuid.UUID) (int64, error)
}

type metodePembayaranGroupRepository struct {
	db *gorm.DB
}

func NewMetodePembayaranGroupRepository(db *gorm.DB) MetodePembayaranGroupRepository {
	return &metodePembayaranGroupRepository{db: db}
}

func (r *metodePembayaranGroupRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MetodePembayaranGroup, int64, error) {
	var groups []models.MetodePembayaranGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MetodePembayaranGroup{})

	// Search by nama
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
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
	if err := query.Offset(offset).Limit(params.PerPage).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *metodePembayaranGroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error) {
	var group models.MetodePembayaranGroup
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *metodePembayaranGroupRepository) FindByIDWithMetode(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error) {
	var group models.MetodePembayaranGroup
	err := r.db.WithContext(ctx).
		Preload("MetodePembayaran", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan ASC")
		}).
		Where("id = ?", id).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *metodePembayaranGroupRepository) Create(ctx context.Context, group *models.MetodePembayaranGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *metodePembayaranGroupRepository) Update(ctx context.Context, group *models.MetodePembayaranGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *metodePembayaranGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MetodePembayaranGroup{}).Error
}

func (r *metodePembayaranGroupRepository) ToggleStatus(ctx context.Context, id uuid.UUID) (*models.MetodePembayaranGroup, error) {
	var group models.MetodePembayaranGroup

	// Find the group
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}

	// Toggle is_active
	group.IsActive = !group.IsActive

	// Update
	if err := r.db.WithContext(ctx).Save(&group).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *metodePembayaranGroupRepository) CheckByName(ctx context.Context, nama string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.MetodePembayaranGroup{}).Where("nama = ?", nama)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *metodePembayaranGroupRepository) CountActiveMetode(ctx context.Context, groupID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.MetodePembayaran{}).
		Where("group_id = ? AND is_active = ?", groupID, true).
		Count(&count).Error

	return count, err
}
