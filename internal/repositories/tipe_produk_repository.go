package repositories

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

// TipeProdukRepository interface for tipe produk operations (Read-only)
// Note: Tipe produk data is managed via migration only (Paletbox, Container, Truckload)
type TipeProdukRepository interface {
	FindAll(ctx context.Context) ([]dto.TipeProdukListDTO, error)
	FindAllWithProduk(ctx context.Context) ([]dto.TipeProdukWithProdukDTO, error)
	FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error)
	GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error)
}

type tipeProdukRepository struct {
	db *gorm.DB
}

func NewTipeProdukRepository(db *gorm.DB) TipeProdukRepository {
	return &tipeProdukRepository{db: db}
}

// FindAll retrieves all tipe produk without pagination
// Returns TipeProdukListDTO ordered by urutan
func (r *tipeProdukRepository) FindAll(ctx context.Context) ([]dto.TipeProdukListDTO, error) {
	var tipes []dto.TipeProdukListDTO

	err := r.db.WithContext(ctx).
		Model(&models.TipeProduk{}).
		Select(`
			id, 
			nama, 
			slug,
			urutan, 
			is_active, 
			updated_at
		`).
		Where("deleted_at IS NULL").
		Order("urutan ASC").
		Scan(&tipes).Error

	if err != nil {
		return nil, err
	}

	return tipes, nil
}

// FindAllWithProduk retrieves all tipe produk with their products
func (r *tipeProdukRepository) FindAllWithProduk(ctx context.Context) ([]dto.TipeProdukWithProdukDTO, error) {
	var result []dto.TipeProdukWithProdukDTO

	// Get all tipe produk
	var tipes []models.TipeProduk
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("urutan ASC").
		Find(&tipes).Error

	if err != nil {
		return nil, err
	}

	// For each tipe produk, get its products
	for _, tipe := range tipes {
		var produk []dto.ProdukBasicDTO

		err := r.db.WithContext(ctx).
			Model(&models.Produk{}).
			Select(`
				id,
				nama,
				slug,
				harga_sebelum_diskon,
				persentase_diskon,
				harga_sesudah_diskon,
				quantity,
				is_active
			`).
			Where("tipe_produk_id = ? AND deleted_at IS NULL", tipe.ID).
			Order("created_at ASC").
			Scan(&produk).Error

		if err != nil {
			return nil, err
		}

		// Ensure empty array instead of null
		if produk == nil {
			produk = []dto.ProdukBasicDTO{}
		}

		result = append(result, dto.TipeProdukWithProdukDTO{
			ID:   tipe.ID,
			Nama: tipe.Nama,
			// Slug:      tipe.Slug,
			// Deskripsi: tipe.Deskripsi,
			Urutan: tipe.Urutan,
			// IsActive: tipe.IsActive,
			Produk: produk,
		})
	}

	return result, nil
}

// FindBySlug retrieves tipe produk by slug
func (r *tipeProdukRepository) FindBySlug(ctx context.Context, slug string) (*models.TipeProduk, error) {
	var tipe models.TipeProduk
	err := r.db.WithContext(ctx).
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&tipe).Error
	if err != nil {
		return nil, err
	}
	return &tipe, nil
}

// GetAllForDropdown retrieves all active tipe produk for dropdown selection
func (r *tipeProdukRepository) GetAllForDropdown(ctx context.Context) ([]models.TipeProduk, error) {
	var tipes []models.TipeProduk
	err := r.db.WithContext(ctx).
		Select("id", "nama", "slug").
		Where("is_active = ?", true).
		Order("urutan ASC").
		Find(&tipes).Error
	return tipes, err
}
