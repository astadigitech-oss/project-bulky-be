package repositories

import (
	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type ForceUpdateRepository interface {
	Create(forceUpdate *models.ForceUpdateApp) error
	Update(forceUpdate *models.ForceUpdateApp) error
	Delete(id string) error
	FindByID(id string) (*models.ForceUpdateApp, error)
	FindAll(page, limit int) ([]models.ForceUpdateApp, int64, error)
	FindActive() (*models.ForceUpdateApp, error)
	SetActive(id string) error
}

type forceUpdateRepository struct {
	db *gorm.DB
}

func NewForceUpdateRepository(db *gorm.DB) ForceUpdateRepository {
	return &forceUpdateRepository{db: db}
}

func (r *forceUpdateRepository) Create(forceUpdate *models.ForceUpdateApp) error {
	return r.db.Create(forceUpdate).Error
}

func (r *forceUpdateRepository) Update(forceUpdate *models.ForceUpdateApp) error {
	return r.db.Save(forceUpdate).Error
}

func (r *forceUpdateRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.ForceUpdateApp{}).Error
}

func (r *forceUpdateRepository) FindByID(id string) (*models.ForceUpdateApp, error) {
	var forceUpdate models.ForceUpdateApp
	err := r.db.Where("id = ?", id).First(&forceUpdate).Error
	if err != nil {
		return nil, err
	}
	return &forceUpdate, nil
}

func (r *forceUpdateRepository) FindAll(page, limit int) ([]models.ForceUpdateApp, int64, error) {
	var forceUpdates []models.ForceUpdateApp
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.ForceUpdateApp{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&forceUpdates).Error
	return forceUpdates, total, err
}

func (r *forceUpdateRepository) FindActive() (*models.ForceUpdateApp, error) {
	var forceUpdate models.ForceUpdateApp
	err := r.db.Where("is_active = ?", true).First(&forceUpdate).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &forceUpdate, nil
}

func (r *forceUpdateRepository) SetActive(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deactivate all records
		if err := tx.Model(&models.ForceUpdateApp{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the selected record
		return tx.Model(&models.ForceUpdateApp{}).Where("id = ?", id).Update("is_active", true).Error
	})
}
