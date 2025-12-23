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
