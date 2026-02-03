package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type BannerTipeProdukRepository interface {
	Create(ctx context.Context, banner *models.BannerTipeProduk) error
	FindByID(ctx context.Context, id string) (*models.BannerTipeProduk, error)
	FindAll(ctx context.Context, params *models.BannerTipeProdukFilterRequest, tipeProdukID string) ([]models.BannerTipeProduk, int64, error)
	FindAllGrouped(ctx context.Context, search string) ([]models.BannerTipeProduk, error)
	FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerTipeProduk, error)
	Update(ctx context.Context, banner *models.BannerTipeProduk) error
	Delete(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
	GetMaxUrutan(ctx context.Context) (int, error)
	GetMaxUrutanByTipeProduk(ctx context.Context, tipeProdukID string) (int, error)
	FindPreviousByUrutan(ctx context.Context, currentUrutan int) (*models.BannerTipeProduk, error)
	FindNextByUrutan(ctx context.Context, currentUrutan int) (*models.BannerTipeProduk, error)
	FindPreviousByUrutanAndTipeProduk(ctx context.Context, tipeProdukID string, currentUrutan int) (*models.BannerTipeProduk, error)
	FindNextByUrutanAndTipeProduk(ctx context.Context, tipeProdukID string, currentUrutan int) (*models.BannerTipeProduk, error)
	ReorderAfterDeleteScoped(ctx context.Context, tipeProdukID string, deletedUrutan int) error
}

type bannerTipeProdukRepository struct {
	db *gorm.DB
}

func NewBannerTipeProdukRepository(db *gorm.DB) BannerTipeProdukRepository {
	return &bannerTipeProdukRepository{db: db}
}

func (r *bannerTipeProdukRepository) Create(ctx context.Context, banner *models.BannerTipeProduk) error {
	return r.db.WithContext(ctx).Create(banner).Error
}

func (r *bannerTipeProdukRepository) FindByID(ctx context.Context, id string) (*models.BannerTipeProduk, error) {
	var banner models.BannerTipeProduk
	err := r.db.WithContext(ctx).Preload("TipeProduk").Where("id = ?", id).First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

func (r *bannerTipeProdukRepository) FindAll(ctx context.Context, params *models.BannerTipeProdukFilterRequest, tipeProdukID string) ([]models.BannerTipeProduk, int64, error) {
	var banners []models.BannerTipeProduk
	var total int64

	query := r.db.WithContext(ctx).Model(&models.BannerTipeProduk{}).Preload("TipeProduk")

	if tipeProdukID != "" {
		query = query.Where("tipe_produk_id = ?", tipeProdukID)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.Search != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := params.SortBy + " " + params.Order
	query = query.Order(orderClause)
	query = query.Offset(params.GetOffset()).Limit(params.PerPage)

	if err := query.Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

// FindAllGrouped - Get all banners without pagination, for grouped response
func (r *bannerTipeProdukRepository) FindAllGrouped(ctx context.Context, search string) ([]models.BannerTipeProduk, error) {
	var banners []models.BannerTipeProduk

	query := r.db.WithContext(ctx).Preload("TipeProduk")

	// Search filter only
	if search != "" {
		query = query.Where("nama ILIKE ?", "%"+search+"%")
	}

	// Order by tipe_produk slug, then urutan
	query = query.Joins("JOIN tipe_produk ON tipe_produk.id = banner_tipe_produk.tipe_produk_id").
		Order("tipe_produk.slug ASC, banner_tipe_produk.urutan ASC")

	if err := query.Find(&banners).Error; err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *bannerTipeProdukRepository) FindByTipeProdukID(ctx context.Context, tipeProdukID string) ([]models.BannerTipeProduk, error) {
	var banners []models.BannerTipeProduk
	err := r.db.WithContext(ctx).
		Where("tipe_produk_id = ?", tipeProdukID).
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&banners).Error
	return banners, err
}

func (r *bannerTipeProdukRepository) Update(ctx context.Context, banner *models.BannerTipeProduk) error {
	return r.db.WithContext(ctx).Omit("TipeProduk").Save(banner).Error
}

func (r *bannerTipeProdukRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BannerTipeProduk{}).Error
}

func (r *bannerTipeProdukRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// Check if banner exists first
			var count int64
			if err := tx.Model(&models.BannerTipeProduk{}).
				Where("id = ?", item.ID).
				Count(&count).Error; err != nil {
				return err
			}

			if count == 0 {
				return gorm.ErrRecordNotFound
			}

			// Update urutan
			result := tx.Model(&models.BannerTipeProduk{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
}

func (r *bannerTipeProdukRepository) GetMaxUrutan(ctx context.Context) (int, error) {
	var maxUrutan int
	err := r.db.WithContext(ctx).
		Model(&models.BannerTipeProduk{}).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan).Error
	return maxUrutan, err
}

func (r *bannerTipeProdukRepository) GetMaxUrutanByTipeProduk(ctx context.Context, tipeProdukID string) (int, error) {
	var maxUrutan int
	err := r.db.WithContext(ctx).
		Model(&models.BannerTipeProduk{}).
		Where("tipe_produk_id = ?", tipeProdukID).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan).Error
	return maxUrutan, err
}

// FindPreviousByUrutan - Find banner with urutan < current (global, no scope)
func (r *bannerTipeProdukRepository) FindPreviousByUrutan(ctx context.Context, currentUrutan int) (*models.BannerTipeProduk, error) {
	var banner models.BannerTipeProduk
	err := r.db.WithContext(ctx).
		Where("urutan < ?", currentUrutan).
		Order("urutan DESC").
		First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindNextByUrutan - Find banner with urutan > current (global, no scope)
func (r *bannerTipeProdukRepository) FindNextByUrutan(ctx context.Context, currentUrutan int) (*models.BannerTipeProduk, error) {
	var banner models.BannerTipeProduk
	err := r.db.WithContext(ctx).
		Where("urutan > ?", currentUrutan).
		Order("urutan ASC").
		First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindPreviousByUrutanAndTipeProduk - Find previous banner within same tipe_produk
func (r *bannerTipeProdukRepository) FindPreviousByUrutanAndTipeProduk(ctx context.Context, tipeProdukID string, currentUrutan int) (*models.BannerTipeProduk, error) {
	var banner models.BannerTipeProduk
	err := r.db.WithContext(ctx).
		Where("tipe_produk_id = ?", tipeProdukID).
		Where("urutan < ?", currentUrutan).
		Order("urutan DESC").
		First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindNextByUrutanAndTipeProduk - Find next banner within same tipe_produk
func (r *bannerTipeProdukRepository) FindNextByUrutanAndTipeProduk(ctx context.Context, tipeProdukID string, currentUrutan int) (*models.BannerTipeProduk, error) {
	var banner models.BannerTipeProduk
	err := r.db.WithContext(ctx).
		Where("tipe_produk_id = ?", tipeProdukID).
		Where("urutan > ?", currentUrutan).
		Order("urutan ASC").
		First(&banner).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// ReorderAfterDeleteScoped - Reorder banners after delete within same tipe_produk
func (r *bannerTipeProdukRepository) ReorderAfterDeleteScoped(ctx context.Context, tipeProdukID string, deletedUrutan int) error {
	return r.db.WithContext(ctx).
		Model(&models.BannerTipeProduk{}).
		Where("tipe_produk_id = ?", tipeProdukID).
		Where("urutan > ?", deletedUrutan).
		Update("urutan", gorm.Expr("urutan - 1")).Error
}
