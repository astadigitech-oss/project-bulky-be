package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

var ErrWarehouseVisitValidation = errors.New("parameter kunjungan gudang tidak valid")

type WarehouseVisitService interface {
	ListPublicAvailability(ctx context.Context, fromDate, toDate string) ([]models.WarehouseVisitDateResponse, error)
	ListAdminDates(ctx context.Context, fromDate, toDate string) ([]models.WarehouseVisitDateResponse, error)
	CreateDate(ctx context.Context, req *models.WarehouseVisitDateInput, adminID uuid.UUID) (*models.WarehouseVisitDate, error)
	UpdateDate(ctx context.Context, id uuid.UUID, req *models.WarehouseVisitDateUpdateInput, adminID uuid.UUID) (*models.WarehouseVisitDate, *models.WarehouseVisitDate, error)
	DeleteDate(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitDate, error)
	ListBookings(ctx context.Context, req *models.WarehouseVisitBookingListRequest) ([]models.WarehouseVisitBookingResponse, *models.PaginationMeta, error)
	RescheduleBooking(ctx context.Context, bookingID, targetDateID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
	SetBookingAttendance(ctx context.Context, bookingID uuid.UUID, status string) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
}

type warehouseVisitService struct {
	repo          repositories.WarehouseVisitRepository
	warehouseRepo repositories.WarehouseRepository
}

func NewWarehouseVisitService(repo repositories.WarehouseVisitRepository, warehouseRepo repositories.WarehouseRepository) WarehouseVisitService {
	return &warehouseVisitService{repo: repo, warehouseRepo: warehouseRepo}
}

func (s *warehouseVisitService) ListPublicAvailability(ctx context.Context, fromDate, toDate string) ([]models.WarehouseVisitDateResponse, error) {
	from, to, err := parseVisitDateRange(fromDate, toDate)
	if err != nil {
		return nil, err
	}
	warehouse, err := s.warehouseRepo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListDates(ctx, warehouse.ID, from, to, true)
}

func (s *warehouseVisitService) ListAdminDates(ctx context.Context, fromDate, toDate string) ([]models.WarehouseVisitDateResponse, error) {
	from, to, err := parseVisitDateRange(fromDate, toDate)
	if err != nil {
		return nil, err
	}
	warehouse, err := s.warehouseRepo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListDates(ctx, warehouse.ID, from, to, false)
}

func (s *warehouseVisitService) CreateDate(ctx context.Context, req *models.WarehouseVisitDateInput, adminID uuid.UUID) (*models.WarehouseVisitDate, error) {
	date, err := parseVisitDate(req.VisitDate)
	if err != nil || req.CrewCount < 0 || (req.IsOpen && req.CrewCount == 0) {
		return nil, ErrWarehouseVisitValidation
	}
	warehouse, err := s.warehouseRepo.FindFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	visitDate := &models.WarehouseVisitDate{
		WarehouseID:    warehouse.ID,
		VisitDate:      date,
		IsOpen:         req.IsOpen,
		CrewCount:      req.CrewCount,
		VisitHoursNote: normalizeVisitHoursNote(req.VisitHoursNote),
		CreatedBy:      &adminID,
		UpdatedBy:      &adminID,
	}
	if err := s.repo.CreateDate(ctx, visitDate); err != nil {
		return nil, err
	}
	return visitDate, nil
}

func (s *warehouseVisitService) UpdateDate(ctx context.Context, id uuid.UUID, req *models.WarehouseVisitDateUpdateInput, adminID uuid.UUID) (*models.WarehouseVisitDate, *models.WarehouseVisitDate, error) {
	if req.CrewCount < 0 || (req.IsOpen && req.CrewCount == 0) {
		return nil, nil, ErrWarehouseVisitValidation
	}
	return s.repo.UpdateDate(ctx, id, req.IsOpen, req.CrewCount, normalizeVisitHoursNote(req.VisitHoursNote), adminID)
}

func (s *warehouseVisitService) DeleteDate(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitDate, error) {
	return s.repo.DeleteDate(ctx, id)
}

func (s *warehouseVisitService) ListBookings(ctx context.Context, req *models.WarehouseVisitBookingListRequest) ([]models.WarehouseVisitBookingResponse, *models.PaginationMeta, error) {
	filter := models.WarehouseVisitBookingFilter{Status: strings.ToUpper(strings.TrimSpace(req.Status))}
	if filter.Status != "" && filter.Status != models.WarehouseVisitBookingConfirmed &&
		filter.Status != models.WarehouseVisitBookingCancelled && filter.Status != models.WarehouseVisitBookingAttended &&
		filter.Status != models.WarehouseVisitBookingNoShow {
		return nil, nil, ErrWarehouseVisitValidation
	}
	if req.FromDate != "" {
		value, err := parseVisitDate(req.FromDate)
		if err != nil {
			return nil, nil, err
		}
		filter.FromDate = &value
	}
	if req.ToDate != "" {
		value, err := parseVisitDate(req.ToDate)
		if err != nil {
			return nil, nil, err
		}
		filter.ToDate = &value
	}
	if filter.FromDate != nil && filter.ToDate != nil && filter.FromDate.After(*filter.ToDate) {
		return nil, nil, ErrWarehouseVisitValidation
	}

	params := req.PaginationRequest
	params.SetDefaults()
	filter.Page, filter.PerPage = params.Page, params.PerPage
	bookings, total, err := s.repo.ListBookings(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return bookings, &meta, nil
}

func (s *warehouseVisitService) RescheduleBooking(ctx context.Context, bookingID, targetDateID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	return s.repo.RescheduleBooking(ctx, bookingID, targetDateID)
}

func (s *warehouseVisitService) CancelBooking(ctx context.Context, bookingID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	return s.repo.CancelBooking(ctx, bookingID)
}

func (s *warehouseVisitService) SetBookingAttendance(ctx context.Context, bookingID uuid.UUID, status string) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != models.WarehouseVisitBookingAttended && status != models.WarehouseVisitBookingNoShow {
		return nil, nil, ErrWarehouseVisitValidation
	}
	return s.repo.SetBookingAttendance(ctx, bookingID, status)
}

func parseVisitDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrWarehouseVisitValidation
	}
	return parsed, nil
}

func parseVisitDateRange(fromValue, toValue string) (time.Time, time.Time, error) {
	from, err := parseVisitDate(fromValue)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := parseVisitDate(toValue)
	if err != nil || from.After(to) {
		return time.Time{}, time.Time{}, ErrWarehouseVisitValidation
	}
	return from, to, nil
}

func normalizeVisitHoursNote(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
