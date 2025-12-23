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
