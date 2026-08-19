package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

// WMSCargoPricedRepository mengelola cache lokal cargo WMS yang sudah
// diberi harga jual — jembatan antara WMS dan produk lokal.
type WMSCargoPricedRepository interface {
	Upsert(ctx context.Context, cargo *models.WMSCargoPriced) error
	FindByCargoID(ctx context.Context, cargoID string) (*models.WMSCargoPriced, error)
	// FindByCode mencari cache cargo berdasarkan "code" (kode bisnis, mis.
	// "004/08/2026") — dipakai untuk auto-attach PDF harga saat create/edit
	// produk, karena field Produk.IDCargo menyimpan code, bukan cargo_id (UUID WMS).
	FindByCode(ctx context.Context, code string) (*models.WMSCargoPriced, error)
	ListUnused(ctx context.Context, search string) ([]models.WMSCargoPriced, error)
	MarkUsed(ctx context.Context, cargoID string, produkID string) error
}

type wmsCargoPricedRepository struct {
	db *gorm.DB
}

func NewWMSCargoPricedRepository(db *gorm.DB) WMSCargoPricedRepository {
	return &wmsCargoPricedRepository{db: db}
}

// Upsert menyimpan/memperbarui cache cargo WMS berdasarkan cargo_id.
func (r *wmsCargoPricedRepository) Upsert(ctx context.Context, cargo *models.WMSCargoPriced) error {
	return r.db.WithContext(ctx).
		Where("cargo_id = ?", cargo.CargoID).
		Assign(cargo).
		FirstOrCreate(cargo).Error
}

func (r *wmsCargoPricedRepository) FindByCargoID(ctx context.Context, cargoID string) (*models.WMSCargoPriced, error) {
	var cargo models.WMSCargoPriced
	if err := r.db.WithContext(ctx).Where("cargo_id = ?", cargoID).First(&cargo).Error; err != nil {
		return nil, err
	}
	return &cargo, nil
}

func (r *wmsCargoPricedRepository) FindByCode(ctx context.Context, code string) (*models.WMSCargoPriced, error) {
	var cargo models.WMSCargoPriced
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&cargo).Error; err != nil {
		return nil, err
	}
	return &cargo, nil
}

// ListUnused daftar cargo yang belum dipakai di produk manapun (untuk
// dropdown "ID Cargo" saat create/edit produk).
func (r *wmsCargoPricedRepository) ListUnused(ctx context.Context, search string) ([]models.WMSCargoPriced, error) {
	var items []models.WMSCargoPriced
	q := r.db.WithContext(ctx).Where("is_used_in_produk = false")
	if search != "" {
		q = q.Where("code ILIKE ?", "%"+search+"%")
	}
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *wmsCargoPricedRepository) MarkUsed(ctx context.Context, cargoID string, produkID string) error {
	return r.db.WithContext(ctx).
		Model(&models.WMSCargoPriced{}).
		Where("cargo_id = ?", cargoID).
		Updates(map[string]interface{}{
			"is_used_in_produk": true,
			"produk_id":         produkID,
		}).Error
}
