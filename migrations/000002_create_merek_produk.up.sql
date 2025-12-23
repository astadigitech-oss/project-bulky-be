CREATE TABLE merek_produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    logo_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_merek_produk_slug ON merek_produk(slug);
CREATE INDEX idx_merek_produk_is_active ON merek_produk(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_merek_produk_nama ON merek_produk(nama);
