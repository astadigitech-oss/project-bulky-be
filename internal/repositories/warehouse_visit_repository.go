package repositories

import (
	"context"
	"errors"
	"sort"
	"time"

	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWarehouseVisitNotFound    = errors.New("data kunjungan gudang tidak ditemukan")
	ErrWarehouseVisitConflict    = errors.New("perubahan tidak dapat dilakukan karena kapasitas atau status booking")
	ErrWarehouseVisitHasBookings = errors.New("jadwal tidak dapat dihapus karena memiliki riwayat booking")
)

type WarehouseVisitRepository interface {
	ListDates(ctx context.Context, warehouseID uuid.UUID, from, to time.Time, openOnly bool) ([]models.WarehouseVisitDateResponse, error)
	CreateDate(ctx context.Context, visitDate *models.WarehouseVisitDate) error
	UpdateDate(ctx context.Context, id uuid.UUID, isOpen bool, crewCount int, visitHoursNote *string, adminID uuid.UUID) (*models.WarehouseVisitDate, *models.WarehouseVisitDate, error)
	DeleteDate(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitDate, error)
	ListBookings(ctx context.Context, filter models.WarehouseVisitBookingFilter) ([]models.WarehouseVisitBookingResponse, int64, error)
	RescheduleBooking(ctx context.Context, id, targetDateID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
	CancelBooking(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
	SetBookingAttendance(ctx context.Context, id uuid.UUID, status string) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error)
}

type warehouseVisitRepository struct {
	db *gorm.DB
}

func NewWarehouseVisitRepository(db *gorm.DB) WarehouseVisitRepository {
	return &warehouseVisitRepository{db: db}
}

func (r *warehouseVisitRepository) ListDates(ctx context.Context, warehouseID uuid.UUID, from, to time.Time, openOnly bool) ([]models.WarehouseVisitDateResponse, error) {
	var dates []models.WarehouseVisitDateResponse
	query := r.db.WithContext(ctx).Table("warehouse_visit_dates AS d").
		Select(`
			d.id::text AS id,
			d.warehouse_id::text AS warehouse_id,
			to_char(d.visit_date, 'YYYY-MM-DD') AS visit_date,
			d.is_open,
			d.crew_count,
			d.crew_count * 2 AS capacity,
			COALESCE(b.booked_count, 0)::bigint AS booked_count,
			GREATEST(d.crew_count * 2 - COALESCE(b.booked_count, 0), 0)::bigint AS remaining_capacity,
			d.visit_hours_note,
			COALESCE(
				NULLIF(BTRIM(d.visit_hours_note), ''),
				CASE WHEN j.is_buka THEN
					to_char(j.jam_buka, 'HH24:MI') || ' - ' || to_char(j.jam_tutup, 'HH24:MI')
				END
			) AS display_hours`).
		Joins(`LEFT JOIN (
			SELECT visit_date_id, COUNT(*) AS booked_count
			FROM warehouse_visit_bookings
			WHERE status <> ?
			GROUP BY visit_date_id
		) AS b ON b.visit_date_id = d.id`, models.WarehouseVisitBookingCancelled).
		Joins(`LEFT JOIN jadwal_gudang AS j
			ON j.warehouse_id = d.warehouse_id
			AND j.hari = EXTRACT(DOW FROM d.visit_date)::int`).
		Where("d.warehouse_id = ? AND d.visit_date >= ? AND d.visit_date <= ?", warehouseID, from, to)
	if openOnly {
		query = query.Where("d.is_open = TRUE")
	}
	if err := query.Order("d.visit_date ASC").Scan(&dates).Error; err != nil {
		return nil, err
	}
	if dates == nil {
		dates = []models.WarehouseVisitDateResponse{}
	}
	return dates, nil
}

func (r *warehouseVisitRepository) CreateDate(ctx context.Context, visitDate *models.WarehouseVisitDate) error {
	if err := r.db.WithContext(ctx).Create(visitDate).Error; err != nil {
		return translateWarehouseVisitDBError(err)
	}
	return nil
}

func (r *warehouseVisitRepository) UpdateDate(ctx context.Context, id uuid.UUID, isOpen bool, crewCount int, visitHoursNote *string, adminID uuid.UUID) (*models.WarehouseVisitDate, *models.WarehouseVisitDate, error) {
	var previous, updated models.WarehouseVisitDate
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&previous, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWarehouseVisitNotFound
			}
			return err
		}

		var bookedCount int64
		if err := tx.Model(&models.WarehouseVisitBooking{}).
			Where("visit_date_id = ? AND status <> ?", id, models.WarehouseVisitBookingCancelled).
			Count(&bookedCount).Error; err != nil {
			return err
		}
		newCapacity := int64(crewCount * 2)
		if newCapacity < bookedCount || (!isOpen && bookedCount > 0) {
			return ErrWarehouseVisitConflict
		}

		if err := tx.Model(&models.WarehouseVisitDate{}).Where("id = ?", id).Updates(map[string]interface{}{
			"is_open":          isOpen,
			"crew_count":       crewCount,
			"visit_hours_note": visitHoursNote,
			"updated_by":       adminID,
		}).Error; err != nil {
			return err
		}
		if err := tx.First(&updated, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, translateWarehouseVisitDBError(err)
	}
	return &previous, &updated, nil
}

func (r *warehouseVisitRepository) DeleteDate(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitDate, error) {
	var previous models.WarehouseVisitDate
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&previous, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWarehouseVisitNotFound
			}
			return err
		}

		var bookingCount int64
		if err := tx.Model(&models.WarehouseVisitBooking{}).
			Where("visit_date_id = ?", id).
			Count(&bookingCount).Error; err != nil {
			return err
		}
		if bookingCount > 0 {
			return ErrWarehouseVisitHasBookings
		}

		return tx.Delete(&previous).Error
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrWarehouseVisitHasBookings
		}
		return nil, translateWarehouseVisitDBError(err)
	}
	return &previous, nil
}

func (r *warehouseVisitRepository) ListBookings(ctx context.Context, filter models.WarehouseVisitBookingFilter) ([]models.WarehouseVisitBookingResponse, int64, error) {
	query := r.db.WithContext(ctx).Table("warehouse_visit_bookings AS b").
		Joins("JOIN buyer AS buyer ON buyer.id = b.buyer_id")
	if filter.FromDate != nil {
		query = query.Where("b.visit_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("b.visit_date <= ?", *filter.ToDate)
	}
	if filter.Status != "" {
		query = query.Where("b.status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var bookings []models.WarehouseVisitBookingResponse
	err := query.Select(`
		b.id::text AS id,
		b.visit_date_id::text AS visit_date_id,
		to_char(b.visit_date, 'YYYY-MM-DD') AS visit_date,
		b.buyer_id::text AS buyer_id,
		buyer.nama AS buyer_name,
		buyer.email AS buyer_email,
		buyer.telepon AS buyer_phone,
		b.status,
		b.created_at,
		b.updated_at`).
		Order("b.visit_date DESC, b.created_at DESC, b.id DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Scan(&bookings).Error
	if err != nil {
		return nil, 0, err
	}
	if bookings == nil {
		bookings = []models.WarehouseVisitBookingResponse{}
	}
	return bookings, total, nil
}

func (r *warehouseVisitRepository) RescheduleBooking(ctx context.Context, id, targetDateID uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	var previous, updated models.WarehouseVisitBooking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&previous, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWarehouseVisitNotFound
			}
			return err
		}
		if previous.Status != models.WarehouseVisitBookingConfirmed {
			return ErrWarehouseVisitConflict
		}
		if previous.VisitDateID == targetDateID {
			return ErrWarehouseVisitConflict
		}

		lockedDates, err := lockVisitDates(tx, previous.VisitDateID, targetDateID)
		if err != nil {
			return err
		}
		var target *models.WarehouseVisitDate
		for i := range lockedDates {
			if lockedDates[i].ID == targetDateID {
				target = &lockedDates[i]
				break
			}
		}
		if target == nil {
			return ErrWarehouseVisitNotFound
		}
		if !target.IsOpen {
			return ErrWarehouseVisitConflict
		}

		var targetBooked int64
		if err := tx.Model(&models.WarehouseVisitBooking{}).
			Where("visit_date_id = ? AND status <> ?", targetDateID, models.WarehouseVisitBookingCancelled).
			Count(&targetBooked).Error; err != nil {
			return err
		}
		if targetBooked >= int64(target.CrewCount*2) {
			return ErrWarehouseVisitConflict
		}

		if err := tx.Model(&models.WarehouseVisitBooking{}).Where("id = ?", id).Updates(map[string]interface{}{
			"visit_date_id": target.ID,
			"visit_date":    target.VisitDate,
		}).Error; err != nil {
			return err
		}
		if err := tx.First(&updated, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, translateWarehouseVisitDBError(err)
	}
	return &previous, &updated, nil
}

func (r *warehouseVisitRepository) CancelBooking(ctx context.Context, id uuid.UUID) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	var previous, updated models.WarehouseVisitBooking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&previous, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWarehouseVisitNotFound
			}
			return err
		}
		if previous.Status != models.WarehouseVisitBookingConfirmed {
			return ErrWarehouseVisitConflict
		}
		if err := tx.Model(&models.WarehouseVisitBooking{}).Where("id = ?", id).
			Update("status", models.WarehouseVisitBookingCancelled).Error; err != nil {
			return err
		}
		if err := tx.First(&updated, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, translateWarehouseVisitDBError(err)
	}
	return &previous, &updated, nil
}

func (r *warehouseVisitRepository) SetBookingAttendance(ctx context.Context, id uuid.UUID, status string) (*models.WarehouseVisitBooking, *models.WarehouseVisitBooking, error) {
	var previous, updated models.WarehouseVisitBooking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&previous, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWarehouseVisitNotFound
			}
			return err
		}
		if previous.Status != models.WarehouseVisitBookingConfirmed {
			return ErrWarehouseVisitConflict
		}
		if err := tx.Model(&models.WarehouseVisitBooking{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}
		if err := tx.First(&updated, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, translateWarehouseVisitDBError(err)
	}
	return &previous, &updated, nil
}

func lockVisitDates(tx *gorm.DB, ids ...uuid.UUID) ([]models.WarehouseVisitDate, error) {
	unique := make(map[uuid.UUID]struct{}, len(ids))
	ordered := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, exists := unique[id]; !exists {
			unique[id] = struct{}{}
			ordered = append(ordered, id)
		}
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].String() < ordered[j].String() })

	locked := make([]models.WarehouseVisitDate, 0, len(ordered))
	for _, id := range ordered {
		var visitDate models.WarehouseVisitDate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&visitDate, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrWarehouseVisitNotFound
			}
			return nil, err
		}
		locked = append(locked, visitDate)
	}
	return locked, nil
}

func translateWarehouseVisitDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrWarehouseVisitConflict
	}
	return err
}
