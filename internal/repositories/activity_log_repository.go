package repositories

import (
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	Create(log *models.ActivityLog) error
	FindAll(filter ActivityLogFilter) ([]models.ActivityLog, int64, error)
	FindByID(id uuid.UUID) (*models.ActivityLog, error)
	FindByEntity(entityType string, entityID uuid.UUID) ([]models.ActivityLog, error)
}

type ActivityLogFilter struct {
	Page    int
	PerPage int
	Search  string // ADMIN, BUYER, SYSTEM (user_type)
	SortBy  string // field untuk sorting
	Order   string // asc atau desc
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *activityLogRepository) FindAll(filter ActivityLogFilter) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	query := r.db.Model(&models.ActivityLog{})

	// Apply search filter (user_type: ADMIN, BUYER, SYSTEM)
	if filter.Search != "" {
		query = query.Where("user_type = ?", filter.Search)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at" // default sort by created_at
	}
	order := filter.Order
	if order == "" {
		order = "desc" // default descending
	}
	orderClause := sortBy + " " + order

	// Apply pagination
	offset := (filter.Page - 1) * filter.PerPage
	query = query.Order(orderClause).Limit(filter.PerPage).Offset(offset)

	// Execute query
	err := query.Find(&logs).Error
	return logs, total, err
}

func (r *activityLogRepository) FindByID(id uuid.UUID) (*models.ActivityLog, error) {
	var log models.ActivityLog
	err := r.db.Where("id = ?", id).First(&log).Error
	return &log, err
}

func (r *activityLogRepository) FindByEntity(entityType string, entityID uuid.UUID) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
