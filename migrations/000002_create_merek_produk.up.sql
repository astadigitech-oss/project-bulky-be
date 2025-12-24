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

-- Table & Column Comments
COMMENT ON TABLE merek_produk IS 'Menyimpan data merek/brand produk';
COMMENT ON COLUMN merek_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN merek_produk.nama IS 'Nama merek';
COMMENT ON COLUMN merek_produk.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN merek_produk.logo_url IS 'URL logo merek';
COMMENT ON COLUMN merek_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN merek_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN merek_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN merek_produk.deleted_at IS 'Soft delete timestamp';
