package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type AlamatBuyerRepository interface {
	Create(ctx context.Context, alamat *models.AlamatBuyer) error
	FindByID(ctx context.Context, id string) (*models.AlamatBuyer, error)
	FindByIDWithWilayah(ctx context.Context, id string) (*models.AlamatBuyer, error)
	FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyer, int64, error)
	FindByBuyerID(ctx context.Context, buyerID string) ([]models.AlamatBuyer, error)
	Update(ctx context.Context, alamat *models.AlamatBuyer) error
	Delete(ctx context.Context, id string) error
	CountByBuyerID(ctx context.Context, buyerID string) (int64, error)
	HasOtherAddresses(ctx context.Context, buyerID, excludeID string) (bool, error)
	UnsetDefaultByBuyerID(ctx context.Context, buyerID string, excludeID *string) error
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
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&alamat).Error
	return &alamat, err
}

func (r *alamatBuyerRepository) FindByIDWithWilayah(ctx context.Context, id string) (*models.AlamatBuyer, error) {
	var alamat models.AlamatBuyer
	err := r.db.WithContext(ctx).
		Preload("Kelurahan").
		Preload("Kelurahan.Kecamatan").
		Preload("Kelurahan.Kecamatan.Kota").
		Preload("Kelurahan.Kecamatan.Kota.Provinsi").
		Where("id = ?", id).First(&alamat).Error
	return &alamat, err
}

func (r *alamatBuyerRepository) FindAll(ctx context.Context, params *models.AlamatBuyerFilterRequest) ([]models.AlamatBuyer, int64, error) {
	var alamats []models.AlamatBuyer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).
		Preload("Kelurahan").
		Preload("Kelurahan.Kecamatan").
		Preload("Kelurahan.Kecamatan.Kota").
		Preload("Kelurahan.Kecamatan.Kota.Provinsi")

	if params.BuyerID != "" {
		query = query.Where("buyer_id = ?", params.BuyerID)
	}

	query.Count(&total)

	err := query.Order("is_default DESC, created_at DESC").
		Offset(params.GetOffset()).
		Limit(params.PerHalaman).
		Find(&alamats).Error

	return alamats, total, err
}

func (r *alamatBuyerRepository) FindByBuyerID(ctx context.Context, buyerID string) ([]models.AlamatBuyer, error) {
	var alamats []models.AlamatBuyer
	err := r.db.WithContext(ctx).
		Preload("Kelurahan").
		Preload("Kelurahan.Kecamatan").
		Preload("Kelurahan.Kecamatan.Kota").
		Preload("Kelurahan.Kecamatan.Kota.Provinsi").
		Where("buyer_id = ?", buyerID).
		Order("is_default DESC, created_at DESC").
		Find(&alamats).Error
	return alamats, err
}

func (r *alamatBuyerRepository) Update(ctx context.Context, alamat *models.AlamatBuyer) error {
	return r.db.WithContext(ctx).Save(alamat).Error
}

func (r *alamatBuyerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AlamatBuyer{}, "id = ?", id).Error
}

func (r *alamatBuyerRepository) CountByBuyerID(ctx context.Context, buyerID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).Where("buyer_id = ?", buyerID).Count(&count).Error
	return count, err
}

func (r *alamatBuyerRepository) HasOtherAddresses(ctx context.Context, buyerID, excludeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).
		Where("buyer_id = ? AND id != ?", buyerID, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *alamatBuyerRepository) UnsetDefaultByBuyerID(ctx context.Context, buyerID string, excludeID *string) error {
	query := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).Where("buyer_id = ?", buyerID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	return query.Update("is_default", false).Error
}
