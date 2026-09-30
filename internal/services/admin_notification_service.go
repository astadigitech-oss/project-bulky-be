package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

var ErrInvalidAdminPushFID = errors.New("fid wajib diisi dan maksimal 255 karakter")

type AdminNotificationService struct {
	repo *repositories.AdminNotificationRepository
}

func NewAdminNotificationService(repo *repositories.AdminNotificationRepository) *AdminNotificationService {
	return &AdminNotificationService{repo: repo}
}

func (s *AdminNotificationService) RegisterDevice(ctx context.Context, adminID uuid.UUID, fid string) error {
	fid = strings.TrimSpace(fid)
	if fid == "" || len(fid) > 255 {
		return ErrInvalidAdminPushFID
	}
	return s.repo.UpsertDevice(ctx, adminID, fid)
}

func (s *AdminNotificationService) UnregisterDevice(ctx context.Context, adminID uuid.UUID, fid string) error {
	fid = strings.TrimSpace(fid)
	if fid == "" || len(fid) > 255 {
		return ErrInvalidAdminPushFID
	}
	return s.repo.DeleteDevice(ctx, adminID, fid)
}

func (s *AdminNotificationService) List(ctx context.Context, adminID uuid.UUID) (*models.AdminNotificationListResponse, error) {
	return s.repo.ListForAdmin(ctx, adminID, 30)
}

func (s *AdminNotificationService) MarkRead(ctx context.Context, adminID, notificationID uuid.UUID) error {
	return s.repo.MarkRead(ctx, adminID, notificationID)
}
