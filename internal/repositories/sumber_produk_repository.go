package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type SumberProdukRepository interface {
	Create(ctx context.Context, sumber *models.SumberProduk) error
	FindByID(ctx context.Context, id string) (*models.SumberProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.SumberProduk, error)
	FindAll(ctx context.Context, params *models.SumberProdukFilterRequest) ([]models.SumberProduk, int64, error)
	Update(ctx context.Context, sumber *models.SumberProduk) error
	Delete(ctx context.Context, sumber *models.SumberProduk) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.SumberProduk, error)
}

type sumberProdukRepository struct {
	db *gorm.DB
}

func NewSumberProdukRepository(db *gorm.DB) SumberProdukRepository {
	return &sumberProdukRepository{db: db}
}

func (r *sumberProdukRepository) Create(ctx context.Context, sumber *models.SumberProduk) error {
	return r.db.WithContext(ctx).Create(sumber).Error
}

func (r *sumberProdukRepository) FindByID(ctx context.Context, id string) (*models.SumberProduk, error) {
	var sumber models.SumberProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sumber).Error
	if err != nil {
		return nil, err
	}
	return &sumber, nil
}

func (r *sumberProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.SumberProduk, error) {
	var sumber models.SumberProduk
	err := r.db.WithContext(ctx).Where("slug_id = ? OR slug_en = ?", slug, slug).First(&sumber).Error
	if err != nil {
		return nil, err
	}
	return &sumber, nil
}

func (r *sumberProdukRepository) FindAll(ctx context.Context, params *models.SumberProdukFilterRequest) ([]models.SumberProduk, int64, error) {
	var sumbers []models.SumberProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SumberProduk{})

	// Search dengan ILIKE pada nama_id dan nama_en
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", search, search)
	}

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

	if err := query.Find(&sumbers).Error; err != nil {
		return nil, 0, err
	}

	return sumbers, total, nil
}

func (r *sumberProdukRepository) Update(ctx context.Context, sumber *models.SumberProduk) error {
	return r.db.WithContext(ctx).Save(sumber).Error
}

func (r *sumberProdukRepository) Delete(ctx context.Context, sumber *models.SumberProduk) error {
	return r.db.WithContext(ctx).Delete(sumber).Error
}

func (r *sumberProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.SumberProduk{}).Where("slug_id = ? OR slug_en = ?", slug, slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *sumberProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.SumberProduk, error) {
	var sumbers []models.SumberProduk
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Order("nama_id ASC").
		Find(&sumbers).Error
	return sumbers, err
}
