CREATE TABLE produk_gambar (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    produk_id UUID NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    gambar_url VARCHAR(500) NOT NULL,
    urutan INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_produk_gambar_produk_id ON produk_gambar(produk_id);
CREATE INDEX idx_produk_gambar_urutan ON produk_gambar(produk_id, urutan);

-- Function to ensure only one primary image per product
CREATE OR REPLACE FUNCTION ensure_single_primary_image()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_primary = true THEN
        UPDATE produk_gambar 
        SET is_primary = false 
        WHERE produk_id = NEW.produk_id AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_single_primary_image
    AFTER INSERT OR UPDATE ON produk_gambar
    FOR EACH ROW
    WHEN (NEW.is_primary = true)
    EXECUTE FUNCTION ensure_single_primary_image();


-- Table & Column Comments
COMMENT ON TABLE produk_gambar IS 'Menyimpan multiple gambar per produk';
COMMENT ON COLUMN produk_gambar.id IS 'Primary key UUID';
COMMENT ON COLUMN produk_gambar.produk_id IS 'FK ke produk';
COMMENT ON COLUMN produk_gambar.gambar_url IS 'URL gambar';
COMMENT ON COLUMN produk_gambar.urutan IS 'Urutan tampilan gambar';
COMMENT ON COLUMN produk_gambar.is_primary IS 'Gambar utama (hanya 1 per produk)';
COMMENT ON COLUMN produk_gambar.created_at IS 'Waktu dibuat';
