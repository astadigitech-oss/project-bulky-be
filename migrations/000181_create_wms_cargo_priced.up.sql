-- migrations/000181_create_wms_cargo_priced.up.sql
-- Cache lokal cargo WMS yang sudah diberi harga jual (already-priced di WMS).
-- Jembatan antara WMS dan produk lokal: dropdown "ID Cargo" saat create/edit
-- produk baca dari tabel ini, dan menyimpan path PDF harga yang sudah
-- di-download dari WMS (pricing_pdf_url) supaya tidak perlu download ulang
-- tiap kali produk dibuat dari cargo yang sama.

CREATE TABLE wms_cargo_priced (
    cargo_id UUID PRIMARY KEY, -- id cargo di sisi WMS (bukan id produk lokal)
    code VARCHAR(100) NOT NULL,
    length_cm DECIMAL(10,2) NOT NULL DEFAULT 0,
    width_cm DECIMAL(10,2) NOT NULL DEFAULT 0,
    height_cm DECIMAL(10,2) NOT NULL DEFAULT 0,
    weight_kg DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_price DECIMAL(15,2) NOT NULL DEFAULT 0,
    pricing_type VARCHAR(20),
    pricing_value DECIMAL(15,2),
    sale_price DECIMAL(15,2),
    priced_at TIMESTAMPTZ,
    pricing_pdf_path VARCHAR(255), -- path lokal relatif terhadap UPLOAD_PATH
    bulky_category JSONB,
    bulky_product_condition JSONB,
    bulky_package_condition JSONB,
    bulky_product_source JSONB,
    bulky_brands JSONB,
    is_used_in_produk BOOLEAN NOT NULL DEFAULT false, -- sudah dipakai di produk lokal (create/update sukses)
    produk_id UUID REFERENCES produk(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wms_cargo_priced_code ON wms_cargo_priced(code);
CREATE INDEX idx_wms_cargo_priced_is_used ON wms_cargo_priced(is_used_in_produk);
CREATE INDEX idx_wms_cargo_priced_produk_id ON wms_cargo_priced(produk_id);
