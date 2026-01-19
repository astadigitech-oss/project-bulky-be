package repositories

import (
	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ModeMaintenanceRepository interface {
	Create(maintenance *models.ModeMaintenance) error
	Update(maintenance *models.ModeMaintenance) error
	Delete(id string) error
	FindByID(id string) (*models.ModeMaintenance, error)
	FindAll(page, limit int) ([]models.ModeMaintenance, int64, error)
	FindActive() (*models.ModeMaintenance, error)
	Activate(id string) error
	Deactivate(id string) error
}

type modeMaintenanceRepository struct {
	db *gorm.DB
}

func NewModeMaintenanceRepository(db *gorm.DB) ModeMaintenanceRepository {
	return &modeMaintenanceRepository{db: db}
}

func (r *modeMaintenanceRepository) Create(maintenance *models.ModeMaintenance) error {
	return r.db.Create(maintenance).Error
}

func (r *modeMaintenanceRepository) Update(maintenance *models.ModeMaintenance) error {
	return r.db.Save(maintenance).Error
}

func (r *modeMaintenanceRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.ModeMaintenance{}).Error
}

func (r *modeMaintenanceRepository) FindByID(id string) (*models.ModeMaintenance, error) {
	var maintenance models.ModeMaintenance
	err := r.db.Where("id = ?", id).First(&maintenance).Error
	if err != nil {
		return nil, err
	}
	return &maintenance, nil
}

func (r *modeMaintenanceRepository) FindAll(page, limit int) ([]models.ModeMaintenance, int64, error) {
	var maintenances []models.ModeMaintenance
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.ModeMaintenance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&maintenances).Error
	return maintenances, total, err
}

func (r *modeMaintenanceRepository) FindActive() (*models.ModeMaintenance, error) {
	var maintenance models.ModeMaintenance
	err := r.db.Where("is_active = ?", true).First(&maintenance).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &maintenance, nil
}

func (r *modeMaintenanceRepository) Activate(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deactivate all records
		if err := tx.Model(&models.ModeMaintenance{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the selected record
		return tx.Model(&models.ModeMaintenance{}).Where("id = ?", id).Update("is_active", true).Error
	})
}

func (r *modeMaintenanceRepository) Deactivate(id string) error {
	return r.db.Model(&models.ModeMaintenance{}).Where("id = ?", id).Update("is_active", false).Error
}
