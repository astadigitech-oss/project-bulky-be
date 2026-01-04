package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type AlamatBuyerRepository interface {
	Create(ctx context.Context, alamat *models.AlamatBuyer) error
	FindByID(ctx context.Context, id string) (*models.AlamatBuyer, error)
	FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyer, int64, error)
	FindByBuyerID(ctx context.Context, buyerID string) ([]models.AlamatBuyer, error)
	Update(ctx context.Context, alamat *models.AlamatBuyer) error
	Delete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, id, buyerID string) error
	CountByBuyerID(ctx context.Context, buyerID string) (int64, error)
}

type alamatBuyerRepository struct {
	db *gorm.DB
}

func NewAlamatBuyerRepository(db *gorm.DB) AlamatBuyerRepository {
	return &alamatBuyerRepository{db: db}
}

func (r *alamatBuyerRepository) Create(ctx context.Context, alamat *models.AlamatBuyer) error {
	return r.db.WithContext(ctx).Create(alamat).Error
}

func (r *alamatBuyerRepository) FindByID(ctx context.Context, id string) (*models.AlamatBuyer, error) {
	var alamat models.AlamatBuyer
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&alamat).Error
	return &alamat, err
}

func (r *alamatBuyerRepository) FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyer, int64, error) {
	var alamat []models.AlamatBuyer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AlamatBuyer{})

	// Filter by buyer_id (required)
	query = query.Where("buyer_id = ?", params.BuyerID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination & ordering
	offset := params.GetOffset()
	query = query.Order("is_default DESC, created_at DESC").
		Offset(offset).
		Limit(params.PerPage)

	if err := query.Find(&alamat).Error; err != nil {
		return nil, 0, err
	}

	return alamat, total, nil
}

func (r *alamatBuyerRepository) FindByBuyerID(ctx context.Context, buyerID string) ([]models.AlamatBuyer, error) {
	var alamat []models.AlamatBuyer
	err := r.db.WithContext(ctx).
		Where("buyer_id = ?", buyerID).
		Order("is_default DESC, created_at DESC").
		Find(&alamat).Error
	return alamat, err
}

func (r *alamatBuyerRepository) Update(ctx context.Context, alamat *models.AlamatBuyer) error {
	return r.db.WithContext(ctx).Save(alamat).Error
}

func (r *alamatBuyerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AlamatBuyer{}, "id = ?", id).Error
}

func (r *alamatBuyerRepository) SetDefault(ctx context.Context, id, buyerID string) error {
	// Trigger di database akan handle unset default lainnya
	return r.db.WithContext(ctx).
		Model(&models.AlamatBuyer{}).
		Where("id = ?", id).
		Update("is_default", true).Error
}

func (r *alamatBuyerRepository) CountByBuyerID(ctx context.Context, buyerID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.AlamatBuyer{}).
		Where("buyer_id = ?", buyerID).
		Count(&count).Error
	return count, err
}
