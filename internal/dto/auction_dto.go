package dto

import (
	"time"
)

// ============================================================
// Query params
// ============================================================

// AuctionListQueryParams parameter list batch (GET /api/panel/auctions).
type AuctionListQueryParams struct {
	Page     int    `query:"page"`
	PerPage  int    `query:"per_page"`
	Search   string `query:"search"`
	Status   string `query:"status"`
	SortBy   string `query:"sort_by"` // created_at_desc | created_at_asc
	SortDesc bool   `query:"sort_desc"`
}

func (p *AuctionListQueryParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 20
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at_desc"
	}
}

// AuctionBidsQueryParams parameter list bid (GET /api/panel/auctions/{id}/bids).
type AuctionBidsQueryParams struct {
	Page     int    `query:"page"`
	PerPage  int    `query:"per_page"`
	Search   string `query:"search"`
	BuyerID  string `query:"buyer_id"`
	SortBy   string `query:"sort_by"` // created_at_desc | amount_desc | amount_asc
	SortDesc bool   `query:"sort_desc"`
}

func (p *AuctionBidsQueryParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 20
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at_desc"
	}
}

// AuctionProductOptionsQueryParams parameter kandidat produk (GET /product-options).
type AuctionProductOptionsQueryParams struct {
	Page        int    `query:"page"`
	PerPage     int    `query:"per_page"`
	Search      string `query:"search"`
	WarehouseID string `query:"warehouse_id"`
	SortBy      string `query:"sort_by"`
	SortDesc    bool   `query:"sort_desc"`
}

func (p *AuctionProductOptionsQueryParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 20
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at_desc"
	}
}

// ============================================================
// Request DTOs
// ============================================================

// AuctionDraftInput payload create/edit draft batch. Field turunan
// (status, grand_total, min_bid_percent, winner, created_by) tidak boleh
// ditulis client; server menghitungnya.
type AuctionDraftInput struct {
	Version               int                `json:"version" validate:"omitempty,min=1"`
	NamaID                string             `json:"nama_id" validate:"required,max=255"`
	NamaEN                *string            `json:"nama_en" validate:"omitempty,max=255"`
	Description           *string            `json:"description" validate:"omitempty,max=5000"`
	WarehouseID           *string            `json:"warehouse_id"`
	OriginType            string             `json:"origin_type" validate:"omitempty,oneof=BULKY_WAREHOUSE SUPPLIER"`
	SupplierName          *string            `json:"supplier_name" validate:"omitempty,max=255"`
	SupplierAddress       *string            `json:"supplier_address" validate:"omitempty,max=1000"`
	SupplierProvinsi      *string            `json:"supplier_provinsi" validate:"omitempty,max=100"`
	SupplierKota          *string            `json:"supplier_kota" validate:"omitempty,max=100"`
	SupplierKecamatan     *string            `json:"supplier_kecamatan" validate:"omitempty,max=100"`
	SupplierKelurahan     *string            `json:"supplier_kelurahan" validate:"omitempty,max=100"`
	SupplierKodePos       *string            `json:"supplier_kode_pos" validate:"omitempty,max=10"`
	SupplierLatitude      *string            `json:"supplier_latitude"`
	SupplierLongitude     *string            `json:"supplier_longitude"`
	KategoriID            *string            `json:"kategori_id"`
	KondisiID             *string            `json:"kondisi_id"`
	KondisiPaketID        *string            `json:"kondisi_paket_id"`
	SumberID              *string            `json:"sumber_id"`
	DiscrepancyPercentage string             `json:"discrepancy_percentage"`
	MerekIDs              []string           `json:"merek_ids"`
	Items                 []AuctionDraftItem `json:"items"`
	PanjangCm             string             `json:"panjang_cm"`
	LebarCm               string             `json:"lebar_cm"`
	TinggiCm              string             `json:"tinggi_cm"`
	BeratKg               string             `json:"berat_kg"`
	ImageAssetIDs         []string           `json:"image_asset_ids"`
	PDFAssetID            *string            `json:"pdf_asset_id"`
}

// AuctionDraftItem dapat merujuk katalog Bulky atau snapshot manual. Item
// MANUAL dipakai untuk barang supplier/consignment yang tidak ada di katalog.
type AuctionDraftItem struct {
	SourceType string `json:"source_type" validate:"omitempty,oneof=CATALOG MANUAL"`
	ProdukID   string `json:"produk_id" validate:"omitempty,uuid"`
	Nama       string `json:"nama" validate:"omitempty,max=255"`
	UnitPrice  string `json:"unit_price" validate:"omitempty"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
}

// AuctionPublishRequest payload publish batch. Wajib membawa version terakhir.
type AuctionPublishRequest struct {
	Version int `json:"version" validate:"required,min=1"`
}

// AuctionWinnerRequest payload konfirmasi winner. Wajib Idempotency-Key header.
type AuctionWinnerRequest struct {
	BidID   string  `json:"bid_id" validate:"required,uuid"`
	Version int     `json:"version" validate:"required,min=1"`
	Note    *string `json:"note" validate:"omitempty,max=1000"`
}

// AuctionOperationRequest payload PATCH /operations. Menerima tepat satu target
// status: payment_status ATAU fulfillment_status.
type AuctionOperationRequest struct {
	Version           int     `json:"version" validate:"required,min=1"`
	PaymentStatus     *string `json:"payment_status" validate:"omitempty,oneof=UNPAID PAID"`
	FulfillmentStatus *string `json:"fulfillment_status" validate:"omitempty,oneof=PENDING PROCESSING COMPLETED"`
	Note              *string `json:"note" validate:"omitempty,max=1000"`
}

// AuctionAssetUploadRequest payload upload aset (multipart file + kind).
type AuctionAssetUploadRequest struct {
	Kind string `form:"kind" validate:"required,oneof=IMAGE PDF"`
}

// AuctionSupplierExcelPreview previews the header row so the admin can map
// workbook columns before importing supplier items.
type AuctionSupplierExcelPreview struct {
	SheetName string                       `json:"sheet_name"`
	HeaderRow int                          `json:"header_row"`
	Columns   []AuctionSupplierExcelColumn `json:"columns"`
}

type AuctionSupplierExcelColumn struct {
	Index   int      `json:"index"`
	Letter  string   `json:"letter"`
	Header  string   `json:"header"`
	Samples []string `json:"samples"`
}

type AuctionSupplierExcelMapping struct {
	NameColumn     int `form:"name_column"`
	PriceColumn    int `form:"price_column"`
	QuantityColumn int `form:"quantity_column"`
	HeaderRow      int `form:"header_row"`
}

type AuctionSupplierExcelImport struct {
	Items []AuctionDraftItem    `json:"items"`
	PDF   *AuctionAssetResponse `json:"pdf"`
}

// ============================================================
// Response DTOs
// ============================================================

// AuctionBatchSummary ringkasan batch untuk list.
type AuctionBatchSummary struct {
	ID           string     `json:"id"`
	Code         string     `json:"code"`
	NamaID       string     `json:"nama_id"`
	ThumbnailURL *string    `json:"thumbnail_url"`
	Status       string     `json:"status"`
	GrandTotal   string     `json:"grand_total"`
	MinBidAmount string     `json:"min_bid_amount"`
	BidderCount  int64      `json:"bidder_count"`
	BidCount     int64      `json:"bid_count"`
	HighestBid   *string    `json:"highest_bid"`
	CreatedAt    time.Time  `json:"created_at"`
	OpenedAt     *time.Time `json:"opened_at"`
	SoldAt       *time.Time `json:"sold_at"`
	Version      int        `json:"version"`
}

// AuctionBatchDetail detail lengkap batch.
type AuctionBatchDetail struct {
	ID                    string                 `json:"id"`
	Code                  string                 `json:"code"`
	NamaID                string                 `json:"nama_id"`
	NamaEN                *string                `json:"nama_en"`
	Description           *string                `json:"description"`
	WarehouseID           *string                `json:"warehouse_id"`
	OriginType            string                 `json:"origin_type"`
	SupplierName          *string                `json:"supplier_name"`
	SupplierAddress       *string                `json:"supplier_address"`
	SupplierProvinsi      *string                `json:"supplier_provinsi"`
	SupplierKota          *string                `json:"supplier_kota"`
	SupplierKecamatan     *string                `json:"supplier_kecamatan"`
	SupplierKelurahan     *string                `json:"supplier_kelurahan"`
	SupplierKodePos       *string                `json:"supplier_kode_pos"`
	SupplierLatitude      *string                `json:"supplier_latitude"`
	SupplierLongitude     *string                `json:"supplier_longitude"`
	KategoriID            *string                `json:"kategori_id"`
	KondisiID             *string                `json:"kondisi_id"`
	KondisiPaketID        *string                `json:"kondisi_paket_id"`
	SumberID              *string                `json:"sumber_id"`
	DiscrepancyPercentage string                 `json:"discrepancy_percentage"`
	Status                string                 `json:"status"`
	GrandTotal            string                 `json:"grand_total"`
	MinBidPercent         string                 `json:"min_bid_percent"`
	MinBidAmount          string                 `json:"min_bid_amount"`
	TotalQuantity         int                    `json:"total_quantity"`
	PanjangCm             string                 `json:"panjang_cm"`
	LebarCm               string                 `json:"lebar_cm"`
	TinggiCm              string                 `json:"tinggi_cm"`
	BeratKg               string                 `json:"berat_kg"`
	VolumeM3              string                 `json:"volume_m3"`
	Version               int                    `json:"version"`
	CreatedBy             string                 `json:"created_by"`
	UpdatedBy             string                 `json:"updated_by"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	OpenedAt              *time.Time             `json:"opened_at"`
	SoldAt                *time.Time             `json:"sold_at"`
	Items                 []AuctionItemSnapshot  `json:"items"`
	BrandIDs              []string               `json:"merek_ids"`
	Images                []AuctionAssetResponse `json:"images"`
	PDF                   *AuctionAssetResponse  `json:"pdf"`
	Winner                *AuctionWinnerResponse `json:"winner"`
	BidderCount           int64                  `json:"bidder_count"`
	BidCount              int64                  `json:"bid_count"`
	HighestBid            *string                `json:"highest_bid"`
	LowestBid             *string                `json:"lowest_bid"`
	ViewCount             int64                  `json:"view_count"`
	RepeatBidderCount     int64                  `json:"repeat_bidder_count"`
	AdditionalBidCount    int64                  `json:"additional_bid_count"`
	TimeToSoldSeconds     *int64                 `json:"time_to_sold_seconds"`
}

// AuctionItemSnapshot snapshot item batch.
type AuctionItemSnapshot struct {
	ProdukID          *string `json:"produk_id"`
	SourceType        string  `json:"source_type"`
	NamaSnapshot      string  `json:"nama_snapshot"`
	Quantity          int     `json:"quantity"`
	UnitPriceSnapshot string  `json:"unit_price_snapshot"`
	SubtotalSnapshot  string  `json:"subtotal_snapshot"`
}

// AuctionBidDetail detail bid untuk panel.
type AuctionBidDetail struct {
	ID                 string             `json:"id"`
	BatchID            string             `json:"batch_id"`
	Buyer              AuctionBuyerSimple `json:"buyer"`
	Sequence           int                `json:"sequence"`
	InputMode          string             `json:"input_mode"`
	InputPercent       *string            `json:"input_percent"`
	EffectivePercent   string             `json:"effective_percent"`
	Amount             string             `json:"amount"`
	GrandTotalSnapshot string             `json:"grand_total_snapshot"`
	CreatedAt          time.Time          `json:"created_at"`
	IsSelected         bool               `json:"is_selected"`
}

// AuctionBuyerSimple info buyer pada bid.
type AuctionBuyerSimple struct {
	ID      string `json:"id"`
	Nama    string `json:"nama"`
	Telepon string `json:"telepon"`
}

// AuctionWinnerResponse info winner pada detail batch.
type AuctionWinnerResponse struct {
	BidID             string             `json:"bid_id"`
	Buyer             AuctionBuyerSimple `json:"buyer"`
	DealAmount        string             `json:"deal_amount"`
	SelectedBy        string             `json:"selected_by"`
	SelectedAt        time.Time          `json:"selected_at"`
	Note              *string            `json:"note"`
	PaymentStatus     string             `json:"payment_status"`
	PaidAt            *time.Time         `json:"paid_at"`
	PaymentNote       *string            `json:"payment_note"`
	FulfillmentStatus string             `json:"fulfillment_status"`
	CompletedAt       *time.Time         `json:"completed_at"`
	FulfillmentNote   *string            `json:"fulfillment_note"`
}

// AuctionAssetResponse aset batch.
type AuctionAssetResponse struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	URL          string `json:"url"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	SizeBytes    int64  `json:"size_bytes"`
}

// AuctionProductOption kandidat produk untuk pemilih produk.
type AuctionProductOption struct {
	ID             string `json:"id"`
	NamaID         string `json:"nama_id"`
	WarehouseID    string `json:"warehouse_id"`
	UnitPrice      string `json:"unit_price"`
	StockAvailable int    `json:"stock_available"`
	IsActive       bool   `json:"is_active"`
	IsSold         bool   `json:"is_sold"`
}
