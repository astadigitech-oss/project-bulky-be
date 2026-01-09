package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

// TipeProdukRepository interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukRepository interface {
	FindByID(ctx context.Context, id string) (*models.TipeProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProduk, int64, error)
	GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error)
}

type tipeProdukRepository struct {
	db *gorm.DB
}

func NewTipeProdukRepository(db *gorm.DB) TipeProdukRepository {
	return &tipeProdukRepository{db: db}
}

// FindByID retrieves a single tipe produk by ID
func (r *tipeProdukRepository) FindByID(ctx context.Context, id string) (*models.TipeProduk, error) {
	var tipe models.TipeProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tipe).Error
	if err != nil {
		return nil, err
	}
	return &tipe, nil
}

// FindBySlug retrieves a single tipe produk by slug
func (r *tipeProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error) {
	var tipe models.TipeProduk
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tipe).Error
	if err != nil {
		return nil, err
	}
	return &tipe, nil
}

// FindAll retrieves all tipe produk with pagination
func (r *tipeProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProduk, int64, error) {
	var tipes []models.TipeProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TipeProduk{})

	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&tipes).Error; err != nil {
		return nil, 0, err
	}

	return tipes, total, nil
}

// GetAllForDropdown retrieves all active tipe produk for dropdown selection
func (r *tipeProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error) {
	var tipes []models.TipeProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&tipes).Error
	return tipes, err
}
