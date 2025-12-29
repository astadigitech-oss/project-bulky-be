package repositories

import (
	"project-bulky-be/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	// Admin
	FindAdminByEmail(email string) (*models.Admin, error)
	FindAdminByID(id uuid.UUID) (*models.Admin, error)
	FindAdminWithRole(id uuid.UUID) (*models.Admin, error)
	UpdateAdminLastLogin(id uuid.UUID) error

	// Buyer
	FindBuyerByEmail(email string) (*models.Buyer, error)
	FindBuyerByID(id uuid.UUID) (*models.Buyer, error)
	UpdateBuyerLastLogin(id uuid.UUID) error

	// Refresh Token
	CreateRefreshToken(token *models.RefreshToken) error
	FindRefreshToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userType models.UserType, userID uuid.UUID) error
	CleanupExpiredTokens() error

	// Role & Permission
	GetRoleWithPermissions(roleID uuid.UUID) (*models.Role, error)
	GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Admin methods
func (r *authRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&admin).Error
	return &admin, err
}

func (r *authRepository) FindAdminByID(id uuid.UUID) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&admin).Error
	return &admin, err
}

func (r *authRepository) FindAdminWithRole(id uuid.UUID) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Preload("Role.Permissions").
		Where("admin.id = ? AND admin.deleted_at IS NULL", id).
		First(&admin).Error
	return &admin, err
}

func (r *authRepository) UpdateAdminLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// Buyer methods
func (r *authRepository) FindBuyerByEmail(email string) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&buyer).Error
	return &buyer, err
}

func (r *authRepository) FindBuyerByID(id uuid.UUID) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&buyer).Error
	return &buyer, err
}

func (r *authRepository) UpdateBuyerLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Buyer{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// Refresh Token methods
func (r *authRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ? AND is_revoked = false AND expired_at > NOW()", token).
		First(&refreshToken).Error
	return &refreshToken, err
}

func (r *authRepository) RevokeRefreshToken(token string) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *authRepository) RevokeAllUserTokens(userType models.UserType, userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_type = ? AND user_id = ? AND is_revoked = false", userType, userID).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *authRepository) CleanupExpiredTokens() error {
	return r.db.Where("expired_at < NOW() OR is_revoked = true").
		Delete(&models.RefreshToken{}).Error
}

// Role & Permission methods
func (r *authRepository) GetRoleWithPermissions(roleID uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").
		Where("id = ? AND deleted_at IS NULL", roleID).
		First(&role).Error
	return &role, err
}

func (r *authRepository) GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permission").
		Joins("INNER JOIN role_permission ON permission.id = role_permission.permission_id").
		Where("role_permission.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}
