package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KecamatanRepository interface {
	Create(ctx context.Context, kecamatan *models.Kecamatan) error
	FindByID(ctx context.Context, id string) (*models.Kecamatan, error)
	FindAll(ctx context.Context, params *models.KecamatanFilterRequest) ([]models.Kecamatan, int64, error)
	FindByKotaID(ctx context.Context, kotaID string) ([]models.Kecamatan, error)
	Update(ctx context.Context, kecamatan *models.Kecamatan) error
	Delete(ctx context.Context, id string) error
	ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error)
	CountKelurahan(ctx context.Context, kecamatanID string) (int64, error)
}

type kecamatanRepository struct {
	db *gorm.DB
}

func NewKecamatanRepository(db *gorm.DB) KecamatanRepository {
	return &kecamatanRepository{db: db}
}

func (r *kecamatanRepository) Create(ctx context.Context, kecamatan *models.Kecamatan) error {
	return r.db.WithContext(ctx).Create(kecamatan).Error
}

func (r *kecamatanRepository) FindByID(ctx context.Context, id string) (*models.Kecamatan, error) {
	var kecamatan models.Kecamatan
	err := r.db.WithContext(ctx).
		Preload("Kota").
		Preload("Kota.Provinsi").
		Where("id = ?", id).First(&kecamatan).Error
	return &kecamatan, err
}

func (r *kecamatanRepository) FindAll(ctx context.Context, params *models.KecamatanFilterRequest) ([]models.Kecamatan, int64, error) {
	var kecamatans []models.Kecamatan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Kecamatan{}).
		Preload("Kota").
		Preload("Kota.Provinsi")

	if params.KotaID != "" {
		query = query.Where("kota_id = ?", params.KotaID)
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
		Find(&kecamatans).Error

	return kecamatans, total, err
}

func (r *kecamatanRepository) FindByKotaID(ctx context.Context, kotaID string) ([]models.Kecamatan, error) {
	var kecamatans []models.Kecamatan
	err := r.db.WithContext(ctx).Where("kota_id = ?", kotaID).Order("nama ASC").Find(&kecamatans).Error
	return kecamatans, err
}

func (r *kecamatanRepository) Update(ctx context.Context, kecamatan *models.Kecamatan) error {
	return r.db.WithContext(ctx).Save(kecamatan).Error
}

func (r *kecamatanRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Kecamatan{}, "id = ?", id).Error
}

func (r *kecamatanRepository) ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Kecamatan{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kecamatanRepository) CountKelurahan(ctx context.Context, kecamatanID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Kelurahan{}).Where("kecamatan_id = ?", kecamatanID).Count(&count).Error
	return count, err
}
