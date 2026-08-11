package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DelivereeVehicleTypeRepository interface {
	FindByID(ctx context.Context, id string) (*models.DelivereeVehicleType, error)
	FindAll(ctx context.Context, params *models.DelivereeVehicleTypeFilterRequest) ([]models.DelivereeVehicleType, int64, error)
	// FindByIDDeliveree mencari satu kendaraan berdasarkan id_deliveree+environment
	// (termasuk yang is_active=false), dipakai saat Sync untuk membedakan created/updated.
	FindByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error)
	// FindActiveByIDDeliveree mencari satu kendaraan AKTIF berdasarkan id_deliveree+environment.
	// Dipakai saat create booking untuk memvalidasi deliveree_vehicle_type_id yang disimpan
	// storefront saat checkout; return nil jika tidak ada (kendaraan dinonaktifkan).
	FindActiveByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error)
	Update(ctx context.Context, vehicle *models.DelivereeVehicleType) error
	// BulkSetActive mengaktifkan/menonaktifkan banyak kendaraan sekaligus
	// (multi-select dari panel admin). Return jumlah baris yang terupdate.
	BulkSetActive(ctx context.Context, ids []string, isActive bool) (int64, error)
	// Upsert menyimpan/update kendaraan hasil Sync berdasarkan (id_deliveree, environment).
	Upsert(ctx context.Context, vehicle *models.DelivereeVehicleType) error
	// FindActiveByEnvironment mengambil semua kendaraan aktif pada environment tertentu,
	// diurutkan ascending berdasarkan kubikasi_max (dipakai untuk hitung threshold & pemilihan kendaraan).
	FindActiveByEnvironment(ctx context.Context, environment string) ([]models.DelivereeVehicleType, error)
	// DeactivateMissing menonaktifkan kendaraan pada environment tertentu yang id_deliveree-nya
	// TIDAK ada di daftar activeDelivereeIDs (dipakai saat Sync untuk soft-disable kendaraan yang sudah tidak ada di API).
	DeactivateMissing(ctx context.Context, environment string, activeDelivereeIDs []int) (int64, error)
	UpdateThresholds(ctx context.Context, id uuid.UUID, thresholdKubikasi, thresholdBerat float64) error
}

type delivereeVehicleTypeRepository struct {
	db *gorm.DB
}

func NewDelivereeVehicleTypeRepository(db *gorm.DB) DelivereeVehicleTypeRepository {
	return &delivereeVehicleTypeRepository{db: db}
}

func (r *delivereeVehicleTypeRepository) FindByID(ctx context.Context, id string) (*models.DelivereeVehicleType, error) {
	var vehicle models.DelivereeVehicleType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&vehicle).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *delivereeVehicleTypeRepository) FindAll(ctx context.Context, params *models.DelivereeVehicleTypeFilterRequest) ([]models.DelivereeVehicleType, int64, error) {
	var vehicles []models.DelivereeVehicleType
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DelivereeVehicleType{})

	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}
	if params.Environment != nil && *params.Environment != "" {
		query = query.Where("environment = ?", *params.Environment)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"nama":         true,
		"kubikasi_max": true,
		"berat_max":    true,
		"updated_at":   true,
		"created_at":   true,
	}
	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "kubikasi_max"
	}
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query = query.Order("environment asc").Order(sortBy + " " + order)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&vehicles).Error; err != nil {
		return nil, 0, err
	}

	return vehicles, total, nil
}

func (r *delivereeVehicleTypeRepository) FindByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error) {
	var vehicle models.DelivereeVehicleType
	err := r.db.WithContext(ctx).
		Where("id_deliveree = ? AND environment = ?", idDeliveree, environment).
		First(&vehicle).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *delivereeVehicleTypeRepository) FindActiveByIDDeliveree(ctx context.Context, idDeliveree int, environment string) (*models.DelivereeVehicleType, error) {
	var vehicle models.DelivereeVehicleType
	err := r.db.WithContext(ctx).
		Where("id_deliveree = ? AND environment = ? AND is_active = true", idDeliveree, environment).
		First(&vehicle).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *delivereeVehicleTypeRepository) Update(ctx context.Context, vehicle *models.DelivereeVehicleType) error {
	return r.db.WithContext(ctx).Save(vehicle).Error
}

func (r *delivereeVehicleTypeRepository) BulkSetActive(ctx context.Context, ids []string, isActive bool) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&models.DelivereeVehicleType{}).
		Where("id IN ?", ids).
		Update("is_active", isActive)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}


func (r *delivereeVehicleTypeRepository) Upsert(ctx context.Context, vehicle *models.DelivereeVehicleType) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id_deliveree"}, {Name: "environment"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"nama", "kubikasi_max", "berat_max", "cargo_length", "cargo_width", "cargo_height",
			"is_active", "last_synced_at", "updated_at",
		}),
	}).Create(vehicle).Error
}

func (r *delivereeVehicleTypeRepository) FindActiveByEnvironment(ctx context.Context, environment string) ([]models.DelivereeVehicleType, error) {
	var vehicles []models.DelivereeVehicleType
	err := r.db.WithContext(ctx).
		Where("environment = ? AND is_active = true", environment).
		Order("kubikasi_max ASC").
		Find(&vehicles).Error
	return vehicles, err
}

func (r *delivereeVehicleTypeRepository) DeactivateMissing(ctx context.Context, environment string, activeDelivereeIDs []int) (int64, error) {
	query := r.db.WithContext(ctx).Model(&models.DelivereeVehicleType{}).
		Where("environment = ? AND is_active = true", environment)

	if len(activeDelivereeIDs) > 0 {
		query = query.Where("id_deliveree NOT IN ?", activeDelivereeIDs)
	}

	result := query.Update("is_active", false)
	return result.RowsAffected, result.Error
}

func (r *delivereeVehicleTypeRepository) UpdateThresholds(ctx context.Context, id uuid.UUID, thresholdKubikasi, thresholdBerat float64) error {
	return r.db.WithContext(ctx).Model(&models.DelivereeVehicleType{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"threshold_kubikasi": thresholdKubikasi,
			"threshold_berat":    thresholdBerat,
		}).Error
}
