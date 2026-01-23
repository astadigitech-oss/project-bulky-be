package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type InformasiPickupRepository interface {
	Get(ctx context.Context) (*models.InformasiPickup, error)
	Update(ctx context.Context, pickup *models.InformasiPickup) error
	UpdateJadwal(ctx context.Context, pickupID string, jadwal []models.JadwalGudang) error
}

type informasiPickupRepository struct {
	db *gorm.DB
}

func NewInformasiPickupRepository(db *gorm.DB) InformasiPickupRepository {
	return &informasiPickupRepository{db: db}
}

func (r *informasiPickupRepository) Get(ctx context.Context) (*models.InformasiPickup, error) {
	var pickup models.InformasiPickup
	err := r.db.WithContext(ctx).
		Preload("JadwalGudang", func(db *gorm.DB) *gorm.DB {
			return db.Order("hari ASC")
		}).
		First(&pickup).Error
	if err != nil {
		return nil, err
	}
	return &pickup, nil
}

func (r *informasiPickupRepository) Update(ctx context.Context, pickup *models.InformasiPickup) error {
	return r.db.WithContext(ctx).Save(pickup).Error
}

func (r *informasiPickupRepository) UpdateJadwal(ctx context.Context, pickupID string, jadwal []models.JadwalGudang) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update each jadwal by hari
		for _, j := range jadwal {
			if err := tx.Model(&models.JadwalGudang{}).
				Where("informasi_pickup_id = ? AND hari = ?", pickupID, j.Hari).
				Updates(map[string]interface{}{
					"jam_buka":  j.JamBuka,
					"jam_tutup": j.JamTutup,
					"is_buka":   j.IsBuka,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
