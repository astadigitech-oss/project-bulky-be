package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type TipeProdukRepository interface {
	Create(ctx context.Context, tipe *models.TipeProduk) error
	FindByID(ctx context.Context, id string) (*models.TipeProduk, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.TipeProduk, int64, error)
	Update(ctx context.Context, tipe *models.TipeProduk) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error)
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
}

type tipeProdukRepository struct {
	db *gorm.DB
}

func NewTipeProdukRepository(db *gorm.DB) TipeProdukRepository {
	return &tipeProdukRepository{db: db}
}

func (r *tipeProdukRepository) Create(ctx context.Context, tipe *models.TipeProduk) error {
	return r.db.WithContext(ctx).Create(tipe).Error
}

func (r *tipeProdukRepository) FindByID(ctx context.Context, id string) (*models.TipeProduk, error) {
	var tipe models.TipeProduk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tipe).Error
	if err != nil {
		return nil, err
	}
	return &tipe, nil
}

func (r *tipeProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error) {
	var tipe models.TipeProduk
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tipe).Error
	if err != nil {
		return nil, err
	}
	return &tipe, nil
}

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

func (r *tipeProdukRepository) Update(ctx context.Context, tipe *models.TipeProduk) error {
	return r.db.WithContext(ctx).Save(tipe).Error
}

func (r *tipeProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.TipeProduk{}).Error
}

func (r *tipeProdukRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.TipeProduk{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *tipeProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error) {
	var tipes []models.TipeProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&tipes).Error
	return tipes, err
}

func (r *tipeProdukRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.TipeProduk{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
