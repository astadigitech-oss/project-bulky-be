package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KondisiPaketRepository interface {
	Create(ctx context.Context, kondisi *models.KondisiPaket) error
	FindByID(ctx context.Context, id string) (*models.KondisiPaket, error)
	FindBySlug(ctx context.Context, slug string) (*models.KondisiPaket, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiPaketSimpleResponse, int64, error)
	Update(ctx context.Context, kondisi *models.KondisiPaket) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.KondisiPaket, error)
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
}

type kondisiPaketRepository struct {
	db *gorm.DB
}

func NewKondisiPaketRepository(db *gorm.DB) KondisiPaketRepository {
	return &kondisiPaketRepository{db: db}
}

func (r *kondisiPaketRepository) Create(ctx context.Context, kondisi *models.KondisiPaket) error {
	return r.db.WithContext(ctx).Create(kondisi).Error
}

func (r *kondisiPaketRepository) FindByID(ctx context.Context, id string) (*models.KondisiPaket, error) {
	var kondisi models.KondisiPaket
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&kondisi).Error
	if err != nil {
		return nil, err
	}
	return &kondisi, nil
}

func (r *kondisiPaketRepository) FindBySlug(ctx context.Context, slug string) (*models.KondisiPaket, error) {
	var kondisi models.KondisiPaket
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&kondisi).Error
	if err != nil {
		return nil, err
	}
	return &kondisi, nil
}

func (r *kondisiPaketRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.KondisiPaketSimpleResponse, int64, error) {
	var kondisis []models.KondisiPaketSimpleResponse
	var total int64

	// query := r.db.WithContext(ctx).Model(&models.KondisiPaket{})
	query := r.db.WithContext(ctx).Model(&models.KondisiPaket{}).
		Select(`
			id,
			nama,
			urutan,
			is_active,
			updated_at
		`)

	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"nama":       true,
		"urutan":     true,
		"is_active":  true,
		"updated_at": true,
	}

	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "nama"
	}

	// Validate order direction
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc" // Default order
	}

	orderClause := sortBy + " " + order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Scan(&kondisis).Error; err != nil {
		return nil, 0, err
	}

	return kondisis, total, nil
}

func (r *kondisiPaketRepository) Update(ctx context.Context, kondisi *models.KondisiPaket) error {
	return r.db.WithContext(ctx).Save(kondisi).Error
}

func (r *kondisiPaketRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.KondisiPaket{}).Error
}

func (r *kondisiPaketRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.KondisiPaket{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kondisiPaketRepository) GetAllForDropdown(ctx context.Context) ([]models.KondisiPaket, error) {
	var kondisis []models.KondisiPaket
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&kondisis).Error
	return kondisis, err
}

func (r *kondisiPaketRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.KondisiPaket{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
