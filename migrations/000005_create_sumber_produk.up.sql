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

-- Table & Column Comments
COMMENT ON TABLE sumber_produk IS 'Menyimpan asal/sumber produk (Supplier Lokal, Import, Buyback, dll)';
COMMENT ON COLUMN sumber_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN sumber_produk.nama IS 'Nama sumber';
COMMENT ON COLUMN sumber_produk.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN sumber_produk.deskripsi IS 'Deskripsi sumber';
COMMENT ON COLUMN sumber_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN sumber_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN sumber_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN sumber_produk.deleted_at IS 'Soft delete timestamp';
