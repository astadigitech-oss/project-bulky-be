package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type PushPromotionNotificationService struct {
	repo       *repositories.PushPromotionNotificationRepository
	pushClient *PushTestService
	now        func() time.Time
}

func NewPushPromotionNotificationService(
	repo *repositories.PushPromotionNotificationRepository,
	pushClient *PushTestService,
) *PushPromotionNotificationService {
	return &PushPromotionNotificationService{
		repo:       repo,
		pushClient: pushClient,
		now:        func() time.Time { return time.Now().UTC() },
	}
}

func (s *PushPromotionNotificationService) Create(ctx context.Context, req *models.CreatePushPromotionRequest, adminID *uuid.UUID) (*models.PushPromotionNotification, error) {
	item := &models.PushPromotionNotification{
		ID:        uuid.New(),
		Nama:      strings.TrimSpace(req.Nama),
		TitleID:   strings.TrimSpace(req.TitleID),
		BodyID:    strings.TrimSpace(req.BodyID),
		TitleEN:   strings.TrimSpace(req.TitleEN),
		BodyEN:    strings.TrimSpace(req.BodyEN),
		DeepLink:  "/products",
		Status:    models.PushPromotionStatusDraft,
		CreatedBy: adminID,
		UpdatedBy: adminID,
	}
	if err := validatePushPromotion(item); err != nil {
		return nil, err
	}
	if req.ScheduledAt != nil && strings.TrimSpace(*req.ScheduledAt) != "" {
		scheduledAt, err := parsePushSchedule(*req.ScheduledAt, s.now())
		if err != nil {
			return nil, err
		}
		item.ScheduledAt = &scheduledAt
		item.Status = models.PushPromotionStatusScheduled
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PushPromotionNotificationService) FindAll(ctx context.Context, params *models.PushPromotionListRequest) (*models.PushPromotionListResponse, error) {
	params.PaginationRequest.SetDefaults()
	if params.Status != "" && !validPushPromotionStatusFilter(params.Status) {
		params.Status = ""
	}
	items, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}
	return &models.PushPromotionListResponse{
		Data: items,
		Meta: models.NewPaginationMeta(params.Page, params.PerPage, total),
	}, nil
}

func (s *PushPromotionNotificationService) FindByID(ctx context.Context, rawID string) (*models.PushPromotionNotification, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.New("ID notifikasi promo tidak valid")
	}
	item, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("notifikasi promo tidak ditemukan")
	}
	return item, err
}

func (s *PushPromotionNotificationService) Update(ctx context.Context, rawID string, req *models.UpdatePushPromotionRequest, adminID *uuid.UUID) (*models.PushPromotionNotification, error) {
	item, err := s.FindByID(ctx, rawID)
	if err != nil {
		return nil, err
	}
	if !editablePushPromotionStatus(item.Status) {
		return nil, errors.New("notifikasi yang sedang dikirim atau sudah terkirim tidak dapat diubah")
	}

	if req.Nama != nil {
		item.Nama = strings.TrimSpace(*req.Nama)
	}
	if req.TitleID != nil {
		item.TitleID = strings.TrimSpace(*req.TitleID)
	}
	if req.BodyID != nil {
		item.BodyID = strings.TrimSpace(*req.BodyID)
	}
	if req.TitleEN != nil {
		item.TitleEN = strings.TrimSpace(*req.TitleEN)
	}
	if req.BodyEN != nil {
		item.BodyEN = strings.TrimSpace(*req.BodyEN)
	}
	// Promo notifications always open the Store product listing.
	item.DeepLink = "/products"
	if req.ScheduledAt != nil {
		if strings.TrimSpace(*req.ScheduledAt) == "" {
			item.ScheduledAt = nil
			item.Status = models.PushPromotionStatusDraft
		} else {
			scheduledAt, parseErr := parsePushSchedule(*req.ScheduledAt, s.now())
			if parseErr != nil {
				return nil, parseErr
			}
			item.ScheduledAt = &scheduledAt
			item.Status = models.PushPromotionStatusScheduled
		}
	}
	item.UpdatedBy = adminID
	if err := validatePushPromotion(item); err != nil {
		return nil, err
	}
	if err := s.repo.SaveEditable(ctx, item); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notifikasi tidak dapat diubah karena statusnya sudah berubah")
		}
		return nil, err
	}
	return s.FindByID(ctx, rawID)
}

func (s *PushPromotionNotificationService) Schedule(ctx context.Context, rawID, rawTime string, adminID *uuid.UUID) (*models.PushPromotionNotification, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.New("ID notifikasi promo tidak valid")
	}
	scheduledAt, err := parsePushSchedule(rawTime, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Schedule(ctx, id, scheduledAt, adminID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notifikasi tidak dapat dijadwalkan pada status saat ini")
		}
		return nil, err
	}
	return s.FindByID(ctx, rawID)
}

func (s *PushPromotionNotificationService) CancelSchedule(ctx context.Context, rawID string, adminID *uuid.UUID) (*models.PushPromotionNotification, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.New("ID notifikasi promo tidak valid")
	}
	if err := s.repo.CancelSchedule(ctx, id, adminID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notifikasi tidak sedang dijadwalkan")
		}
		return nil, err
	}
	return s.FindByID(ctx, rawID)
}

func (s *PushPromotionNotificationService) Delete(ctx context.Context, rawID string) error {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return errors.New("ID notifikasi promo tidak valid")
	}
	if err := s.repo.DeleteDraft(ctx, id); errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("hanya draft atau notifikasi gagal yang dapat dihapus")
	} else {
		return err
	}
}

func (s *PushPromotionNotificationService) SendNow(ctx context.Context, rawID string) (*models.PushPromotionNotification, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.New("ID notifikasi promo tidak valid")
	}
	item, err := s.repo.ClaimForImmediateSend(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("notifikasi tidak dapat dikirim pada status saat ini")
	}
	if err != nil {
		return nil, err
	}
	return s.deliver(ctx, item)
}

func (s *PushPromotionNotificationService) StartScheduler(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[push-promotion-scheduler] scheduler dihentikan")
			return
		case <-ticker.C:
			s.processDue(ctx)
		}
	}
}

func (s *PushPromotionNotificationService) processDue(ctx context.Context) {
	items, err := s.repo.ClaimDue(ctx, 20)
	if err != nil {
		log.Printf("[push-promotion-scheduler] gagal mengambil jadwal: %v", err)
		return
	}
	for i := range items {
		sendCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		if _, err := s.deliver(sendCtx, &items[i]); err != nil {
			log.Printf("[push-promotion-scheduler] notifikasi %s gagal: %v", items[i].ID, err)
		}
		cancel()
	}
}

func (s *PushPromotionNotificationService) deliver(ctx context.Context, item *models.PushPromotionNotification) (*models.PushPromotionNotification, error) {
	result, err := s.pushClient.SendPromotion(ctx, item.ID.String(), item.TitleID, item.BodyID, item.TitleEN, item.BodyEN, "/products")
	if err != nil {
		markCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if markErr := s.repo.MarkFailed(markCtx, item.ID, err); markErr != nil {
			log.Printf("[push-promotion] gagal menyimpan status gagal untuk %s: %v", item.ID, markErr)
		}
		return nil, fmt.Errorf("gagal mengirim notifikasi promo: %w", err)
	}
	if err := s.repo.MarkSent(ctx, item.ID, result.SuccessCount, result.FailureCount); err != nil {
		return nil, err
	}
	return s.FindByID(ctx, item.ID.String())
}

func validatePushPromotion(item *models.PushPromotionNotification) error {
	switch {
	case item.Nama == "":
		return errors.New("nama promo wajib diisi")
	case item.TitleID == "" || item.BodyID == "":
		return errors.New("judul dan isi bahasa Indonesia wajib diisi")
	case item.TitleEN == "" || item.BodyEN == "":
		return errors.New("judul dan isi bahasa Inggris wajib diisi")
	case pushTextLength(item.TitleID) > 50 || pushTextLength(item.TitleEN) > 50:
		return errors.New("judul notifikasi maksimal 50 karakter")
	case pushTextLength(item.BodyID) > 120 || pushTextLength(item.BodyEN) > 120:
		return errors.New("isi notifikasi maksimal 120 karakter")
	case len(item.Nama) > 120:
		return errors.New("nama promo maksimal 120 karakter")
	case item.DeepLink != "/products":
		return errors.New("tujuan notifikasi promo harus halaman produk Store")
	}
	return nil
}

func pushTextLength(value string) int {
	length := 0
	for _, char := range value {
		length += utf16.RuneLen(char)
	}
	return length
}

func parsePushSchedule(raw string, now time.Time) (time.Time, error) {
	value, err := time.Parse(time.RFC3339, strings.TrimSpace(raw))
	if err != nil {
		return time.Time{}, errors.New("jadwal harus berupa tanggal dan waktu ISO 8601")
	}
	value = value.UTC()
	if !value.After(now) {
		return time.Time{}, errors.New("waktu jadwal harus berada di masa depan")
	}
	return value, nil
}

func editablePushPromotionStatus(status string) bool {
	return status == models.PushPromotionStatusDraft ||
		status == models.PushPromotionStatusScheduled ||
		status == models.PushPromotionStatusFailed
}

func validPushPromotionStatusFilter(status string) bool {
	switch status {
	case models.PushPromotionStatusDraft, models.PushPromotionStatusScheduled, models.PushPromotionStatusSending, models.PushPromotionStatusSent, models.PushPromotionStatusFailed:
		return true
	default:
		return false
	}
}
