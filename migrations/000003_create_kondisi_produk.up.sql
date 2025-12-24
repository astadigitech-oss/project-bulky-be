CREATE TABLE kondisi_produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT,
    urutan INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kondisi_produk_slug ON kondisi_produk(slug);
CREATE INDEX idx_kondisi_produk_urutan ON kondisi_produk(urutan);

-- Table & Column Comments
COMMENT ON TABLE kondisi_produk IS 'Menyimpan kondisi fisik produk (Baru, Bekas Seperti Baru, Bekas Layak, dll)';
COMMENT ON COLUMN kondisi_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN kondisi_produk.nama IS 'Nama kondisi';
COMMENT ON COLUMN kondisi_produk.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN kondisi_produk.deskripsi IS 'Deskripsi kondisi';
COMMENT ON COLUMN kondisi_produk.urutan IS 'Urutan tampilan';
COMMENT ON COLUMN kondisi_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN kondisi_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kondisi_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kondisi_produk.deleted_at IS 'Soft delete timestamp';
