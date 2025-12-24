CREATE TABLE kelurahan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kecamatan_id UUID NOT NULL REFERENCES kecamatan(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(15) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_kelurahan_kecamatan_id ON kelurahan(kecamatan_id);
CREATE INDEX idx_kelurahan_nama ON kelurahan(nama);
CREATE INDEX idx_kelurahan_kode ON kelurahan(kode) WHERE kode IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_kelurahan_updated_at
    BEFORE UPDATE ON kelurahan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE kelurahan IS 'Master data kelurahan/desa di Indonesia. Berelasi ke kecamatan. Level terbawah dalam hierarki wilayah.';
COMMENT ON COLUMN kelurahan.id IS 'Primary key UUID';
COMMENT ON COLUMN kelurahan.kecamatan_id IS 'Foreign key ke tabel kecamatan';
COMMENT ON COLUMN kelurahan.nama IS 'Nama kelurahan/desa';
COMMENT ON COLUMN kelurahan.kode IS 'Kode BPS kelurahan (10 digit), nullable untuk auto-generate';
COMMENT ON COLUMN kelurahan.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kelurahan.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kelurahan.deleted_at IS 'Soft delete timestamp';
