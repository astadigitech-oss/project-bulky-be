package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ProdukDokumenRepository interface {
	Create(ctx context.Context, dokumen *models.ProdukDokumen) error
	FindByID(ctx context.Context, id string) (*models.ProdukDokumen, error)
	FindByProdukID(ctx context.Context, produkID string) ([]models.ProdukDokumen, error)
	Delete(ctx context.Context, id string) error
}

type produkDokumenRepository struct {
	db *gorm.DB
}

func NewProdukDokumenRepository(db *gorm.DB) ProdukDokumenRepository {
	return &produkDokumenRepository{db: db}
}

func (r *produkDokumenRepository) Create(ctx context.Context, dokumen *models.ProdukDokumen) error {
	return r.db.WithContext(ctx).Create(dokumen).Error
}

func (r *produkDokumenRepository) FindByID(ctx context.Context, id string) (*models.ProdukDokumen, error) {
	var dokumen models.ProdukDokumen
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dokumen).Error
	if err != nil {
		return nil, err
	}
	return &dokumen, nil
}

func (r *produkDokumenRepository) FindByProdukID(ctx context.Context, produkID string) ([]models.ProdukDokumen, error) {
	var dokumens []models.ProdukDokumen
	err := r.db.WithContext(ctx).
		Where("produk_id = ?", produkID).
		Order("created_at ASC").
		Find(&dokumens).Error
	return dokumens, err
}

func (r *produkDokumenRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ProdukDokumen{}).Error
}
