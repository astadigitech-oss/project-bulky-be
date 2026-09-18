-- ============================================================
-- PRD 02 — Model Database Lelang (Auction)
-- Membuat seluruh tabel auction untuk module lelang admin.
-- Repo: project-bulky-be
-- ============================================================

-- Enable UUID extension (idempotent)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ------------------------------------------------------------
-- 1. auction_batches — master batch lelang
-- ------------------------------------------------------------
CREATE TABLE auction_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    nama_id VARCHAR(255) NOT NULL,
    nama_en VARCHAR(255),
    description TEXT,
    warehouse_id UUID REFERENCES warehouse(id),
    kategori_id UUID REFERENCES kategori_produk(id),
    kondisi_id UUID REFERENCES kondisi_produk(id),
    kondisi_paket_id UUID REFERENCES kondisi_paket(id),
    sumber_id UUID REFERENCES sumber_produk(id),
    discrepancy_percentage NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (discrepancy_percentage >= 0 AND discrepancy_percentage <= 100),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT','OPEN','SOLD')),
    grand_total NUMERIC(18,0) NOT NULL DEFAULT 0 CHECK (grand_total >= 0),
    min_bid_percent NUMERIC(12,4) NOT NULL DEFAULT 0.1 CHECK (min_bid_percent >= 0),
    total_quantity INTEGER NOT NULL DEFAULT 0 CHECK (total_quantity >= 0),
    panjang_cm NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (panjang_cm >= 0),
    lebar_cm NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (lebar_cm >= 0),
    tinggi_cm NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (tinggi_cm >= 0),
    berat_kg NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (berat_kg >= 0),
    volume_m3 NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (volume_m3 >= 0),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    created_by UUID NOT NULL REFERENCES admin(id),
    updated_by UUID NOT NULL REFERENCES admin(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    opened_at TIMESTAMPTZ,
    sold_at TIMESTAMPTZ,
    CHECK (status IN ('OPEN','SOLD') OR opened_at IS NULL),
    CHECK (status <> 'SOLD' OR sold_at IS NOT NULL)
);

CREATE INDEX idx_auction_batches_status_created ON auction_batches(status, created_at, id);
CREATE INDEX idx_auction_batches_created ON auction_batches(created_at DESC, id);
CREATE INDEX idx_auction_batches_nama ON auction_batches(nama_id);
CREATE INDEX idx_auction_batches_warehouse ON auction_batches(warehouse_id);
CREATE INDEX idx_auction_batches_opened ON auction_batches(opened_at);

CREATE TRIGGER update_auction_batches_updated_at
    BEFORE UPDATE ON auction_batches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ------------------------------------------------------------
-- 2. auction_batch_items — item batch + snapshot harga/nama
-- ------------------------------------------------------------
CREATE TABLE auction_batch_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    produk_id UUID NOT NULL REFERENCES produk(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    nama_snapshot VARCHAR(255) NOT NULL,
    unit_price_snapshot NUMERIC(18,0) NOT NULL CHECK (unit_price_snapshot >= 0),
    subtotal_snapshot NUMERIC(18,0) NOT NULL CHECK (subtotal_snapshot >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (batch_id, produk_id)
);

CREATE INDEX idx_auction_batch_items_batch ON auction_batch_items(batch_id);
CREATE INDEX idx_auction_batch_items_produk ON auction_batch_items(produk_id);

-- ------------------------------------------------------------
-- 3. auction_batch_brands — pivot batch <-> merek
-- ------------------------------------------------------------
CREATE TABLE auction_batch_brands (
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    merek_id UUID NOT NULL REFERENCES merek_produk(id) ON DELETE RESTRICT,
    PRIMARY KEY (batch_id, merek_id)
);

CREATE INDEX idx_auction_batch_brands_merek ON auction_batch_brands(merek_id);

-- ------------------------------------------------------------
-- 4. auction_assets — aset yang diupload admin (gambar/PDF)
-- ------------------------------------------------------------
CREATE TABLE auction_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    uploaded_by UUID NOT NULL REFERENCES admin(id) ON DELETE RESTRICT,
    kind VARCHAR(10) NOT NULL CHECK (kind IN ('IMAGE','PDF')),
    storage_key TEXT NOT NULL,
    original_name TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auction_assets_uploaded_by ON auction_assets(uploaded_by);
CREATE INDEX idx_auction_assets_kind ON auction_assets(kind);

-- ------------------------------------------------------------
-- 5. auction_batch_assets — pivot batch <-> asset (urutan galeri)
-- ------------------------------------------------------------
CREATE TABLE auction_batch_assets (
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    asset_id UUID NOT NULL REFERENCES auction_assets(id) ON DELETE RESTRICT,
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    PRIMARY KEY (batch_id, asset_id)
);

CREATE INDEX idx_auction_batch_assets_asset ON auction_batch_assets(asset_id);

-- ------------------------------------------------------------
-- 6. auction_shipping_quotes — quote ongkir (ditulis storefront)
-- ------------------------------------------------------------
CREATE TABLE auction_shipping_quotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE RESTRICT,
    provider VARCHAR(50) NOT NULL,
    service VARCHAR(100) NOT NULL,
    origin_snapshot JSONB NOT NULL,
    destination_snapshot JSONB NOT NULL,
    physical_snapshot JSONB NOT NULL,
    amount NUMERIC(18,0) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    quoted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_auction_shipping_quotes_batch ON auction_shipping_quotes(batch_id);
CREATE INDEX idx_auction_shipping_quotes_buyer ON auction_shipping_quotes(buyer_id);

-- ------------------------------------------------------------
-- 7. auction_bids — bid buyer (ditulis storefront; immutable)
-- ------------------------------------------------------------
CREATE TABLE auction_bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE RESTRICT,
    sequence INTEGER NOT NULL CHECK (sequence >= 1),
    input_mode VARCHAR(10) NOT NULL CHECK (input_mode IN ('AMOUNT','PERCENT')),
    input_percent NUMERIC(12,4) CHECK (input_percent IS NULL OR input_percent >= 0),
    amount NUMERIC(18,0) NOT NULL CHECK (amount >= 0),
    grand_total_snapshot NUMERIC(18,0) NOT NULL CHECK (grand_total_snapshot >= 0),
    shipping_quote_id UUID REFERENCES auction_shipping_quotes(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (batch_id, buyer_id, sequence),
    UNIQUE (batch_id, id),
    CHECK (input_mode = 'AMOUNT' OR input_percent IS NOT NULL)
);

CREATE INDEX idx_auction_bids_batch_created ON auction_bids(batch_id, created_at, id);
CREATE INDEX idx_auction_bids_batch_amount ON auction_bids(batch_id, amount, id);
CREATE INDEX idx_auction_bids_buyer_batch ON auction_bids(buyer_id, batch_id, created_at, id);

-- ------------------------------------------------------------
-- 8. auction_winners — pemenang (satu per batch)
-- ------------------------------------------------------------
CREATE TABLE auction_winners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT UNIQUE,
    bid_id UUID NOT NULL,
    deal_amount NUMERIC(18,0) NOT NULL CHECK (deal_amount >= 0),
    selected_by UUID NOT NULL REFERENCES admin(id) ON DELETE RESTRICT,
    selected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    note TEXT,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'UNPAID' CHECK (payment_status IN ('UNPAID','PAID')),
    paid_at TIMESTAMPTZ,
    payment_note TEXT,
    fulfillment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (fulfillment_status IN ('PENDING','PROCESSING','COMPLETED')),
    completed_at TIMESTAMPTZ,
    fulfillment_note TEXT,
    CHECK (payment_status = 'PAID' OR paid_at IS NULL),
    CHECK (fulfillment_status = 'COMPLETED' OR completed_at IS NULL),
    FOREIGN KEY (batch_id, bid_id) REFERENCES auction_bids(batch_id, id) ON DELETE RESTRICT
);

CREATE INDEX idx_auction_winners_bid ON auction_winners(bid_id);

-- ------------------------------------------------------------
-- 9. auction_stock_reservations — reservasi stok saat publish
-- ------------------------------------------------------------
CREATE TABLE auction_stock_reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    produk_id UUID NOT NULL REFERENCES produk(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status VARCHAR(10) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','CONSUMED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    consumed_at TIMESTAMPTZ,
    UNIQUE (batch_id, produk_id)
);

CREATE INDEX idx_auction_stock_reservations_produk_status ON auction_stock_reservations(produk_id, status);
CREATE INDEX idx_auction_stock_reservations_batch ON auction_stock_reservations(batch_id);

-- ------------------------------------------------------------
-- 10. auction_events — event view (ditulis storefront; dedup)
-- ------------------------------------------------------------
CREATE TABLE auction_events (
    id UUID PRIMARY KEY,
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    event_type VARCHAR(30) NOT NULL CHECK (event_type IN ('BATCH_VIEW')),
    buyer_id UUID REFERENCES buyer(id) ON DELETE RESTRICT,
    occurred_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auction_events_batch_type_occurred ON auction_events(batch_id, event_type, occurred_at);

-- ------------------------------------------------------------
-- 11. auction_audit_logs — audit append-only admin
-- ------------------------------------------------------------
CREATE TABLE auction_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES auction_batches(id) ON DELETE RESTRICT,
    actor_admin_id UUID NOT NULL REFERENCES admin(id) ON DELETE RESTRICT,
    action VARCHAR(30) NOT NULL,
    before_data JSONB,
    after_data JSONB NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auction_audit_logs_batch ON auction_audit_logs(batch_id, created_at, id);

-- ------------------------------------------------------------
-- 12. auction_idempotency — idempotency key untuk mutasi
-- ------------------------------------------------------------
CREATE TABLE auction_idempotency (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id UUID NOT NULL,
    actor_type VARCHAR(10) NOT NULL CHECK (actor_type IN ('ADMIN','BUYER')),
    operation VARCHAR(50) NOT NULL,
    key VARCHAR(128) NOT NULL,
    request_hash VARCHAR(64) NOT NULL,
    response_status INTEGER NOT NULL,
    response_body JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (actor_type, actor_id, operation, key)
);

CREATE INDEX idx_auction_idempotency_actor ON auction_idempotency(actor_type, actor_id, operation);
