package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/models"
)

type PushPromotionNotificationRepository struct {
	db *gorm.DB
}

func NewPushPromotionNotificationRepository(db *gorm.DB) *PushPromotionNotificationRepository {
	return &PushPromotionNotificationRepository{db: db}
}

func (r *PushPromotionNotificationRepository) Create(ctx context.Context, item *models.PushPromotionNotification) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PushPromotionNotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PushPromotionNotification, error) {
	var item models.PushPromotionNotification
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PushPromotionNotificationRepository) FindAll(ctx context.Context, params *models.PushPromotionListRequest) ([]models.PushPromotionNotification, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.PushPromotionNotification{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.PushPromotionNotification
	err := query.Order("created_at DESC").
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Find(&items).Error
	return items, total, err
}

func (r *PushPromotionNotificationRepository) Save(ctx context.Context, item *models.PushPromotionNotification) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *PushPromotionNotificationRepository) SaveEditable(ctx context.Context, item *models.PushPromotionNotification) error {
	result := r.db.WithContext(ctx).Model(&models.PushPromotionNotification{}).
		Where("id = ? AND status IN ?", item.ID, []string{
			models.PushPromotionStatusDraft,
			models.PushPromotionStatusScheduled,
			models.PushPromotionStatusFailed,
		}).
		Updates(map[string]any{
			"nama":         item.Nama,
			"title_id":     item.TitleID,
			"body_id":      item.BodyID,
			"title_en":     item.TitleEN,
			"body_en":      item.BodyEN,
			"deep_link":    item.DeepLink,
			"scheduled_at": item.ScheduledAt,
			"status":       item.Status,
			"updated_by":   item.UpdatedBy,
			"updated_at":   time.Now().UTC(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PushPromotionNotificationRepository) DeleteDraft(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND status IN ?", id, []string{models.PushPromotionStatusDraft, models.PushPromotionStatusFailed}).
		Delete(&models.PushPromotionNotification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PushPromotionNotificationRepository) Schedule(ctx context.Context, id uuid.UUID, scheduledAt time.Time, adminID *uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.PushPromotionNotification{}).
		Where("id = ? AND status IN ?", id, []string{models.PushPromotionStatusDraft, models.PushPromotionStatusScheduled, models.PushPromotionStatusFailed}).
		Updates(map[string]any{
			"status":       models.PushPromotionStatusScheduled,
			"scheduled_at": scheduledAt,
			"updated_by":   adminID,
			"last_error":   nil,
			"updated_at":   time.Now().UTC(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PushPromotionNotificationRepository) CancelSchedule(ctx context.Context, id uuid.UUID, adminID *uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.PushPromotionNotification{}).
		Where("id = ? AND status = ?", id, models.PushPromotionStatusScheduled).
		Updates(map[string]any{
			"status":       models.PushPromotionStatusDraft,
			"scheduled_at": nil,
			"updated_by":   adminID,
			"updated_at":   time.Now().UTC(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PushPromotionNotificationRepository) ClaimForImmediateSend(ctx context.Context, id uuid.UUID) (*models.PushPromotionNotification, error) {
	var item models.PushPromotionNotification
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.PushPromotionNotification{}).
			Where("id = ? AND status IN ?", id, []string{models.PushPromotionStatusDraft, models.PushPromotionStatusScheduled, models.PushPromotionStatusFailed}).
			Updates(map[string]any{
				"status":     models.PushPromotionStatusSending,
				"last_error": nil,
				"updated_at": time.Now().UTC(),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&item, "id = ?", id).Error
	})
	return &item, err
}

func (r *PushPromotionNotificationRepository) ClaimDue(ctx context.Context, limit int) ([]models.PushPromotionNotification, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var items []models.PushPromotionNotification
	err := r.db.WithContext(ctx).Raw(
		"WITH due AS ("+
			" SELECT id FROM push_promotion_notifications"+
			" WHERE status = ? AND scheduled_at <= NOW()"+
			" ORDER BY scheduled_at ASC, id ASC"+
			" LIMIT ? FOR UPDATE SKIP LOCKED"+
			")"+
			" UPDATE push_promotion_notifications AS p"+
			" SET status = ?, last_error = NULL, updated_at = NOW()"+
			" FROM due WHERE p.id = due.id"+
			" RETURNING p.*",
		models.PushPromotionStatusScheduled, limit, models.PushPromotionStatusSending,
	).Scan(&items).Error
	return items, err
}

func (r *PushPromotionNotificationRepository) MarkSent(ctx context.Context, id uuid.UUID, successCount, failureCount int) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&models.PushPromotionNotification{}).
		Where("id = ? AND status = ?", id, models.PushPromotionStatusSending).
		Updates(map[string]any{
			"status":        models.PushPromotionStatusSent,
			"sent_at":       now,
			"success_count": successCount,
			"failure_count": failureCount,
			"last_error":    nil,
			"updated_at":    now,
		}).Error
}

func (r *PushPromotionNotificationRepository) MarkFailed(ctx context.Context, id uuid.UUID, sendErr error) error {
	if sendErr == nil {
		return errors.New("send error wajib diisi")
	}
	message := sendErr.Error()
	return r.db.WithContext(ctx).Model(&models.PushPromotionNotification{}).
		Where("id = ? AND status = ?", id, models.PushPromotionStatusSending).
		Updates(map[string]any{
			"status":     models.PushPromotionStatusFailed,
			"last_error": message,
			"updated_at": time.Now().UTC(),
		}).Error
}
