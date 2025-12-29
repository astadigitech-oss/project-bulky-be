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
	Page          int
	PerPage       int
	UserType      string
	UserID        *uuid.UUID
	Action        string
	Modul         string
	TanggalDari   string
	TanggalSampai string
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

	// Apply filters
	if filter.UserType != "" {
		query = query.Where("user_type = ?", filter.UserType)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Modul != "" {
		query = query.Where("modul = ?", filter.Modul)
	}
	if filter.TanggalDari != "" {
		query = query.Where("created_at >= ?", filter.TanggalDari)
	}
	if filter.TanggalSampai != "" {
		query = query.Where("created_at <= ?", filter.TanggalSampai)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PerPage
	query = query.Order("created_at DESC").Limit(filter.PerPage).Offset(offset)

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
