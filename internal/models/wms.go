package models

import "time"

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
