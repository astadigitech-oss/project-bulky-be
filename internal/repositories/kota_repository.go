package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KotaRepository interface {
	Create(ctx context.Context, kota *models.Kota) error
	FindByID(ctx context.Context, id string) (*models.Kota, error)
	FindAll(ctx context.Context, params *models.KotaFilterRequest) ([]models.Kota, int64, error)
	FindByProvinsiID(ctx context.Context, provinsiID string) ([]models.Kota, error)
	Update(ctx context.Context, kota *models.Kota) error
	Delete(ctx context.Context, id string) error
	ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error)
	CountKecamatan(ctx context.Context, kotaID string) (int64, error)
}

type kotaRepository struct {
	db *gorm.DB
}

func NewKotaRepository(db *gorm.DB) KotaRepository {
	return &kotaRepository{db: db}
}

func (r *kotaRepository) Create(ctx context.Context, kota *models.Kota) error {
	return r.db.WithContext(ctx).Create(kota).Error
}

func (r *kotaRepository) FindByID(ctx context.Context, id string) (*models.Kota, error) {
	var kota models.Kota
	err := r.db.WithContext(ctx).Preload("Provinsi").Where("id = ?", id).First(&kota).Error
	return &kota, err
}

func (r *kotaRepository) FindAll(ctx context.Context, params *models.KotaFilterRequest) ([]models.Kota, int64, error) {
	var kotas []models.Kota
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Kota{}).Preload("Provinsi")

	if params.ProvinsiID != "" {
		query = query.Where("provinsi_id = ?", params.ProvinsiID)
	}

	if params.Cari != "" {
		search := "%" + params.Cari + "%"
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", search, search)
	}

	query.Count(&total)

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	err := query.Order(orderClause).
		Offset(params.GetOffset()).
		Limit(params.PerHalaman).
		Find(&kotas).Error

	return kotas, total, err
}

func (r *kotaRepository) FindByProvinsiID(ctx context.Context, provinsiID string) ([]models.Kota, error) {
	var kotas []models.Kota
	err := r.db.WithContext(ctx).Where("provinsi_id = ?", provinsiID).Order("nama ASC").Find(&kotas).Error
	return kotas, err
}

func (r *kotaRepository) Update(ctx context.Context, kota *models.Kota) error {
	return r.db.WithContext(ctx).Save(kota).Error
}

func (r *kotaRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Kota{}, "id = ?", id).Error
}

func (r *kotaRepository) ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Kota{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kotaRepository) CountKecamatan(ctx context.Context, kotaID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Kecamatan{}).Where("kota_id = ?", kotaID).Count(&count).Error
	return count, err
}
