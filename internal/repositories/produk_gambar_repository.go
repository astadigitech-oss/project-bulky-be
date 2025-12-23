package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ProdukGambarRepository interface {
	Create(ctx context.Context, gambar *models.ProdukGambar) error
	FindByID(ctx context.Context, id string) (*models.ProdukGambar, error)
	FindByProdukID(ctx context.Context, produkID string) ([]models.ProdukGambar, error)
	Update(ctx context.Context, gambar *models.ProdukGambar) error
	Delete(ctx context.Context, id string) error
	CountByProdukID(ctx context.Context, produkID string) (int64, error)
	UpdateOrder(ctx context.Context, items []models.ReorderItem) error
	SetPrimary(ctx context.Context, produkID, gambarID string) error
}

type produkGambarRepository struct {
	db *gorm.DB
}

func NewProdukGambarRepository(db *gorm.DB) ProdukGambarRepository {
	return &produkGambarRepository{db: db}
}

func (r *produkGambarRepository) Create(ctx context.Context, gambar *models.ProdukGambar) error {
	return r.db.WithContext(ctx).Create(gambar).Error
}

func (r *produkGambarRepository) FindByID(ctx context.Context, id string) (*models.ProdukGambar, error) {
	var gambar models.ProdukGambar
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&gambar).Error
	if err != nil {
		return nil, err
	}
	return &gambar, nil
}

func (r *produkGambarRepository) FindByProdukID(ctx context.Context, produkID string) ([]models.ProdukGambar, error) {
	var gambars []models.ProdukGambar
	err := r.db.WithContext(ctx).
		Where("produk_id = ?", produkID).
		Order("urutan ASC").
		Find(&gambars).Error
	return gambars, err
}

func (r *produkGambarRepository) Update(ctx context.Context, gambar *models.ProdukGambar) error {
	return r.db.WithContext(ctx).Save(gambar).Error
}

func (r *produkGambarRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ProdukGambar{}).Error
}

func (r *produkGambarRepository) CountByProdukID(ctx context.Context, produkID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProdukGambar{}).Where("produk_id = ?", produkID).Count(&count).Error
	return count, err
}

func (r *produkGambarRepository) UpdateOrder(ctx context.Context, items []models.ReorderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&models.ProdukGambar{}).
				Where("id = ?", item.ID).
				Update("urutan", item.Urutan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *produkGambarRepository) SetPrimary(ctx context.Context, produkID, gambarID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset all to non-primary
		if err := tx.Model(&models.ProdukGambar{}).
			Where("produk_id = ?", produkID).
			Update("is_primary", false).Error; err != nil {
			return err
		}
		// Set the selected one as primary
		return tx.Model(&models.ProdukGambar{}).
			Where("id = ?", gambarID).
			Update("is_primary", true).Error
	})
}
