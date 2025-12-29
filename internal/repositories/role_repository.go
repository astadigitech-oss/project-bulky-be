package repositories

import (
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]models.Role, error)
	FindByID(id uuid.UUID) (*models.Role, error)
	FindByIDWithPermissions(id uuid.UUID) (*models.Role, error)
	FindByKode(kode string) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uuid.UUID) error

	// Permission operations
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemoveAllPermissions(roleID uuid.UUID) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	return &role, err
}

func (r *roleRepository) FindByIDWithPermissions(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&role).Error
	return &role, err
}

func (r *roleRepository) FindByKode(kode string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("kode = ? AND deleted_at IS NULL", kode).First(&role).Error
	return &role, err
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Role{}).Error
}

func (r *roleRepository) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	// Remove existing permissions
	if err := r.RemoveAllPermissions(roleID); err != nil {
		return err
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := r.db.Create(&rolePermission).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *roleRepository) RemoveAllPermissions(roleID uuid.UUID) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}
