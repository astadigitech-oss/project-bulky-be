CREATE TABLE tipe_produk (
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

CREATE INDEX idx_tipe_produk_slug ON tipe_produk(slug);
CREATE INDEX idx_tipe_produk_urutan ON tipe_produk(urutan);

CREATE TRIGGER update_tipe_produk_updated_at
    BEFORE UPDATE ON tipe_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default data
INSERT INTO tipe_produk (nama, slug, deskripsi, urutan) VALUES
('Paletbox', 'paletbox', 'Produk dalam kemasan paletbox', 1),
('Container', 'container', 'Produk dalam kemasan container', 2),
('Truckload', 'truckload', 'Produk dalam kemasan truckload', 3);


-- Table & Column Comments
COMMENT ON TABLE tipe_produk IS 'Menyimpan tipe produk (Paletbox, Container, Truckload)';
COMMENT ON COLUMN tipe_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN tipe_produk.nama IS 'Nama tipe produk';
COMMENT ON COLUMN tipe_produk.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN tipe_produk.deskripsi IS 'Deskripsi tipe produk';
COMMENT ON COLUMN tipe_produk.urutan IS 'Urutan tampilan';
COMMENT ON COLUMN tipe_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN tipe_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN tipe_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN tipe_produk.deleted_at IS 'Soft delete timestamp';
