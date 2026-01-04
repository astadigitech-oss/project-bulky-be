package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type MerekProdukRepository interface {
	Create(ctx context.Context, merek *models.MerekProduk) error
	FindByID(ctx context.Context, id string) (*models.MerekProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.MerekProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProduk, int64, error)
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

func (r *merekProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.MerekProduk, int64, error) {
	var mereks []models.MerekProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MerekProduk{})

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

	if err := query.Find(&mereks).Error; err != nil {
		return nil, 0, err
	}

	return mereks, total, nil
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
