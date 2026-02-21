package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type DisclaimerRepository interface {
	Create(ctx context.Context, disclaimer *models.Disclaimer) error
	FindByID(ctx context.Context, id string) (*models.Disclaimer, error)
	FindBySlug(ctx context.Context, slug string) (*models.Disclaimer, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Disclaimer, int64, error)
	Update(ctx context.Context, disclaimer *models.Disclaimer) error
	Delete(ctx context.Context, disclaimer *models.Disclaimer) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetActive(ctx context.Context) (*models.Disclaimer, error)
}

type disclaimerRepository struct {
	db *gorm.DB
}

func NewDisclaimerRepository(db *gorm.DB) DisclaimerRepository {
	return &disclaimerRepository{db: db}
}

func (r *disclaimerRepository) Create(ctx context.Context, disclaimer *models.Disclaimer) error {
	return r.db.WithContext(ctx).Create(disclaimer).Error
}

func (r *disclaimerRepository) FindByID(ctx context.Context, id string) (*models.Disclaimer, error) {
	var disclaimer models.Disclaimer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&disclaimer).Error
	if err != nil {
		return nil, err
	}
	return &disclaimer, nil
}

func (r *disclaimerRepository) FindBySlug(ctx context.Context, slug string) (*models.Disclaimer, error) {
	var disclaimer models.Disclaimer
	err := r.db.WithContext(ctx).Where("(slug_id = ? OR slug_en = ?) AND is_active = ?", slug, slug, true).First(&disclaimer).Error
	if err != nil {
		return nil, err
	}
	return &disclaimer, nil
}

func (r *disclaimerRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Disclaimer, int64, error) {
	var disclaimers []models.Disclaimer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Disclaimer{})

	if params.Search != "" {
		query = query.Where("judul ILIKE ? OR judul_en ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"judul":      true,
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

	orderClause := sortBy + " " + order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&disclaimers).Error; err != nil {
		return nil, 0, err
	}

	return disclaimers, total, nil
}

func (r *disclaimerRepository) Update(ctx context.Context, disclaimer *models.Disclaimer) error {
	return r.db.WithContext(ctx).Save(disclaimer).Error
}

func (r *disclaimerRepository) Delete(ctx context.Context, disclaimer *models.Disclaimer) error {
	// Manual update slug untuk soft delete
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		shortID := disclaimer.ID.String()[:8]
		suffix := "_deleted_" + shortID

		updates := map[string]interface{}{
			"deleted_at": now,
		}

		if disclaimer.Slug != nil && *disclaimer.Slug != "" {
			v := *disclaimer.Slug + suffix
			updates["slug"] = v
		}
		if disclaimer.SlugID != nil && *disclaimer.SlugID != "" {
			v := *disclaimer.SlugID + suffix
			updates["slug_id"] = v
		}
		if disclaimer.SlugEN != nil && *disclaimer.SlugEN != "" {
			v := *disclaimer.SlugEN + suffix
			updates["slug_en"] = v
		}

		if err := tx.Model(disclaimer).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *disclaimerRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Disclaimer{}).Where("slug_id = ? OR slug_en = ?", slug, slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *disclaimerRepository) GetActive(ctx context.Context) (*models.Disclaimer, error) {
	var disclaimer models.Disclaimer
	err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&disclaimer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &disclaimer, nil
}
