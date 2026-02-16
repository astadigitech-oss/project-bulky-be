CREATE TABLE kupon_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kupon_id UUID NOT NULL REFERENCES kupon(id),
    buyer_id UUID NOT NULL REFERENCES buyer(id),
    pesanan_id UUID NOT NULL REFERENCES pesanan(id),
    kode_kupon VARCHAR(50) NOT NULL,
    nilai_potongan DECIMAL(15,2) NOT NULL CHECK (nilai_potongan > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kupon_usage_kupon_id ON kupon_usage(kupon_id);
CREATE INDEX idx_kupon_usage_buyer_id ON kupon_usage(buyer_id);
CREATE INDEX idx_kupon_usage_pesanan_id ON kupon_usage(pesanan_id);
CREATE INDEX idx_kupon_usage_created_at ON kupon_usage(created_at DESC);

COMMENT ON TABLE kupon_usage IS 'Tracking penggunaan kupon pada setiap order';
COMMENT ON COLUMN kupon_usage.kupon_id IS 'Foreign key ke tabel kupon';
COMMENT ON COLUMN kupon_usage.buyer_id IS 'Foreign key ke tabel buyer yang menggunakan';
COMMENT ON COLUMN kupon_usage.pesanan_id IS 'Foreign key ke tabel pesanan';
COMMENT ON COLUMN kupon_usage.kode_kupon IS 'Snapshot kode kupon saat digunakan';
COMMENT ON COLUMN kupon_usage.nilai_potongan IS 'Nilai potongan yang diterapkan';
