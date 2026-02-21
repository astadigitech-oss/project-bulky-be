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
	FindAll(ctx context.Context, params *models.KategoriProdukFilterRequest) ([]models.KategoriProduk, int64, error)
	Update(ctx context.Context, kategori *models.KategoriProduk) error
	Delete(ctx context.Context, kategori *models.KategoriProduk) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.KategoriProduk, error)
	FindActiveByIDs(ctx context.Context, ids []string) ([]models.KategoriProduk, error)
	FindAllActiveForDropdown(ctx context.Context) ([]models.KategoriProduk, error)
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
	err := r.db.WithContext(ctx).Where("slug_id = ? OR slug_en = ?", slug, slug).First(&kategori).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *kategoriProdukRepository) FindAll(ctx context.Context, params *models.KategoriProdukFilterRequest) ([]models.KategoriProduk, int64, error) {
	var kategoris []models.KategoriProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KategoriProduk{})

	// Search dengan ILIKE pada nama_id dan nama_en
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("nama_id ILIKE ? OR nama_en ILIKE ?", search, search)
	}

	// Active filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
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

func (r *kategoriProdukRepository) Delete(ctx context.Context, kategori *models.KategoriProduk) error {
	return r.db.WithContext(ctx).Delete(kategori).Error
}

func (r *kategoriProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.KategoriProduk{}).Where("slug_id = ? OR slug_en = ?", slug, slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kategoriProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.KategoriProduk, error) {
	var kategoris []models.KategoriProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama_id", "nama_en", "slug").
		Where("is_active = ?", true).
		Order("nama_id ASC").
		Find(&kategoris).Error
	return kategoris, err
}

func (r *kategoriProdukRepository) FindActiveByIDs(ctx context.Context, ids []string) ([]models.KategoriProduk, error) {
	var kategoris []models.KategoriProduk
	err := r.db.WithContext(ctx).
		Where("id IN ? AND is_active = ?", ids, true).
		Find(&kategoris).Error
	return kategoris, err
}

func (r *kategoriProdukRepository) FindAllActiveForDropdown(ctx context.Context) ([]models.KategoriProduk, error) {
	var kategoris []models.KategoriProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama_id", "nama_en", "slug").
		Where("is_active = ?", true).
		Order("created_at ASC, nama_id ASC").
		Find(&kategoris).Error
	return kategoris, err
}
