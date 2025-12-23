CREATE TABLE produk_dokumen (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    produk_id UUID NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    nama_dokumen VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    tipe_file VARCHAR(50) NOT NULL,
    ukuran_file INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_produk_dokumen_produk_id ON produk_dokumen(produk_id);
