package repositories

import (
	"context"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JadwalGudangRepository interface {
	FindByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]models.JadwalGudang, error)
	UpdateBatch(ctx context.Context, warehouseID uuid.UUID, jadwal []models.JadwalGudang) error
}

type jadwalGudangRepository struct {
	db *gorm.DB
}

func NewJadwalGudangRepository(db *gorm.DB) JadwalGudangRepository {
	return &jadwalGudangRepository{db: db}
}

func (r *jadwalGudangRepository) FindByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]models.JadwalGudang, error) {
	var jadwal []models.JadwalGudang
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ?", warehouseID).
		Order("hari ASC").
		Find(&jadwal).Error
	return jadwal, err
}

func (r *jadwalGudangRepository) UpdateBatch(ctx context.Context, warehouseID uuid.UUID, jadwal []models.JadwalGudang) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing jadwal
		if err := tx.Where("warehouse_id = ?", warehouseID).Delete(&models.JadwalGudang{}).Error; err != nil {
			return err
		}

		// Insert new jadwal
		for _, j := range jadwal {
			j.WarehouseID = warehouseID
			if err := tx.Create(&j).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
