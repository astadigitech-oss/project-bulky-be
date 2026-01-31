package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type HeroSectionRepository interface {
	Create(ctx context.Context, hero *models.HeroSection) error
	FindByID(ctx context.Context, id string) (*models.HeroSection, error)
	FindAll(ctx context.Context, params *models.HeroSectionFilterRequest) ([]models.HeroSection, int64, error)
	Update(ctx context.Context, hero *models.HeroSection) error
	Delete(ctx context.Context, id string) error
	GetVisibleHero(ctx context.Context) (*models.HeroSection, error)
	CheckDateRangeOverlap(ctx context.Context, tanggalMulai, tanggalSelesai *time.Time, excludeID *string) (bool, error)
}

type heroSectionRepository struct {
	db *gorm.DB
}

func NewHeroSectionRepository(db *gorm.DB) HeroSectionRepository {
	return &heroSectionRepository{db: db}
}

func (r *heroSectionRepository) Create(ctx context.Context, hero *models.HeroSection) error {
	return r.db.WithContext(ctx).Create(hero).Error
}

func (r *heroSectionRepository) FindByID(ctx context.Context, id string) (*models.HeroSection, error) {
	var hero models.HeroSection
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&hero).Error
	return &hero, err
}

func (r *heroSectionRepository) FindAll(ctx context.Context, params *models.HeroSectionFilterRequest) ([]models.HeroSection, int64, error) {
	var heroes []models.HeroSection
	var total int64

	query := r.db.WithContext(ctx).Model(&models.HeroSection{})

	// Search
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	// Filter is_default
	if params.IsDefault != nil {
		query = query.Where("is_default = ?", *params.IsDefault)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination & ordering
	offset := params.GetOffset()
	orderBy := params.SortBy + " " + params.Order
	query = query.Order(orderBy).Offset(offset).Limit(params.PerPage)

	if err := query.Find(&heroes).Error; err != nil {
		return nil, 0, err
	}

	return heroes, total, nil
}

func (r *heroSectionRepository) Update(ctx context.Context, hero *models.HeroSection) error {
	return r.db.WithContext(ctx).Save(hero).Error
}

func (r *heroSectionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.HeroSection{}, "id = ?", id).Error
}

func (r *heroSectionRepository) GetVisibleHero(ctx context.Context) (*models.HeroSection, error) {
	var hero models.HeroSection
	now := time.Now()

	// Priority: Scheduled hero (dalam rentang tanggal) > Default hero
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where(`
			is_default = true
			OR (
				tanggal_mulai IS NOT NULL 
				AND tanggal_selesai IS NOT NULL
				AND tanggal_mulai <= ? 
				AND tanggal_selesai >= ?
			)
		`, now, now).
		Order("tanggal_mulai DESC NULLS LAST"). // Prioritize scheduled over default
		First(&hero).Error

	return &hero, err
}

func (r *heroSectionRepository) CheckDateRangeOverlap(ctx context.Context, tanggalMulai, tanggalSelesai *time.Time, excludeID *string) (bool, error) {
	// Skip validation if no date range
	if tanggalMulai == nil || tanggalSelesai == nil {
		return false, nil
	}

	query := r.db.WithContext(ctx).Model(&models.HeroSection{}).
		Where("deleted_at IS NULL").
		Where(`
			(tanggal_mulai IS NOT NULL AND tanggal_selesai IS NOT NULL) AND
			(
				(tanggal_mulai <= ? AND tanggal_selesai >= ?) OR
				(tanggal_mulai <= ? AND tanggal_selesai >= ?) OR
				(tanggal_mulai >= ? AND tanggal_selesai <= ?)
			)
		`, tanggalMulai, tanggalMulai, tanggalSelesai, tanggalSelesai, tanggalMulai, tanggalSelesai)

	// Exclude current record if updating
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
