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
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.HeroSection, int64, error)
	Update(ctx context.Context, hero *models.HeroSection) error
	Delete(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
	GetVisibleHero(ctx context.Context) (*models.HeroSection, error)
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

func (r *heroSectionRepository) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.HeroSection, int64, error) {
	var heroes []models.HeroSection
	var total int64

	query := r.db.WithContext(ctx).Model(&models.HeroSection{})

	// Search
	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	// Filter is_active
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
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

func (r *heroSectionRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.HeroSection{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *heroSectionRepository) GetVisibleHero(ctx context.Context) (*models.HeroSection, error) {
	var hero models.HeroSection
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("tanggal_mulai IS NULL OR tanggal_mulai <= ?", now).
		Where("tanggal_selesai IS NULL OR tanggal_selesai >= ?", now).
		Order("urutan ASC").
		First(&hero).Error

	return &hero, err
}
