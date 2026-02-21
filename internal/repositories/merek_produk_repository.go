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
	FindAll(ctx context.Context, params *models.MerekProdukFilterRequest) ([]models.MerekProdukSimpleResponse, int64, error)
	Update(ctx context.Context, merek *models.MerekProduk) error
	Delete(ctx context.Context, merek *models.MerekProduk) error
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
	err := r.db.WithContext(ctx).Where("slug_id = ? OR slug_en = ?", slug, slug).First(&merek).Error
	if err != nil {
		return nil, err
	}
	return &merek, nil
}

func (r *merekProdukRepository) FindAll(ctx context.Context, params *models.MerekProdukFilterRequest) ([]models.MerekProdukSimpleResponse, int64, error) {
	type row struct {
		ID        string    `json:"id"`
		NamaID    string    `json:"nama_id"`
		NamaEN    *string   `json:"nama_en"`
		LogoURL   *string   `json:"logo_url"`
		IsActive  bool      `json:"is_active"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	var mereks []row
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.MerekProduk{}).
		Select("id, nama_id, nama_en, logo_url, is_active, updated_at")

	// Search dengan ILIKE pada nama_id dan nama_en
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", search, search)
	}

	// IsActive filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// sort_by: is_active atau updated_at
	sortColumn := "updated_at"
	if params.SortBy == "is_active" {
		sortColumn = "is_active"
	}
	orderDir := "DESC"
	if params.Order == "asc" {
		orderDir = "ASC"
	}
	query = query.Order(sortColumn + " " + orderDir)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Scan(&mereks).Error; err != nil {
		return nil, 0, err
	}

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

func (r *merekProdukRepository) Delete(ctx context.Context, merek *models.MerekProduk) error {
	return r.db.WithContext(ctx).Delete(merek).Error
}

func (r *merekProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.MerekProduk{}).Where("slug_id = ? OR slug_en = ?", slug, slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *merekProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.MerekProduk, error) {
	var mereks []models.MerekProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama_id", "nama_en").
		Where("is_active = ?", true).
		Order("nama_id ASC").
		Find(&mereks).Error
	return mereks, err
}
