package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PesananPembayaran struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PesananID           uuid.UUID       `gorm:"type:uuid;not null" json:"pesanan_id"`
	BuyerID             uuid.UUID       `gorm:"type:uuid;not null" json:"buyer_id"`
	MetodePembayaranID  *uuid.UUID      `gorm:"type:uuid" json:"metode_pembayaran_id"`
	Jumlah              decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"jumlah"`
	Status              PaymentStatus   `gorm:"type:payment_status;not null;default:'PENDING'" json:"status"`
	XenditInvoiceID     *string         `gorm:"type:varchar(100)" json:"xendit_invoice_id"`
	XenditExternalID    *string         `gorm:"type:varchar(100);unique" json:"xendit_external_id"`
	XenditPaymentURL    *string         `gorm:"type:text" json:"xendit_payment_url"`
	XenditPaymentMethod *string         `gorm:"type:varchar(50)" json:"xendit_payment_method"`
	ExpiredAt           *time.Time      `json:"expired_at"`
	PaidAt              *time.Time      `json:"paid_at"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Buyer            Buyer             `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	MetodePembayaran *MetodePembayaran `gorm:"foreignKey:MetodePembayaranID" json:"metode_pembayaran,omitempty"`
}

func (PesananPembayaran) TableName() string {
	return "pesanan_pembayaran"
}

// Request DTOs
type CreatePembayaranRequest struct {
	PesananID          string             `json:"pesanan_id" binding:"required,uuid"`
	MetodePembayaranID string             `json:"metode_pembayaran_id" binding:"required,uuid"`
	SplitPayments      []SplitPaymentItem `json:"split_payments,omitempty"`
}

type SplitPaymentItem struct {
	BuyerID string  `json:"buyer_id" binding:"required,uuid"`
	Jumlah  float64 `json:"jumlah" binding:"required,min=0"`
}

// Response DTOs
type PesananPembayaranResponse struct {
	ID                  string                    `json:"id"`
	PesananID           string                    `json:"pesanan_id"`
	BuyerID             string                    `json:"buyer_id"`
	MetodePembayaranID  *string                   `json:"metode_pembayaran_id"`
	Jumlah              float64                   `json:"jumlah"`
	Status              string                    `json:"status"`
	XenditInvoiceID     *string                   `json:"xendit_invoice_id"`
	XenditExternalID    *string                   `json:"xendit_external_id"`
	XenditPaymentURL    *string                   `json:"xendit_payment_url"`
	XenditPaymentMethod *string                   `json:"xendit_payment_method"`
	ExpiredAt           *time.Time                `json:"expired_at"`
	PaidAt              *time.Time                `json:"paid_at"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
	Buyer               *BuyerResponse            `json:"buyer,omitempty"`
	MetodePembayaran    *MetodePembayaranResponse `json:"metode_pembayaran,omitempty"`
}
