package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type MerekProdukRepository interface {
	Create(ctx context.Context, merek *models.MerekProduk) error
	FindByID(ctx context.Context, id string) (*models.MerekProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.MerekProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProdukSimpleResponse, int64, error)
	Update(ctx context.Context, merek *models.MerekProduk) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.MerekProduk, error)
}

type merekProdukRepository struct {
	db *gorm.DB
}

func NewMerekProdukRepository(db *gorm.DB) MerekProdukRepository {
	return &merekProdukRepository{db: db}
}

func (r *merekProdukRepository) Create(ctx context.Context, merek *models.MerekProduk) error {
	return r.db.WithContext(ctx).Create(merek).Error
}

func (r *merekProdukRepository) FindByID(ctx context.Context, id string) (*models.MerekProduk, error) {
	var merek models.MerekProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&merek).Error
	if err != nil {
		return nil, err
	}
	return &merek, nil
}

func (r *merekProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.MerekProduk, error) {
	var merek models.MerekProduk
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&merek).Error
	if err != nil {
		return nil, err
	}
	return &merek, nil
}

func (r *merekProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProdukSimpleResponse, int64, error) {
	var mereks []struct {
		ID        string    `json:"id"`
		NamaID    string    `json:"nama_id"`
		NamaEN    *string   `json:"nama_en"`
		LogoURL   *string   `json:"logo_url"`
		IsActive  bool      `json:"is_active"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	var total int64

	// Base query with simplified fields for list view
	query := r.db.WithContext(ctx).
		Model(&models.MerekProduk{}).
		Select(`
			id,
			nama_id,
			nama_en,
			logo_url,
			is_active,
			updated_at
		`)

	// Search filter
	if params.Search != "" {
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// IsActive filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Valid sort fields for list response
	validSortFields := map[string]bool{
		"id":         true,
		"nama_id":    true,
		"logo_url":   true,
		"is_active":  true,
		"updated_at": true,
	}

	// Validate sort_by field
	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "nama_id" // Default sort field
	}

	// Validate order direction
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc" // Default order
	}

	orderClause := sortBy + " " + order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Scan(&mereks).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	responses := make([]models.MerekProdukSimpleResponse, len(mereks))
	for i, m := range mereks {
		responses[i] = models.MerekProdukSimpleResponse{
			ID: m.ID,
			Nama: models.TranslatableString{
				ID: m.NamaID,
				EN: m.NamaEN,
			},
			LogoURL:   m.LogoURL,
			IsActive:  m.IsActive,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return responses, total, nil
}

func (r *merekProdukRepository) Update(ctx context.Context, merek *models.MerekProduk) error {
	return r.db.WithContext(ctx).Save(merek).Error
}

func (r *merekProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MerekProduk{}).Error
}

func (r *merekProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.MerekProduk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *merekProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.MerekProduk, error) {
	var mereks []models.MerekProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("nama ASC").
		Find(&mereks).Error
	return mereks, err
}
