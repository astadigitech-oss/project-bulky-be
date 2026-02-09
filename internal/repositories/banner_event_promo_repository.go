package repositories

import (
	"context"
	"fmt"
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerEventPromoRepository interface {
	Create(ctx context.Context, banner *models.BannerEventPromo) error
	CreateWithKategori(ctx context.Context, banner *models.BannerEventPromo, kategoriIDs []uuid.UUID) error
	FindByID(ctx context.Context, id string) (*models.BannerEventPromo, error)
	FindAll(ctx context.Context, params *models.BannerEventPromoFilterRequest) ([]models.BannerEventPromo, int64, error)
	Update(ctx context.Context, banner *models.BannerEventPromo) error
	UpdateWithKategori(ctx context.Context, banner *models.BannerEventPromo, kategoriIDs []uuid.UUID) error
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

func (r *bannerEventPromoRepository) CreateWithKategori(ctx context.Context, banner *models.BannerEventPromo, kategoriIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create banner
		if err := tx.Create(banner).Error; err != nil {
			return err
		}

		// Create pivot table entries
		for _, kategoriID := range kategoriIDs {
			relation := &models.BannerEPKategoriProduk{
				BannerID:         banner.ID,
				KategoriProdukID: kategoriID,
			}
			if err := tx.Create(relation).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bannerEventPromoRepository) FindByID(ctx context.Context, id string) (*models.BannerEventPromo, error) {
	var banner models.BannerEventPromo
	err := r.db.WithContext(ctx).
		Preload("KategoriList").
		Where("id = ?", id).
		First(&banner).Error
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

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination & ordering with preload
	offset := params.GetOffset()
	orderBy := params.SortBy + " " + params.Order
	query = query.Preload("KategoriList").
		Order(orderBy).
		Offset(offset).
		Limit(params.PerPage)

	if err := query.Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

func (r *bannerEventPromoRepository) Update(ctx context.Context, banner *models.BannerEventPromo) error {
	return r.db.WithContext(ctx).Save(banner).Error
}

func (r *bannerEventPromoRepository) UpdateWithKategori(ctx context.Context, banner *models.BannerEventPromo, kategoriIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update banner
		if err := tx.Save(banner).Error; err != nil {
			return err
		}

		// Delete existing pivot relations
		if err := tx.Where("banner_id = ?", banner.ID).
			Delete(&models.BannerEPKategoriProduk{}).Error; err != nil {
			return err
		}

		// Create new pivot relations
		for _, kategoriID := range kategoriIDs {
			relation := &models.BannerEPKategoriProduk{
				BannerID:         banner.ID,
				KategoriProdukID: kategoriID,
			}
			if err := tx.Create(relation).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
		Preload("KategoriList").
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
