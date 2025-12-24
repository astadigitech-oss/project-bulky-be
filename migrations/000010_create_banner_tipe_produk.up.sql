CREATE TABLE banner_tipe_produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tipe_produk_id UUID NOT NULL REFERENCES tipe_produk(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    gambar_url VARCHAR(500) NOT NULL,
    urutan INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_banner_tipe_produk_tipe_id ON banner_tipe_produk(tipe_produk_id);
CREATE INDEX idx_banner_tipe_produk_urutan ON banner_tipe_produk(urutan);
CREATE INDEX idx_banner_tipe_produk_is_active ON banner_tipe_produk(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER update_banner_tipe_produk_updated_at
    BEFORE UPDATE ON banner_tipe_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- Table & Column Comments
COMMENT ON TABLE banner_tipe_produk IS 'Menyimpan banner untuk setiap tipe produk';
COMMENT ON COLUMN banner_tipe_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN banner_tipe_produk.tipe_produk_id IS 'FK ke tipe_produk';
COMMENT ON COLUMN banner_tipe_produk.nama IS 'Nama/judul banner';
COMMENT ON COLUMN banner_tipe_produk.gambar_url IS 'URL gambar banner';
COMMENT ON COLUMN banner_tipe_produk.urutan IS 'Urutan tampilan';
COMMENT ON COLUMN banner_tipe_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN banner_tipe_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN banner_tipe_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN banner_tipe_produk.deleted_at IS 'Soft delete timestamp';
