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


-- Table & Column Comments
COMMENT ON TABLE produk_dokumen IS 'Menyimpan dokumen/file PDF per produk';
COMMENT ON COLUMN produk_dokumen.id IS 'Primary key UUID';
COMMENT ON COLUMN produk_dokumen.produk_id IS 'FK ke produk';
COMMENT ON COLUMN produk_dokumen.nama_dokumen IS 'Nama file/dokumen';
COMMENT ON COLUMN produk_dokumen.file_url IS 'URL file';
COMMENT ON COLUMN produk_dokumen.tipe_file IS 'Tipe file (pdf, xlsx, dll)';
COMMENT ON COLUMN produk_dokumen.ukuran_file IS 'Ukuran file dalam bytes';
COMMENT ON COLUMN produk_dokumen.created_at IS 'Waktu dibuat';
