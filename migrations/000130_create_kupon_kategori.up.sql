CREATE TABLE kupon_kategori (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kupon_id UUID NOT NULL REFERENCES kupon(id) ON DELETE CASCADE,
    kategori_id UUID NOT NULL REFERENCES kategori_produk(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT kupon_kategori_unique UNIQUE (kupon_id, kategori_id)
);

CREATE INDEX idx_kupon_kategori_kupon_id ON kupon_kategori(kupon_id);
CREATE INDEX idx_kupon_kategori_kategori_id ON kupon_kategori(kategori_id);

COMMENT ON TABLE kupon_kategori IS 'Pivot table untuk pembatasan kategori produk pada kupon';
COMMENT ON COLUMN kupon_kategori.kupon_id IS 'Foreign key ke tabel kupon';
COMMENT ON COLUMN kupon_kategori.kategori_id IS 'Foreign key ke tabel kategori_produk';
