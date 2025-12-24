package repositories

import (
	"context"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type AdminSessionRepository interface {
	Create(ctx context.Context, session *models.AdminSession) error
	FindByToken(ctx context.Context, token string) (*models.AdminSession, error)
	FindByAdminID(ctx context.Context, adminID string) ([]models.AdminSession, error)
	Delete(ctx context.Context, id string) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteByAdminID(ctx context.Context, adminID string) error
	DeleteExpired(ctx context.Context) error
}

type adminSessionRepository struct {
	db *gorm.DB
}

func NewAdminSessionRepository(db *gorm.DB) AdminSessionRepository {
	return &adminSessionRepository{db: db}
}

func (r *adminSessionRepository) Create(ctx context.Context, session *models.AdminSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *adminSessionRepository) FindByToken(ctx context.Context, token string) (*models.AdminSession, error) {
	var session models.AdminSession
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *adminSessionRepository) FindByAdminID(ctx context.Context, adminID string) ([]models.AdminSession, error) {
	var sessions []models.AdminSession
	err := r.db.WithContext(ctx).Where("admin_id = ?", adminID).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *adminSessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.AdminSession{}).Error
}

func (r *adminSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.AdminSession{}).Error
}

func (r *adminSessionRepository) DeleteByAdminID(ctx context.Context, adminID string) error {
	return r.db.WithContext(ctx).Where("admin_id = ?", adminID).Delete(&models.AdminSession{}).Error
}

func (r *adminSessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < NOW()").Delete(&models.AdminSession{}).Error
}
