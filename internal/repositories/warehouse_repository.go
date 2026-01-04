package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *models.Warehouse) error
	FindByID(ctx context.Context, id string) (*models.Warehouse, error)
	FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.Warehouse, int64, error)
	Update(ctx context.Context, warehouse *models.Warehouse) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error)
	GetAllForDropdown(ctx context.Context) ([]models.Warehouse, error)
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *warehouseRepository) FindByID(ctx context.Context, id string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindAll(ctx context.Context, params *models.PaginationRequest, kota string) ([]models.Warehouse, int64, error) {
	var warehouses []models.Warehouse
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Warehouse{})

	if params.Search != "" {
		query = query.Where("nama ILIKE ? OR kota ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if kota != "" {
		query = query.Where("kota ILIKE ?", "%"+kota+"%")
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

	if err := query.Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *warehouseRepository) Update(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *warehouseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Warehouse{}).Error
}

func (r *warehouseRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Warehouse{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *warehouseRepository) GetAllForDropdown(ctx context.Context) ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("nama ASC").
		Find(&warehouses).Error
	return warehouses, err
}
