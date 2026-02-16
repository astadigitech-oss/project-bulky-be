package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PesananAdminQueryParams query parameters for admin pesanan list
type PesananAdminQueryParams struct {
	Page          int    `form:"page" binding:"required,min=1"`
	PerPage       int    `form:"per_page" binding:"required,min=1,max=100"`
	Cari          string `form:"cari"`
	OrderStatus   string `form:"order_status" binding:"omitempty,oneof=PENDING PROCESSING READY SHIPPED COMPLETED CANCELLED"`
	PaymentStatus string `form:"payment_status" binding:"omitempty,oneof=PENDING PAID EXPIRED FAILED"`
	DeliveryType  string `form:"delivery_type" binding:"omitempty,oneof=PICKUP DELIVEREE FORWARDER"`
	TanggalDari   string `form:"tanggal_dari"`
	TanggalSampai string `form:"tanggal_sampai"`
	SortBy        string `form:"sort_by"`
	SortOrder     string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for query params
func (p *PesananAdminQueryParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 10
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}

// PesananAdminListResponse simple response for pesanan list (admin)
type PesananAdminListResponse struct {
	ID              uuid.UUID                 `json:"id"`
	Kode            string                    `json:"kode"`
	Buyer           PesananAdminBuyerResponse `json:"buyer"`
	DeliveryType    string                    `json:"delivery_type"`
	PaymentType     string                    `json:"payment_type"`
	PaymentStatus   string                    `json:"payment_status"`
	OrderStatus     string                    `json:"order_status"`
	TotalItem       int                       `json:"total_item"`
	BiayaProduk     decimal.Decimal           `json:"biaya_produk"`
	BiayaPengiriman decimal.Decimal           `json:"biaya_pengiriman"`
	BiayaPPN        decimal.Decimal           `json:"biaya_ppn"`
	TotalBayar      decimal.Decimal           `json:"total_bayar"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// PesananAdminDetailResponse detailed response for pesanan (admin)
type PesananAdminDetailResponse struct {
	ID               uuid.UUID                           `json:"id"`
	Kode             string                              `json:"kode"`
	Buyer            PesananAdminBuyerDetailResponse     `json:"buyer"`
	AlamatPengiriman *PesananAdminAlamatResponse         `json:"alamat_pengiriman"`
	DeliveryType     string                              `json:"delivery_type"`
	PaymentType      string                              `json:"payment_type"`
	PaymentStatus    string                              `json:"payment_status"`
	OrderStatus      string                              `json:"order_status"`
	Items            []PesananAdminItemResponse          `json:"items"`
	Pembayaran       []PesananAdminPembayaranResponse    `json:"pembayaran"`
	StatusHistory    []PesananAdminStatusHistoryResponse `json:"status_history"`
	BiayaProduk      decimal.Decimal                     `json:"biaya_produk"`
	BiayaPengiriman  decimal.Decimal                     `json:"biaya_pengiriman"`
	BiayaPPN         decimal.Decimal                     `json:"biaya_ppn"`
	PotonganKupon    decimal.Decimal                     `json:"potongan_kupon"`
	TotalBayar       decimal.Decimal                     `json:"total_bayar"`
	CatatanBuyer     *string                             `json:"catatan_buyer"`
	CatatanAdmin     *string                             `json:"catatan_admin"`
	CreatedAt        time.Time                           `json:"created_at"`
	UpdatedAt        time.Time                           `json:"updated_at"`
}

// PesananAdminBuyerResponse buyer info for list
type PesananAdminBuyerResponse struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Email string    `json:"email"`
}

// PesananAdminBuyerDetailResponse buyer info with phone for detail
type PesananAdminBuyerDetailResponse struct {
	ID      uuid.UUID `json:"id"`
	Nama    string    `json:"nama"`
	Email   string    `json:"email"`
	Telepon string    `json:"telepon"`
}

// PesananAdminAlamatResponse alamat info
type PesananAdminAlamatResponse struct {
	ID            uuid.UUID `json:"id"`
	Label         string    `json:"label"`
	NamaPenerima  string    `json:"nama_penerima"`
	Telepon       string    `json:"telepon"`
	AlamatLengkap string    `json:"alamat_lengkap"`
	Kota          string    `json:"kota"`
	Provinsi      string    `json:"provinsi"`
	KodePos       string    `json:"kode_pos"`
}

// PesananAdminItemResponse pesanan item info
type PesananAdminItemResponse struct {
	ID           uuid.UUID                      `json:"id"`
	Produk       PesananAdminItemProdukResponse `json:"produk"`
	NamaProduk   string                         `json:"nama_produk"`
	Qty          int                            `json:"qty"`
	HargaSatuan  decimal.Decimal                `json:"harga_satuan"`
	DiskonSatuan decimal.Decimal                `json:"diskon_satuan"`
	Subtotal     decimal.Decimal                `json:"subtotal"`
}

// PesananAdminItemProdukResponse produk info in item
type PesananAdminItemProdukResponse struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	GambarURL *string   `json:"gambar_url"`
}

// PesananAdminPembayaranResponse pembayaran info
type PesananAdminPembayaranResponse struct {
	ID               uuid.UUID                            `json:"id"`
	BuyerID          uuid.UUID                            `json:"buyer_id"`
	NamaPembayar     string                               `json:"nama_pembayar"`
	Jumlah           decimal.Decimal                      `json:"jumlah"`
	MetodePembayaran PesananAdminMetodePembayaranResponse `json:"metode_pembayaran"`
	Status           string                               `json:"status"`
	PaidAt           *time.Time                           `json:"paid_at"`
	XenditInvoiceID  *string                              `json:"xendit_invoice_id"`
}

// PesananAdminMetodePembayaranResponse metode pembayaran info
type PesananAdminMetodePembayaranResponse struct {
	ID   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
	Kode string    `json:"kode"`
}

// PesananAdminStatusHistoryResponse status history info
type PesananAdminStatusHistoryResponse struct {
	StatusFrom *string   `json:"status_from"`
	StatusTo   string    `json:"status_to"`
	StatusType string    `json:"status_type"`
	Note       *string   `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

// UpdatePesananStatusRequest request for update pesanan status
type UpdatePesananStatusRequest struct {
	OrderStatus string  `json:"order_status" binding:"required,oneof=PROCESSING READY SHIPPED COMPLETED CANCELLED"`
	Note        *string `json:"note" binding:"omitempty,max=500"`
}

// UpdatePesananStatusResponse response after update status
type UpdatePesananStatusResponse struct {
	ID             uuid.UUID `json:"id"`
	Kode           string    `json:"kode"`
	OrderStatus    string    `json:"order_status"`
	PreviousStatus string    `json:"previous_status"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
}

// PesananStatisticsResponse response for pesanan statistics
type PesananStatisticsResponse struct {
	TotalPesanan     int64            `json:"total_pesanan"`
	TotalRevenue     decimal.Decimal  `json:"total_revenue"`
	PerStatus        map[string]int64 `json:"per_status"`
	PerDeliveryType  map[string]int64 `json:"per_delivery_type"`
	PerPaymentStatus map[string]int64 `json:"per_payment_status"`
}
