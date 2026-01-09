package repositories

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TipeProdukRepository interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*dto.TipeProdukDetailDTO, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]dto.TipeProdukListDTO, int64, error)
	GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error)
}

type tipeProdukRepository struct {
	db *gorm.DB
}

func NewTipeProdukRepository(db *gorm.DB) TipeProdukRepository {
	return &tipeProdukRepository{db: db}
}

// FindByID retrieves a single tipe produk by ID with complete details
// Returns TipeProdukDetailDTO with all 10 fields including jumlah_produk
func (r *tipeProdukRepository) FindByID(ctx context.Context, id uuid.UUID) (*dto.TipeProdukDetailDTO, error) {
	var tipe dto.TipeProdukDetailDTO

	err := r.db.WithContext(ctx).
		Model(&models.TipeProduk{}).
		Select(`
			id, 
			nama, 
			slug, 
			deskripsi,
			urutan, 
			is_active, 
			created_at,
			updated_at
		`).
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(&tipe).Error

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
// Returns TipeProdukListDTO with simplified 8 fields for list view
func (r *tipeProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]dto.TipeProdukListDTO, int64, error) {
	var tipes []dto.TipeProdukListDTO
	var total int64

	// Base query for select simplified fields
	query := r.db.WithContext(ctx).
		Model(&models.TipeProduk{}).
		Select(`
			id, 
			nama, 
			slug, 
			urutan, 
			is_active, 
			updated_at
		`).
		Where("deleted_at IS NULL")

	// Search filter
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
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
		"id":   true,
		"nama": true,
		// "slug":          true,
		// "icon_url": true,
		"urutan": true,
		// "jumlah_produk": true,
		// "is_active":     true,
		"updated_at": true,
	}

	// Validate sort_by field
	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "urutan" // Default sort field
	}

	// Validate order direction
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc" // Default order
	}

	orderClause := sortBy + " " + order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Scan(&tipes).Error; err != nil {
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
