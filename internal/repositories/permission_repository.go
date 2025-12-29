package repositories

import (
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindAll() ([]models.Permission, error)
	FindByID(id uuid.UUID) (*models.Permission, error)
	FindByModul() (map[string][]models.Permission, error)
	FindByIDs(ids []uuid.UUID) ([]models.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Order("modul, nama").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByID(id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	return &permission, err
}

func (r *permissionRepository) FindByModul() (map[string][]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Order("modul, nama").Find(&permissions).Error; err != nil {
		return nil, err
	}

	// Group by modul
	result := make(map[string][]models.Permission)
	for _, perm := range permissions {
		result[perm.Modul] = append(result[perm.Modul], perm)
	}

	return result, nil
}

func (r *permissionRepository) FindByIDs(ids []uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}
