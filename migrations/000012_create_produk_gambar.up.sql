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
