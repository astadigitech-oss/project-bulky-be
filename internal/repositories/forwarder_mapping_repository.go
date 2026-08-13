package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ForwarderMappingRepository interface {
	FindCities(ctx context.Context, params *models.ForwarderCityFilterRequest) ([]models.ForwarderCityMapping, int64, error)
	FindSubdistricts(ctx context.Context, params *models.ForwarderSubdistrictFilterRequest) ([]models.ForwarderSubdistrictMapping, int64, error)
	// UpsertCities menyimpan/update mapping kota hasil sync berdasarkan unique
	// constraint kota_pattern. Return (created, updated).
	UpsertCities(ctx context.Context, items []models.ForwarderCityMapping) (int, int, error)
	// UpsertSubdistricts menyimpan/update mapping kecamatan hasil sync berdasarkan
	// unique constraint (kecamatan_pattern, forwarder_city_id). Return (created, updated).
	UpsertSubdistricts(ctx context.Context, items []models.ForwarderSubdistrictMapping) (int, int, error)
}

type forwarderMappingRepository struct {
	db *gorm.DB
}

func NewForwarderMappingRepository(db *gorm.DB) ForwarderMappingRepository {
	return &forwarderMappingRepository{db: db}
}

func (r *forwarderMappingRepository) FindCities(ctx context.Context, params *models.ForwarderCityFilterRequest) ([]models.ForwarderCityMapping, int64, error) {
	var items []models.ForwarderCityMapping
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ForwarderCityMapping{})

	if params.Search != "" {
		query = query.Where("kota_pattern ILIKE ? OR forwarder_city_name ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"kota_pattern":        true,
		"forwarder_city_name": true,
		"forwarder_city_id":   true,
		"updated_at":          true,
		"created_at":          true,
	}
	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "kota_pattern"
	}
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query = query.Order(sortBy + " " + order)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *forwarderMappingRepository) FindSubdistricts(ctx context.Context, params *models.ForwarderSubdistrictFilterRequest) ([]models.ForwarderSubdistrictMapping, int64, error) {
	var items []models.ForwarderSubdistrictMapping
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ForwarderSubdistrictMapping{})

	if params.Search != "" {
		query = query.Where("kecamatan_pattern ILIKE ? OR forwarder_subdistrict_name ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.ForwarderCityID != nil && *params.ForwarderCityID > 0 {
		query = query.Where("forwarder_city_id = ?", *params.ForwarderCityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortFields := map[string]bool{
		"kecamatan_pattern":          true,
		"forwarder_subdistrict_name": true,
		"forwarder_city_id":          true,
		"forwarder_subdistrict_id":   true,
		"updated_at":                 true,
		"created_at":                 true,
	}
	sortBy := params.SortBy
	if !validSortFields[sortBy] {
		sortBy = "kecamatan_pattern"
	}
	order := params.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query = query.Order(sortBy + " " + order)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *forwarderMappingRepository) UpsertCities(ctx context.Context, items []models.ForwarderCityMapping) (int, int, error) {
	if len(items) == 0 {
		return 0, 0, nil
	}

	created, updated := 0, 0
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range items {
			item := &items[i]
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "kota_pattern"}},
				DoUpdates: clause.AssignmentColumns([]string{"forwarder_city_id", "forwarder_city_name", "updated_at"}),
			}).Create(item)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 1 {
				created++
			} else {
				updated++
			}
		}
		return nil
	})
	return created, updated, err
}

func (r *forwarderMappingRepository) UpsertSubdistricts(ctx context.Context, items []models.ForwarderSubdistrictMapping) (int, int, error) {
	if len(items) == 0 {
		return 0, 0, nil
	}

	created, updated := 0, 0
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range items {
			item := &items[i]
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "kecamatan_pattern"}, {Name: "forwarder_city_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"forwarder_subdistrict_id", "forwarder_subdistrict_name", "updated_at"}),
			}).Create(item)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 1 {
				created++
			} else {
				updated++
			}
		}
		return nil
	})
	return created, updated, err
}
