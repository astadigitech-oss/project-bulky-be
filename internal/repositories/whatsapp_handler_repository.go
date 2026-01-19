package repositories

import (
	"context"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WhatsAppHandlerRepository interface {
	Create(ctx context.Context, handler *models.WhatsAppHandler) error
	FindAll(ctx context.Context, params *models.WhatsAppHandlerFilterRequest) ([]models.WhatsAppHandler, int64, error)
	FindAllSimple(ctx context.Context) ([]models.WhatsAppHandler, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.WhatsAppHandler, error)
	Update(ctx context.Context, handler *models.WhatsAppHandler) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActive(ctx context.Context) (*models.WhatsAppHandler, error)
	SetActive(ctx context.Context, id uuid.UUID) error
	DeactivateAll(ctx context.Context) error
}

type whatsAppHandlerRepository struct {
	db *gorm.DB
}

func NewWhatsAppHandlerRepository(db *gorm.DB) WhatsAppHandlerRepository {
	return &whatsAppHandlerRepository{db: db}
}

func (r *whatsAppHandlerRepository) Create(ctx context.Context, handler *models.WhatsAppHandler) error {
	return r.db.WithContext(ctx).Create(handler).Error
}

func (r *whatsAppHandlerRepository) FindAll(ctx context.Context, params *models.WhatsAppHandlerFilterRequest) ([]models.WhatsAppHandler, int64, error) {
	var items []models.WhatsAppHandler
	var total int64

	query := r.db.WithContext(ctx).Model(&models.WhatsAppHandler{})

	// Search
	if params.Search != "" {
		query = query.Where("nomor_wa ILIKE ? OR pesan_awal ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Filter by active status
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate sort fields
	validSortFields := map[string]bool{
		"nomor_wa":   true,
		"is_active":  true,
		"created_at": true,
		"updated_at": true,
	}

	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "updated_at"
	}

	order := params.Order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Sorting
	orderClause := sortBy + " " + order
	query = query.Order(orderClause)

	// Pagination
	if params.PerPage > 0 {
		query = query.Offset(params.GetOffset()).Limit(params.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *whatsAppHandlerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.WhatsAppHandler, error) {
	var handler models.WhatsAppHandler
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&handler).Error; err != nil {
		return nil, err
	}
	return &handler, nil
}

func (r *whatsAppHandlerRepository) Update(ctx context.Context, handler *models.WhatsAppHandler) error {
	return r.db.WithContext(ctx).Save(handler).Error
}

func (r *whatsAppHandlerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.WhatsAppHandler{}).Error
}

func (r *whatsAppHandlerRepository) GetActive(ctx context.Context) (*models.WhatsAppHandler, error) {
	var handler models.WhatsAppHandler
	if err := r.db.WithContext(ctx).Where("is_active = ? AND deleted_at IS NULL", true).First(&handler).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &handler, nil
}

func (r *whatsAppHandlerRepository) SetActive(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deactivate all
		if err := tx.Model(&models.WhatsAppHandler{}).
			Where("deleted_at IS NULL").
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the selected one
		if err := tx.Model(&models.WhatsAppHandler{}).
			Where("id = ?", id).
			Update("is_active", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *whatsAppHandlerRepository) FindAllSimple(ctx context.Context) ([]models.WhatsAppHandler, error) {
	var items []models.WhatsAppHandler
	if err := r.db.WithContext(ctx).Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *whatsAppHandlerRepository) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&models.WhatsAppHandler{}).
		Where("deleted_at IS NULL").
		Update("is_active", false).Error
}
