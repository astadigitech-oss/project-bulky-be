package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KondisiProdukRepository interface {
	Create(ctx context.Context, kondisi *models.KondisiProduk) error
	FindByID(ctx context.Context, id string) (*models.KondisiProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.KondisiProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiProduk, int64, error)
	Update(ctx context.Context, kondisi *models.KondisiProduk) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.KondisiProduk, error)
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
}

type kondisiProdukRepository struct {
	db *gorm.DB
}

func NewKondisiProdukRepository(db *gorm.DB) KondisiProdukRepository {
	return &kondisiProdukRepository{db: db}
}

func (r *kondisiProdukRepository) Create(ctx context.Context, kondisi *models.KondisiProduk) error {
	return r.db.WithContext(ctx).Create(kondisi).Error
}

func (r *kondisiProdukRepository) FindByID(ctx context.Context, id string) (*models.KondisiProduk, error) {
	var kondisi models.KondisiProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&kondisi).Error
	if err != nil {
		return nil, err
	}
	return &kondisi, nil
}

func (r *kondisiProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.KondisiProduk, error) {
	var kondisi models.KondisiProduk
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&kondisi).Error
	if err != nil {
		return nil, err
	}
	return &kondisi, nil
}

func (r *kondisiProdukRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiProduk, int64, error) {
	var kondisis []models.KondisiProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KondisiProduk{})

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

	if err := query.Find(&kondisis).Error; err != nil {
		return nil, 0, err
	}

	return kondisis, total, nil
}

func (r *kondisiProdukRepository) Update(ctx context.Context, kondisi *models.KondisiProduk) error {
	return r.db.WithContext(ctx).Save(kondisi).Error
}

func (r *kondisiProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.KondisiProduk{}).Error
}

func (r *kondisiProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.KondisiProduk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kondisiProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.KondisiProduk, error) {
	var kondisis []models.KondisiProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&kondisis).Error
	return kondisis, err
}

func (r *kondisiProdukRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.KondisiProduk{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
