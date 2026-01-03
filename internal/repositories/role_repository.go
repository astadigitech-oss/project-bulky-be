package repositories

import (
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]models.Role, error)
	FindAllWithParams(params dto.RoleQueryParams) ([]models.Role, int64, error)
	FindByID(id uuid.UUID) (*models.Role, error)
	FindByIDWithPermissions(id uuid.UUID) (*models.Role, error)
	FindByKode(kode string) (*models.Role, error)
	FindByKodeExcludeID(kode string, excludeID uuid.UUID) (*models.Role, error)
	CountAdminByRoleID(roleID uuid.UUID) (int64, error)
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

func (r *roleRepository) FindAllWithParams(params dto.RoleQueryParams) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.Model(&models.Role{}).Where("deleted_at IS NULL")

	// Filter by search
	if params.Cari != "" {
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", "%"+params.Cari+"%", "%"+params.Cari+"%")
	}

	// Filter by is_active
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Set default values
	if params.Halaman == 0 {
		params.Halaman = 1
	}
	if params.PerHalaman == 0 {
		params.PerHalaman = 10
	}
	if params.UrutBerdasarkan == "" {
		params.UrutBerdasarkan = "created_at"
	}
	if params.Urutan == "" {
		params.Urutan = "desc"
	}

	// Sort
	orderClause := params.UrutBerdasarkan + " " + params.Urutan
	query = query.Order(orderClause)

	// Pagination
	offset := (params.Halaman - 1) * params.PerHalaman
	query = query.Offset(offset).Limit(params.PerHalaman)

	err := query.Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) FindByKodeExcludeID(kode string, excludeID uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("kode = ? AND id != ? AND deleted_at IS NULL", kode, excludeID).First(&role).Error
	return &role, err
}

func (r *roleRepository) CountAdminByRoleID(roleID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Admin{}).Where("role_id = ? AND deleted_at IS NULL", roleID).Count(&count).Error
	return count, err
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
