package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DeliveryType string
type PaymentType string
type PaymentStatus string
type OrderStatus string

const (
	DeliveryTypePickup    DeliveryType = "PICKUP"
	DeliveryTypeDeliveree DeliveryType = "DELIVEREE"
	DeliveryTypeForwarder DeliveryType = "FORWARDER"

	PaymentTypeRegular PaymentType = "REGULAR"
	PaymentTypeSplit   PaymentType = "SPLIT"

	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPartial  PaymentStatus = "PARTIAL"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusExpired  PaymentStatus = "EXPIRED"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"

	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusReady      OrderStatus = "READY"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type Pesanan struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Kode                string          `gorm:"type:varchar(20);not null;unique" json:"kode"`
	BuyerID             uuid.UUID       `gorm:"type:uuid;not null" json:"buyer_id"`
	DeliveryType        DeliveryType    `gorm:"type:delivery_type;not null" json:"delivery_type"`
	AlamatBuyerID       *uuid.UUID      `gorm:"type:uuid" json:"alamat_buyer_id"`
	PaymentType         PaymentType     `gorm:"type:payment_type;not null;default:'REGULAR'" json:"payment_type"`
	PaymentStatus       PaymentStatus   `gorm:"type:payment_status;not null;default:'PENDING'" json:"payment_status"`
	OrderStatus         OrderStatus     `gorm:"type:order_status;not null;default:'PENDING'" json:"order_status"`
	BiayaProduk         decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0" json:"biaya_produk"`
	BiayaPengiriman     decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0" json:"biaya_pengiriman"`
	BiayaPPN            decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0" json:"biaya_ppn"`
	BiayaLainnya        decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0" json:"biaya_lainnya"`
	Total               decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0" json:"total"`
	Catatan             *string         `gorm:"type:text" json:"catatan"`
	CatatanAdmin        *string         `gorm:"type:text" json:"catatan_admin"`
	ExpiredAt           *time.Time      `gorm:"type:timestamptz" json:"expired_at"`
	PaidAt              *time.Time      `gorm:"type:timestamptz" json:"paid_at"`
	ProcessedAt         *time.Time      `gorm:"type:timestamptz" json:"processed_at"`
	ReadyAt             *time.Time      `gorm:"type:timestamptz" json:"ready_at"`
	ShippedAt           *time.Time      `gorm:"type:timestamptz" json:"shipped_at"`
	CompletedAt         *time.Time      `gorm:"type:timestamptz" json:"completed_at"`
	CancelledAt         *time.Time      `gorm:"type:timestamptz" json:"cancelled_at"`
	CancelledReason     *string         `gorm:"type:text" json:"cancelled_reason"`
	DelivereeBookingID  *string         `gorm:"type:varchar(100)" json:"deliveree_booking_id"`
	ForwarderTrackingNo *string         `gorm:"type:varchar(100)" json:"forwarder_tracking_no"`
	CreatedAt           time.Time       `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt  `gorm:"type:timestamptz;index" json:"-"`

	// Relations
	Buyer       Buyer               `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	AlamatBuyer *AlamatBuyer        `gorm:"foreignKey:AlamatBuyerID" json:"alamat_buyer,omitempty"`
	Items       []PesananItem       `gorm:"foreignKey:PesananID" json:"items,omitempty"`
	Pembayaran  []PesananPembayaran `gorm:"foreignKey:PesananID" json:"pembayaran,omitempty"`
}

func (Pesanan) TableName() string {
	return "pesanan"
}

// Request DTOs
type CreatePesananRequest struct {
	BuyerID         string                     `json:"buyer_id" binding:"required,uuid"`
	DeliveryType    string                     `json:"delivery_type" binding:"required,oneof=PICKUP DELIVEREE FORWARDER"`
	AlamatBuyerID   *string                    `json:"alamat_buyer_id" binding:"omitempty,uuid"`
	PaymentType     string                     `json:"payment_type" binding:"required,oneof=REGULAR SPLIT"`
	Items           []CreatePesananItemRequest `json:"items" binding:"required,min=1,dive"`
	BiayaPengiriman float64                    `json:"biaya_pengiriman"`
	BiayaLainnya    float64                    `json:"biaya_lainnya"`
	Catatan         *string                    `json:"catatan"`
}

type UpdatePesananRequest struct {
	OrderStatus         *string `json:"order_status" binding:"omitempty,oneof=PENDING PROCESSING READY SHIPPED COMPLETED CANCELLED"`
	CatatanAdmin        *string `json:"catatan_admin"`
	CancelledReason     *string `json:"cancelled_reason"`
	DelivereeBookingID  *string `json:"deliveree_booking_id"`
	ForwarderTrackingNo *string `json:"forwarder_tracking_no"`
}

// Response DTOs
type PesananResponse struct {
	ID                  string                      `json:"id"`
	Kode                string                      `json:"kode"`
	BuyerID             string                      `json:"buyer_id"`
	DeliveryType        string                      `json:"delivery_type"`
	AlamatBuyerID       *string                     `json:"alamat_buyer_id"`
	PaymentType         string                      `json:"payment_type"`
	PaymentStatus       string                      `json:"payment_status"`
	OrderStatus         string                      `json:"order_status"`
	BiayaProduk         float64                     `json:"biaya_produk"`
	BiayaPengiriman     float64                     `json:"biaya_pengiriman"`
	BiayaPPN            float64                     `json:"biaya_ppn"`
	BiayaLainnya        float64                     `json:"biaya_lainnya"`
	Total               float64                     `json:"total"`
	Catatan             *string                     `json:"catatan"`
	CatatanAdmin        *string                     `json:"catatan_admin"`
	ExpiredAt           *time.Time                  `json:"expired_at"`
	PaidAt              *time.Time                  `json:"paid_at"`
	ProcessedAt         *time.Time                  `json:"processed_at"`
	ReadyAt             *time.Time                  `json:"ready_at"`
	ShippedAt           *time.Time                  `json:"shipped_at"`
	CompletedAt         *time.Time                  `json:"completed_at"`
	CancelledAt         *time.Time                  `json:"cancelled_at"`
	CancelledReason     *string                     `json:"cancelled_reason"`
	DelivereeBookingID  *string                     `json:"deliveree_booking_id"`
	ForwarderTrackingNo *string                     `json:"forwarder_tracking_no"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
	Buyer               *BuyerResponse              `json:"buyer,omitempty"`
	AlamatBuyer         *AlamatBuyerResponse        `json:"alamat_buyer,omitempty"`
	Items               []PesananItemResponse       `json:"items,omitempty"`
	Pembayaran          []PesananPembayaranResponse `json:"pembayaran,omitempty"`
}
