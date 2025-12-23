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
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.SumberProduk, int64, error)
	Update(ctx context.Context, sumber *models.SumberProduk) error
	Delete(ctx context.Context, id string) error
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
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&sumber).Error
	if err != nil {
		return nil, err
	}
	return &sumber, nil
}

func (r *sumberProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.SumberProduk, int64, error) {
	var sumbers []models.SumberProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SumberProduk{})

	if params.Cari != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Cari+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerHalaman)

	if err := query.Find(&sumbers).Error; err != nil {
		return nil, 0, err
	}

	return sumbers, total, nil
}

func (r *sumberProdukRepository) Update(ctx context.Context, sumber *models.SumberProduk) error {
	return r.db.WithContext(ctx).Save(sumber).Error
}

func (r *sumberProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.SumberProduk{}).Error
}

func (r *sumberProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.SumberProduk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *sumberProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.SumberProduk, error) {
	var sumbers []models.SumberProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("nama ASC").
		Find(&sumbers).Error
	return sumbers, err
}
