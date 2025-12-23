CREATE TABLE sumber_produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_sumber_produk_slug ON sumber_produk(slug);
CREATE INDEX idx_sumber_produk_is_active ON sumber_produk(is_active) WHERE deleted_at IS NULL;
