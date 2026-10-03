package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/models"
)

type AdminNotificationRepository struct{ db *gorm.DB }

func NewAdminNotificationRepository(db *gorm.DB) *AdminNotificationRepository {
	return &AdminNotificationRepository{db: db}
}

func (r *AdminNotificationRepository) UpsertDevice(ctx context.Context, adminID uuid.UUID, fid string) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO admin_push_devices (admin_id, fid, platform, last_seen_at, created_at, updated_at)
		VALUES (?, ?, 'web', NOW(), NOW(), NOW())
		ON CONFLICT (fid) DO UPDATE
		SET admin_id = EXCLUDED.admin_id,
		    platform = EXCLUDED.platform,
		    last_seen_at = NOW(),
		    updated_at = NOW()
	`, adminID, fid).Error
}

func (r *AdminNotificationRepository) DeleteDevice(ctx context.Context, adminID uuid.UUID, fid string) error {
	return r.db.WithContext(ctx).Where("admin_id = ? AND fid = ?", adminID, fid).Delete(&models.AdminPushDevice{}).Error
}

func (r *AdminNotificationRepository) ListForAdmin(ctx context.Context, adminID uuid.UUID, limit int) (*models.AdminNotificationListResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 30
	}
	response := &models.AdminNotificationListResponse{Notifications: []models.AdminNotification{}}
	if err := r.db.WithContext(ctx).Where("admin_id = ?", adminID).
		Order("created_at DESC").Limit(limit).Find(&response.Notifications).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.AdminNotification{}).
		Where("admin_id = ? AND is_read = FALSE", adminID).Count(&response.UnreadCount).Error; err != nil {
		return nil, err
	}
	return response, nil
}

func (r *AdminNotificationRepository) MarkRead(ctx context.Context, adminID, notificationID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.AdminNotification{}).
		Where("admin_id = ? AND id = ?", adminID, notificationID).
		Updates(map[string]any{"is_read": true, "read_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		var exists int64
		if err := r.db.WithContext(ctx).Model(&models.AdminNotification{}).
			Where("admin_id = ? AND id = ?", adminID, notificationID).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return gorm.ErrRecordNotFound
		}
	}
	return nil
}
