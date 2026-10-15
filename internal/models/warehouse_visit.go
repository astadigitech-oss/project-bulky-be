package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	WarehouseVisitBookingConfirmed = "CONFIRMED"
	WarehouseVisitBookingCancelled = "CANCELLED"
	WarehouseVisitBookingAttended  = "ATTENDED"
	WarehouseVisitBookingNoShow    = "NO_SHOW"
)

type WarehouseVisitDate struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	WarehouseID    uuid.UUID  `gorm:"type:uuid;not null" json:"warehouse_id"`
	VisitDate      time.Time  `gorm:"type:date;not null" json:"visit_date"`
	IsOpen         bool       `gorm:"not null;default:true" json:"is_open"`
	CrewCount      int        `gorm:"not null;default:0" json:"crew_count"`
	VisitHoursNote *string    `gorm:"type:text" json:"visit_hours_note,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (WarehouseVisitDate) TableName() string { return "warehouse_visit_dates" }

type WarehouseVisitBooking struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	VisitDateID uuid.UUID `gorm:"type:uuid;not null" json:"visit_date_id"`
	VisitDate   time.Time `gorm:"type:date;not null" json:"visit_date"`
	BuyerID     uuid.UUID `gorm:"type:uuid;not null" json:"buyer_id"`
	Status      string    `gorm:"type:varchar(20);not null;default:CONFIRMED" json:"status"`
	CreatedAt   time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (WarehouseVisitBooking) TableName() string { return "warehouse_visit_bookings" }

type WarehouseVisitDateInput struct {
	VisitDate      string  `json:"visit_date" binding:"required"`
	IsOpen         bool    `json:"is_open"`
	CrewCount      int     `json:"crew_count" binding:"min=0"`
	VisitHoursNote *string `json:"visit_hours_note"`
}

type WarehouseVisitDateUpdateInput struct {
	IsOpen         bool    `json:"is_open"`
	CrewCount      int     `json:"crew_count" binding:"min=0"`
	VisitHoursNote *string `json:"visit_hours_note"`
}

type WarehouseVisitDateResponse struct {
	ID                string  `json:"id"`
	WarehouseID       string  `json:"warehouse_id"`
	VisitDate         string  `json:"visit_date"`
	IsOpen            bool    `json:"is_open"`
	CrewCount         int     `json:"crew_count"`
	Capacity          int     `json:"capacity"`
	BookedCount       int64   `json:"booked_count"`
	RemainingCapacity int64   `json:"remaining_capacity"`
	VisitHoursNote    *string `json:"visit_hours_note,omitempty"`
	DisplayHours      *string `json:"display_hours,omitempty"`
}

type WarehouseVisitBookingListRequest struct {
	PaginationRequest
	FromDate string `query:"from"`
	ToDate   string `query:"to"`
	Status   string `query:"status"`
}

type WarehouseVisitBookingResponse struct {
	ID          string    `json:"id"`
	VisitDateID string    `json:"visit_date_id"`
	VisitDate   string    `json:"visit_date"`
	BuyerID     string    `json:"buyer_id"`
	BuyerName   string    `json:"buyer_name"`
	BuyerEmail  string    `json:"buyer_email"`
	BuyerPhone  *string   `json:"buyer_phone,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WarehouseVisitRescheduleRequest struct {
	VisitDateID string `json:"visit_date_id" binding:"required,uuid"`
}

type WarehouseVisitAttendanceRequest struct {
	Status string `json:"status" binding:"required,oneof=ATTENDED NO_SHOW"`
}

type WarehouseVisitBookingFilter struct {
	FromDate *time.Time
	ToDate   *time.Time
	Status   string
	Page     int
	PerPage  int
}
