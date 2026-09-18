package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// JSONMap adalah map[string]interface{} yang dapat dibaca/ditulis sebagai JSONB
// di PostgreSQL. GORM tidak dapat mem-scan jsonb langsung ke map, sehingga perlu
// implementasi sql.Scanner & driver.Valuer.
type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Status batch lelang.
const (
	AuctionBatchStatusDRAFT = "DRAFT"
	AuctionBatchStatusOPEN  = "OPEN"
	AuctionBatchStatusSOLD  = "SOLD"
)

// Input mode bid.
const (
	AuctionBidModeAmount  = "AMOUNT"
	AuctionBidModePercent = "PERCENT"
)

// AuctionBatch adalah master batch lelang. Metadata snapshot (kategori,
// warehouse, fisik, total) dibekukan saat publish. Status hanya DRAFT → OPEN → SOLD.
type AuctionBatch struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code                  string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	NamaID                string          `gorm:"type:varchar(255);not null" json:"nama_id"`
	NamaEN                *string         `gorm:"type:varchar(255)" json:"nama_en"`
	Description           *string         `gorm:"type:text" json:"description"`
	WarehouseID           *uuid.UUID      `gorm:"type:uuid" json:"warehouse_id"`
	KategoriID            *uuid.UUID      `gorm:"type:uuid" json:"kategori_id"`
	KondisiID             *uuid.UUID      `gorm:"type:uuid" json:"kondisi_id"`
	KondisiPaketID        *uuid.UUID      `gorm:"type:uuid" json:"kondisi_paket_id"`
	SumberID              *uuid.UUID      `gorm:"type:uuid" json:"sumber_id"`
	DiscrepancyPercentage decimal.Decimal `gorm:"type:numeric(5,2);not null;default:0" json:"discrepancy_percentage"`
	Status                string          `gorm:"type:varchar(20);not null;default:DRAFT" json:"status"`
	GrandTotal            decimal.Decimal `gorm:"type:numeric(18,0);not null;default:0" json:"grand_total"`
	MinBidPercent         decimal.Decimal `gorm:"type:numeric(12,4);not null;default:0.1" json:"min_bid_percent"`
	TotalQuantity         int             `gorm:"not null;default:0" json:"total_quantity"`
	PanjangCm             decimal.Decimal `gorm:"type:numeric(12,3);not null;default:0" json:"panjang_cm"`
	LebarCm               decimal.Decimal `gorm:"type:numeric(12,3);not null;default:0" json:"lebar_cm"`
	TinggiCm              decimal.Decimal `gorm:"type:numeric(12,3);not null;default:0" json:"tinggi_cm"`
	BeratKg               decimal.Decimal `gorm:"type:numeric(12,3);not null;default:0" json:"berat_kg"`
	VolumeM3              decimal.Decimal `gorm:"type:numeric(12,3);not null;default:0" json:"volume_m3"`
	Version               int             `gorm:"not null;default:1" json:"version"`
	CreatedBy             uuid.UUID       `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy             uuid.UUID       `gorm:"type:uuid;not null" json:"updated_by"`
	CreatedAt             time.Time       `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time       `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	OpenedAt              *time.Time      `gorm:"type:timestamptz" json:"opened_at"`
	SoldAt                *time.Time      `gorm:"type:timestamptz" json:"sold_at"`

	// Relations
	Items        []AuctionBatchItem        `gorm:"foreignKey:BatchID" json:"items,omitempty"`
	Brands       []AuctionBatchBrand       `gorm:"foreignKey:BatchID" json:"brands,omitempty"`
	BatchAssets  []AuctionBatchAsset       `gorm:"foreignKey:BatchID" json:"batch_assets,omitempty"`
	Winner       *AuctionWinner            `gorm:"foreignKey:BatchID" json:"winner,omitempty"`
	Reservations []AuctionStockReservation `gorm:"foreignKey:BatchID" json:"-"`
}

func (AuctionBatch) TableName() string {
	return "auction_batches"
}

// AuctionBatchItem adalah item produk di dalam batch beserta snapshot
// nama/harga pada saat draft disimpan. Snapshot dibekukan saat publish.
type AuctionBatchItem struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID           uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex:idx_batch_produk" json:"batch_id"`
	ProdukID          uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex:idx_batch_produk" json:"produk_id"`
	Quantity          int             `gorm:"not null" json:"quantity"`
	NamaSnapshot      string          `gorm:"type:varchar(255);not null" json:"nama_snapshot"`
	UnitPriceSnapshot decimal.Decimal `gorm:"type:numeric(18,0);not null" json:"unit_price_snapshot"`
	SubtotalSnapshot  decimal.Decimal `gorm:"type:numeric(18,0);not null" json:"subtotal_snapshot"`
	CreatedAt         time.Time       `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (AuctionBatchItem) TableName() string {
	return "auction_batch_items"
}

// AuctionBatchBrand adalah pivot batch <-> merek.
type AuctionBatchBrand struct {
	BatchID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"batch_id"`
	MerekID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"merek_id"`
}

func (AuctionBatchBrand) TableName() string {
	return "auction_batch_brands"
}

// AuctionAsset adalah aset (gambar/PDF) yang diupload admin untuk batch.
type AuctionAsset struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UploadedBy   uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	Kind         string    `gorm:"type:varchar(10);not null" json:"kind"`
	StorageKey   string    `gorm:"type:text;not null" json:"storage_key"`
	OriginalName string    `gorm:"type:text;not null" json:"original_name"`
	MimeType     string    `gorm:"type:varchar(100);not null" json:"mime_type"`
	SizeBytes    int64     `gorm:"not null" json:"size_bytes"`
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (AuctionAsset) TableName() string {
	return "auction_assets"
}

// AuctionBatchAsset adalah pivot batch <-> asset dengan urutan galeri.
type AuctionBatchAsset struct {
	BatchID   uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"batch_id"`
	AssetID   uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"asset_id"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
}

func (AuctionBatchAsset) TableName() string {
	return "auction_batch_assets"
}

// AuctionBid adalah bid buyer. Immutable; ditulis oleh storefront backend.
type AuctionBid struct {
	ID                 uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID            uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_bid_batch_seq" json:"batch_id"`
	BuyerID            uuid.UUID        `gorm:"type:uuid;not null" json:"buyer_id"`
	Sequence           int              `gorm:"not null;uniqueIndex:idx_bid_batch_seq" json:"sequence"`
	InputMode          string           `gorm:"type:varchar(10);not null" json:"input_mode"`
	InputPercent       *decimal.Decimal `gorm:"type:numeric(12,4)" json:"input_percent"`
	Amount             decimal.Decimal  `gorm:"type:numeric(18,0);not null" json:"amount"`
	GrandTotalSnapshot decimal.Decimal  `gorm:"type:numeric(18,0);not null" json:"grand_total_snapshot"`
	ShippingQuoteID    *uuid.UUID       `gorm:"type:uuid" json:"shipping_quote_id"`
	CreatedAt          time.Time        `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`

	// Relations
	Buyer *Buyer `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

func (AuctionBid) TableName() string {
	return "auction_bids"
}

// AuctionWinner adalah pemenang batch. Satu per batch; deal = nominal bid.
type AuctionWinner struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID           uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"batch_id"`
	BidID             uuid.UUID       `gorm:"type:uuid;not null" json:"bid_id"`
	DealAmount        decimal.Decimal `gorm:"type:numeric(18,0);not null" json:"deal_amount"`
	SelectedBy        uuid.UUID       `gorm:"type:uuid;not null" json:"selected_by"`
	SelectedAt        time.Time       `gorm:"type:timestamptz;not null;default:now()" json:"selected_at"`
	Note              *string         `gorm:"type:text" json:"note"`
	PaymentStatus     string          `gorm:"type:varchar(20);not null;default:UNPAID" json:"payment_status"`
	PaidAt            *time.Time      `gorm:"type:timestamptz" json:"paid_at"`
	PaymentNote       *string         `gorm:"type:text" json:"payment_note"`
	FulfillmentStatus string          `gorm:"type:varchar(20);not null;default:PENDING" json:"fulfillment_status"`
	CompletedAt       *time.Time      `gorm:"type:timestamptz" json:"completed_at"`
	FulfillmentNote   *string         `gorm:"type:text" json:"fulfillment_note"`
}

func (AuctionWinner) TableName() string {
	return "auction_winners"
}

// AuctionStockReservation adalah reservasi stok produk saat publish.
type AuctionStockReservation struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID    uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_reserv_batch_produk" json:"batch_id"`
	ProdukID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_reserv_batch_produk" json:"produk_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	Status     string     `gorm:"type:varchar(10);not null;default:ACTIVE" json:"status"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	ConsumedAt *time.Time `gorm:"type:timestamptz" json:"consumed_at"`
}

func (AuctionStockReservation) TableName() string {
	return "auction_stock_reservations"
}

// AuctionShippingQuote adalah quote ongkir yang dihasilkan storefront.
type AuctionShippingQuote struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID             uuid.UUID       `gorm:"type:uuid;not null" json:"batch_id"`
	BuyerID             uuid.UUID       `gorm:"type:uuid;not null" json:"buyer_id"`
	Provider            string          `gorm:"type:varchar(50);not null" json:"provider"`
	Service             string          `gorm:"type:varchar(100);not null" json:"service"`
	OriginSnapshot      JSONMap         `gorm:"type:jsonb;not null" json:"origin_snapshot"`
	DestinationSnapshot JSONMap         `gorm:"type:jsonb;not null" json:"destination_snapshot"`
	PhysicalSnapshot    JSONMap         `gorm:"type:jsonb;not null" json:"physical_snapshot"`
	Amount              decimal.Decimal `gorm:"type:numeric(18,0);not null" json:"amount"`
	Currency            string          `gorm:"type:varchar(10);not null;default:IDR" json:"currency"`
	QuotedAt            time.Time       `gorm:"type:timestamptz;not null;default:now()" json:"quoted_at"`
	ExpiresAt           *time.Time      `gorm:"type:timestamptz" json:"expires_at"`
}

func (AuctionShippingQuote) TableName() string {
	return "auction_shipping_quotes"
}

// AuctionEvent adalah event view batch yang ditulis storefront.
type AuctionEvent struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	BatchID    uuid.UUID  `gorm:"type:uuid;not null" json:"batch_id"`
	EventType  string     `gorm:"type:varchar(30);not null" json:"event_type"`
	BuyerID    *uuid.UUID `gorm:"type:uuid" json:"buyer_id"`
	OccurredAt time.Time  `gorm:"type:timestamptz;not null" json:"occurred_at"`
	ReceivedAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"received_at"`
}

func (AuctionEvent) TableName() string {
	return "auction_events"
}

// AuctionAuditLog adalah audit append-only aksi admin.
type AuctionAuditLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BatchID      uuid.UUID `gorm:"type:uuid;not null" json:"batch_id"`
	ActorAdminID uuid.UUID `gorm:"type:uuid;not null" json:"actor_admin_id"`
	Action       string    `gorm:"type:varchar(30);not null" json:"action"`
	BeforeData   JSONMap   `gorm:"type:jsonb" json:"before_data"`
	AfterData    JSONMap   `gorm:"type:jsonb;not null" json:"after_data"`
	Note         *string   `gorm:"type:text" json:"note"`
	CreatedAt    time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (AuctionAuditLog) TableName() string {
	return "auction_audit_logs"
}

// AuctionIdempotency menyimpan hasil mutasi berdasarkan key untuk dedup.
type AuctionIdempotency struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ActorID        uuid.UUID `gorm:"type:uuid;not null" json:"actor_id"`
	ActorType      string    `gorm:"type:varchar(10);not null" json:"actor_type"`
	Operation      string    `gorm:"type:varchar(50);not null" json:"operation"`
	Key            string    `gorm:"type:varchar(128);not null" json:"key"`
	RequestHash    string    `gorm:"type:varchar(64);not null" json:"request_hash"`
	ResponseStatus int       `gorm:"not null" json:"response_status"`
	ResponseBody   JSONMap   `gorm:"type:jsonb;not null" json:"response_body"`
	CreatedAt      time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
}

func (AuctionIdempotency) TableName() string {
	return "auction_idempotency"
}
