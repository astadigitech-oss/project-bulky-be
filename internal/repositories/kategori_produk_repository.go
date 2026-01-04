package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KategoriProdukRepository interface {
	Create(ctx context.Context, kategori *models.KategoriProduk) error
	FindByID(ctx context.Context, id string) (*models.KategoriProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.KategoriProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KategoriProduk, int64, error)
	Update(ctx context.Context, kategori *models.KategoriProduk) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.KategoriProduk, error)
}

type kategoriProdukRepository struct {
	db *gorm.DB
}

func NewKategoriProdukRepository(db *gorm.DB) KategoriProdukRepository {
	return &kategoriProdukRepository{db: db}
}

func (r *kategoriProdukRepository) Create(ctx context.Context, kategori *models.KategoriProduk) error {
	return r.db.WithContext(ctx).Create(kategori).Error
}

func (r *kategoriProdukRepository) FindByID(ctx context.Context, id string) (*models.KategoriProduk, error) {
	var kategori models.KategoriProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&kategori).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.KategoriProduk, error) {
	var kategori models.KategoriProduk
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&kategori).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KategoriProduk, int64, error) {
	var kategoris []models.KategoriProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KategoriProduk{})

	// Search filter
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	// Active filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderClause := params.SortBy + " " + params.Order
	query = query.Order(orderClause)

	// Pagination
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&kategoris).Error; err != nil {
		return nil, 0, err
	}

	return kategoris, total, nil
}

func (r *kategoriProdukRepository) Update(ctx context.Context, kategori *models.KategoriProduk) error {
	return r.db.WithContext(ctx).Save(kategori).Error
}

func (r *kategoriProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.KategoriProduk{}).Error
}

func (r *kategoriProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.KategoriProduk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kategoriProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.KategoriProduk, error) {
	var kategoris []models.KategoriProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("nama ASC").
		Find(&kategoris).Error
	return kategoris, err
}
