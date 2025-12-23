CREATE TABLE produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(255) NOT NULL,
    slug VARCHAR(280) NOT NULL UNIQUE,
    id_cargo VARCHAR(50) UNIQUE,
    kategori_id UUID NOT NULL REFERENCES kategori_produk(id),
    merek_id UUID REFERENCES merek_produk(id),
    kondisi_id UUID NOT NULL REFERENCES kondisi_produk(id),
    kondisi_paket_id UUID NOT NULL REFERENCES kondisi_paket(id),
    sumber_id UUID REFERENCES sumber_produk(id),
    warehouse_id UUID NOT NULL REFERENCES warehouse(id),
    tipe_produk_id UUID NOT NULL REFERENCES tipe_produk(id),
    harga_sebelum_diskon DECIMAL(15,2) NOT NULL CHECK (harga_sebelum_diskon > 0),
    persentase_diskon DECIMAL(5,2) DEFAULT 0 CHECK (persentase_diskon >= 0 AND persentase_diskon <= 100),
    harga_sesudah_diskon DECIMAL(15,2) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_terjual INTEGER DEFAULT 0 CHECK (quantity_terjual >= 0),
    discrepancy TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Basic indexes
CREATE INDEX idx_produk_slug ON produk(slug);
CREATE INDEX idx_produk_id_cargo ON produk(id_cargo) WHERE id_cargo IS NOT NULL;

-- Foreign key indexes
CREATE INDEX idx_produk_kategori_id ON produk(kategori_id);
CREATE INDEX idx_produk_merek_id ON produk(merek_id);
CREATE INDEX idx_produk_kondisi_id ON produk(kondisi_id);
CREATE INDEX idx_produk_kondisi_paket_id ON produk(kondisi_paket_id);
CREATE INDEX idx_produk_sumber_id ON produk(sumber_id);
CREATE INDEX idx_produk_warehouse_id ON produk(warehouse_id);
CREATE INDEX idx_produk_tipe_produk_id ON produk(tipe_produk_id);

-- Filter indexes
CREATE INDEX idx_produk_is_active ON produk(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_produk_harga ON produk(harga_sesudah_diskon);
CREATE INDEX idx_produk_created_at ON produk(created_at DESC);

-- Composite indexes for common queries
CREATE INDEX idx_produk_kategori_active ON produk(kategori_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_produk_tipe_active ON produk(tipe_produk_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_produk_warehouse_active ON produk(warehouse_id, is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER update_produk_updated_at
    BEFORE UPDATE ON produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
