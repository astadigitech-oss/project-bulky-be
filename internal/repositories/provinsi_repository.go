package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ProvinsiRepository interface {
	Create(ctx context.Context, provinsi *models.Provinsi) error
	FindByID(ctx context.Context, id string) (*models.Provinsi, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Provinsi, int64, error)
	FindAllDropdown(ctx context.Context) ([]models.Provinsi, error)
	Update(ctx context.Context, provinsi *models.Provinsi) error
	Delete(ctx context.Context, id string) error
	ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error)
	CountKota(ctx context.Context, provinsiID string) (int64, error)
}

type provinsiRepository struct {
	db *gorm.DB
}

func NewProvinsiRepository(db *gorm.DB) ProvinsiRepository {
	return &provinsiRepository{db: db}
}

func (r *provinsiRepository) Create(ctx context.Context, provinsi *models.Provinsi) error {
	return r.db.WithContext(ctx).Create(provinsi).Error
}

func (r *provinsiRepository) FindByID(ctx context.Context, id string) (*models.Provinsi, error) {
	var provinsi models.Provinsi
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&provinsi).Error
	return &provinsi, err
}

func (r *provinsiRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.Provinsi, int64, error) {
	var provinsis []models.Provinsi
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Provinsi{})

	if params.Cari != "" {
		search := "%" + params.Cari + "%"
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", search, search)
	}

	query.Count(&total)

	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	err := query.Order(orderClause).
		Offset(params.GetOffset()).
		Limit(params.PerHalaman).
		Find(&provinsis).Error

	return provinsis, total, err
}

func (r *provinsiRepository) FindAllDropdown(ctx context.Context) ([]models.Provinsi, error) {
	var provinsis []models.Provinsi
	err := r.db.WithContext(ctx).Order("nama ASC").Find(&provinsis).Error
	return provinsis, err
}

func (r *provinsiRepository) Update(ctx context.Context, provinsi *models.Provinsi) error {
	return r.db.WithContext(ctx).Save(provinsi).Error
}

func (r *provinsiRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Provinsi{}, "id = ?", id).Error
}

func (r *provinsiRepository) ExistsByKode(ctx context.Context, kode string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Provinsi{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *provinsiRepository) CountKota(ctx context.Context, provinsiID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Kota{}).Where("provinsi_id = ?", provinsiID).Count(&count).Error
	return count, err
}
