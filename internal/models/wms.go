package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WMSConnectionInfo hasil cek koneksi/identitas client dari WMS
// (GET /api/integration/me). Field Data dibiarkan bertipe interface{} karena
// bentuk persisnya belum didokumentasikan detail oleh tim WMS.
type WMSConnectionInfo struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ========================================
// Sync Produk Palet dari WMS — Cargo Ready-to-Price
// ========================================

// WMSCargoListFilterRequest filter untuk daftar cargo yang siap diberi harga
// (GET /api/integration/cargos/ready-to-price di sisi WMS).
type WMSCargoListFilterRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

func (p *WMSCargoListFilterRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 25
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

// WMSCargoRef referensi ringkas (id + nama) dari master data WMS, mis. kategori,
// kondisi, sumber, brand. Nullable di sisi WMS.
type WMSCargoRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WMSCargoPricingResponse 1 item cargo siap diberi harga (ukuran fisik lengkap,
// belum pernah dihargai, belum disinkronkan) dari WMS.
type WMSCargoPricingResponse struct {
	ID                    string         `json:"id"`
	Code                  string         `json:"code"`
	LengthCM              float64        `json:"length_cm"`
	WidthCM               float64        `json:"width_cm"`
	HeightCM              float64        `json:"height_cm"`
	WeightKG              float64        `json:"weight_kg"`
	TotalPrice            float64        `json:"total_price"`
	BulkyCategory         *WMSCargoRef   `json:"bulky_category"`
	BulkyProductCondition *WMSCargoRef   `json:"bulky_product_condition"`
	BulkyPackageCondition *WMSCargoRef   `json:"bulky_package_condition"`
	BulkyProductSource    *WMSCargoRef   `json:"bulky_product_source"`
	BulkyBrands           []*WMSCargoRef `json:"bulky_brands"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

// wmsCargoListEnvelope bentuk respons mentah dari WMS untuk
// GET /api/integration/cargos/ready-to-price:
// {"success":true,"message":"...","data":[...],"meta":{"page":1,"limit":25,"total_items":1,"total_page":1}}
type WMSCargoListEnvelope struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    []WMSCargoPricingResponse `json:"data"`
	Meta    WMSPaginationMetaRaw      `json:"meta"`
}

// WMSPaginationMetaRaw meta pagination sebagaimana dikembalikan WMS (nama
// field beda dari PaginationMeta internal kita — page/limit, bukan
// current_page/per_page).
type WMSPaginationMetaRaw struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPage  int   `json:"total_page"`
}

// ========================================
// Tetapkan Harga Cargo (discount % atau fix rupiah)
// ========================================

// SetWMSCargoPriceRequest body untuk POST /api/integration/cargos/{id}/price
// di sisi WMS.
type SetWMSCargoPriceRequest struct {
	Type  string  `json:"type" binding:"required,oneof=discount fix"`
	Value float64 `json:"value" binding:"required,gt=0"`
}

// WMSCargoPriceResponse hasil penetapan harga cargo dari WMS.
type WMSCargoPriceResponse struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	PricingType   string    `json:"pricing_type"`
	PricingValue  float64   `json:"pricing_value"`
	TotalPrice    float64   `json:"total_price"`
	SalePrice     float64   `json:"sale_price"`
	PricedAt      time.Time `json:"priced_at"`
	PricingPDFURL string    `json:"pricing_pdf_url"`
}

// wmsCargoPriceEnvelope bentuk respons mentah dari WMS untuk
// POST /api/integration/cargos/{id}/price.
type WMSCargoPriceEnvelope struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    WMSCargoPriceResponse `json:"data"`
}

// wmsErrorEnvelope bentuk respons error dari WMS:
// {"success":false,"message":"...","errors":"..."}
type WMSErrorEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

// ========================================
// Cargo Sudah Diberi Harga, Belum Dikonfirmasi Sinkron (Already-Priced)
// ========================================

// WMSCargoPricedResponse 1 item cargo yang sudah diberi harga di WMS tapi
// belum dikonfirmasi sinkron (belum lewat POST /cargos/{id}/status). Dipakai
// sebagai sumber dropdown "ID Cargo" saat create/edit produk lokal.
type WMSCargoPricedResponse struct {
	ID                    string         `json:"id"`
	Code                  string         `json:"code"`
	LengthCM              float64        `json:"length_cm"`
	WidthCM               float64        `json:"width_cm"`
	HeightCM              float64        `json:"height_cm"`
	WeightKG              float64        `json:"weight_kg"`
	TotalPrice            float64        `json:"total_price"`
	PricingType           string         `json:"pricing_type"`
	PricingValue          float64        `json:"pricing_value"`
	SalePrice             float64        `json:"sale_price"`
	PricedAt              time.Time      `json:"priced_at"`
	PricingPDFURL         string         `json:"pricing_pdf_url"`
	BulkyCategory         *WMSCargoRef   `json:"bulky_category"`
	BulkyProductCondition *WMSCargoRef   `json:"bulky_product_condition"`
	BulkyPackageCondition *WMSCargoRef   `json:"bulky_package_condition"`
	BulkyProductSource    *WMSCargoRef   `json:"bulky_product_source"`
	BulkyBrands           []*WMSCargoRef `json:"bulky_brands"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

// wmsCargoPricedListEnvelope bentuk respons mentah dari WMS untuk
// GET /api/integration/cargos/already-priced.
type WMSCargoPricedListEnvelope struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []WMSCargoPricedResponse `json:"data"`
	Meta    WMSPaginationMetaRaw     `json:"meta"`
}

// ========================================
// Tandai Cargo Sudah Dikonfirmasi Sinkron (is_sync = true)
// ========================================

// WMSCargoSyncStatusResponse hasil penandaan cargo sudah disinkron dari WMS.
type WMSCargoSyncStatusResponse struct {
	ID       string    `json:"id"`
	Code     string    `json:"code"`
	IsSync   bool      `json:"is_sync"`
	SyncedAt time.Time `json:"synced_at"`
}

// wmsCargoSyncStatusEnvelope bentuk respons mentah dari WMS untuk
// POST /api/integration/cargos/{id}/status.
type WMSCargoSyncStatusEnvelope struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    WMSCargoSyncStatusResponse `json:"data"`
}

// ========================================
// Cache Lokal Cargo WMS Sudah Diberi Harga (wms_cargo_priced)
// ========================================

// WMSCargoPriced cache lokal 1 cargo WMS yang sudah diberi harga jual —
// jembatan antara WMS dan produk lokal. Diisi otomatis saat admin menekan
// "Simpan Harga Jual" di modal sync (termasuk PDF harga yang sudah
// di-download & disimpan ke storage lokal). Dropdown "ID Cargo" saat
// create/edit produk membaca dari tabel ini.
type WMSCargoPriced struct {
	CargoID               uuid.UUID       `gorm:"type:uuid;primary_key;column:cargo_id" json:"cargo_id"`
	Code                  string          `gorm:"type:varchar(100);not null" json:"code"`
	LengthCM              float64         `gorm:"type:decimal(10,2);not null;default:0" json:"length_cm"`
	WidthCM               float64         `gorm:"type:decimal(10,2);not null;default:0" json:"width_cm"`
	HeightCM              float64         `gorm:"type:decimal(10,2);not null;default:0" json:"height_cm"`
	WeightKG              float64         `gorm:"type:decimal(10,2);not null;default:0" json:"weight_kg"`
	TotalPrice            float64         `gorm:"type:decimal(15,2);not null;default:0" json:"total_price"`
	PricingType           *string         `gorm:"type:varchar(20)" json:"pricing_type"`
	PricingValue          *float64        `gorm:"type:decimal(15,2)" json:"pricing_value"`
	SalePrice             *float64        `gorm:"type:decimal(15,2)" json:"sale_price"`
	PricedAt              *time.Time      `gorm:"type:timestamptz" json:"priced_at"`
	PricingPDFPath        *string         `gorm:"type:varchar(255);column:pricing_pdf_path" json:"pricing_pdf_path"`
	BulkyCategory         json.RawMessage `gorm:"type:jsonb" json:"bulky_category,omitempty"`
	BulkyProductCondition json.RawMessage `gorm:"type:jsonb" json:"bulky_product_condition,omitempty"`
	BulkyPackageCondition json.RawMessage `gorm:"type:jsonb" json:"bulky_package_condition,omitempty"`
	BulkyProductSource    json.RawMessage `gorm:"type:jsonb" json:"bulky_product_source,omitempty"`
	BulkyBrands           json.RawMessage `gorm:"type:jsonb" json:"bulky_brands,omitempty"`
	IsUsedInProduk        bool            `gorm:"not null;default:false;column:is_used_in_produk" json:"is_used_in_produk"`
	ProdukID              *uuid.UUID      `gorm:"type:uuid;column:produk_id" json:"produk_id"`
	CreatedAt             time.Time       `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time       `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
}

func (WMSCargoPriced) TableName() string {
	return "wms_cargo_priced"
}

// WMSCargoPricedListResponse 1 item cargo WMS yang sudah di-price, siap
// dipilih di dropdown "ID Cargo" saat create/edit produk lokal.
type WMSCargoPricedListResponse struct {
	CargoID        string  `json:"cargo_id"`
	Code           string  `json:"code"`
	LengthCM       float64 `json:"length_cm"`
	WidthCM        float64 `json:"width_cm"`
	HeightCM       float64 `json:"height_cm"`
	WeightKG       float64 `json:"weight_kg"`
	TotalPrice     float64 `json:"total_price"`
	SalePrice      float64 `json:"sale_price"`
	PricingPDFURL  *string `json:"pricing_pdf_url"`
	IsUsedInProduk bool    `json:"is_used_in_produk"`
}
