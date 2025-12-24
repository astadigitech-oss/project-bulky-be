package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type KelurahanRepository interface {
	Create(ctx context.Context, kelurahan *models.Kelurahan) error
	FindByID(ctx context.Context, id string) (*models.Kelurahan, error)
	FindByIDWithFullHierarchy(ctx context.Context, id string) (*models.Kelurahan, error)
	FindAll(ctx context.Context, params *models.KelurahanFilterRequest) ([]models.Kelurahan, int64, error)
	FindByKecamatanID(ctx context.Context, kecamatanID string) ([]models.Kelurahan, error)
	Update(ctx context.Context, kelurahan *models.Kelurahan) error
	Delete(ctx context.Context, id string) error
	ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error)
	CountAlamatBuyer(ctx context.Context, kelurahanID string) (int64, error)
}

type kelurahanRepository struct {
	db *gorm.DB
}

func NewKelurahanRepository(db *gorm.DB) KelurahanRepository {
	return &kelurahanRepository{db: db}
}

func (r *kelurahanRepository) Create(ctx context.Context, kelurahan *models.Kelurahan) error {
	return r.db.WithContext(ctx).Create(kelurahan).Error
}

func (r *kelurahanRepository) FindByID(ctx context.Context, id string) (*models.Kelurahan, error) {
	var kelurahan models.Kelurahan
	err := r.db.WithContext(ctx).
		Preload("Kecamatan").
		Preload("Kecamatan.Kota").
		Preload("Kecamatan.Kota.Provinsi").
		Where("id = ?", id).First(&kelurahan).Error
	return &kelurahan, err
}

func (r *kelurahanRepository) FindByIDWithFullHierarchy(ctx context.Context, id string) (*models.Kelurahan, error) {
	var kelurahan models.Kelurahan
	err := r.db.WithContext(ctx).
		Preload("Kecamatan").
		Preload("Kecamatan.Kota").
		Preload("Kecamatan.Kota.Provinsi").
		Where("id = ?", id).First(&kelurahan).Error
	return &kelurahan, err
}

func (r *kelurahanRepository) FindAll(ctx context.Context, params *models.KelurahanFilterRequest) ([]models.Kelurahan, int64, error) {
	var kelurahans []models.Kelurahan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Kelurahan{}).
		Preload("Kecamatan").
		Preload("Kecamatan.Kota").
		Preload("Kecamatan.Kota.Provinsi")

	if params.KecamatanID != "" {
		query = query.Where("kecamatan_id = ?", params.KecamatanID)
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
		Find(&kelurahans).Error

	return kelurahans, total, err
}

func (r *kelurahanRepository) FindByKecamatanID(ctx context.Context, kecamatanID string) ([]models.Kelurahan, error) {
	var kelurahans []models.Kelurahan
	err := r.db.WithContext(ctx).Where("kecamatan_id = ?", kecamatanID).Order("nama ASC").Find(&kelurahans).Error
	return kelurahans, err
}

func (r *kelurahanRepository) Update(ctx context.Context, kelurahan *models.Kelurahan) error {
	return r.db.WithContext(ctx).Save(kelurahan).Error
}

func (r *kelurahanRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Kelurahan{}, "id = ?", id).Error
}

func (r *kelurahanRepository) ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Kelurahan{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *kelurahanRepository) CountAlamatBuyer(ctx context.Context, kelurahanID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).Where("kelurahan_id = ?", kelurahanID).Count(&count).Error
	return count, err
}
