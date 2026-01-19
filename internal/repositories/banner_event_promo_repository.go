package repositories

import (
	"context"
	"fmt"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type BannerEventPromoRepository interface {
	Create(ctx context.Context, banner *models.BannerEventPromo) error
	FindByID(ctx context.Context, id string) (*models.BannerEventPromo, error)
	FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromo, int64, error)
	Update(ctx context.Context, banner *models.BannerEventPromo) error
	Delete(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
	GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromo, error)
	GetMaxUrutan(ctx context.Context) (int, error)
}

type bannerEventPromoRepository struct {
	db *gorm.DB
}

func NewBannerEventPromoRepository(db *gorm.DB) BannerEventPromoRepository {
	return &bannerEventPromoRepository{db: db}
}

func (r *bannerEventPromoRepository) Create(ctx context.Context, banner *models.BannerEventPromo) error {
	return r.db.WithContext(ctx).Create(banner).Error
}

func (r *bannerEventPromoRepository) FindByID(ctx context.Context, id string) (*models.BannerEventPromo, error) {
	var banner models.BannerEventPromo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&banner).Error
	return &banner, err
}

func (r *bannerEventPromoRepository) FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromo, int64, error) {
	var banners []models.BannerEventPromo
	var total int64

	query := r.db.WithContext(ctx).Model(&models.BannerEventPromo{})

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

	if err := query.Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

func (r *bannerEventPromoRepository) Update(ctx context.Context, banner *models.BannerEventPromo) error {
	return r.db.WithContext(ctx).Save(banner).Error
}

func (r *bannerEventPromoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.BannerEventPromo{}, "id = ?", id).Error
}

func (r *bannerEventPromoRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// Check if banner exists first
			var count int64
			if err := tx.Model(&models.BannerEventPromo{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("banner dengan ID %s tidak ditemukan", item.ID)
			}

			// Update urutan
			if err := tx.Model(&models.BannerEventPromo{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *bannerEventPromoRepository) GetVisibleBanners(ctx context.Context) ([]models.BannerEventPromo, error) {
	var banners []models.BannerEventPromo
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("tanggal_mulai IS NULL OR tanggal_mulai <= ?", now).
		Where("tanggal_selesai IS NULL OR tanggal_selesai >= ?", now).
		Order("urutan ASC").
		Find(&banners).Error

	return banners, err
}

func (r *bannerEventPromoRepository) GetMaxUrutan(ctx context.Context) (int, error) {
	var maxUrutan int
	err := r.db.WithContext(ctx).
		Model(&models.BannerEventPromo{}).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan).Error
	return maxUrutan, err
}
