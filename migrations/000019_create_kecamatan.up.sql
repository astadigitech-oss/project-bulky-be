CREATE TABLE kecamatan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kota_id UUID NOT NULL REFERENCES kota(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_kecamatan_kota_id ON kecamatan(kota_id);
CREATE INDEX idx_kecamatan_nama ON kecamatan(nama);
CREATE INDEX idx_kecamatan_kode ON kecamatan(kode) WHERE kode IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_kecamatan_updated_at
    BEFORE UPDATE ON kecamatan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE kecamatan IS 'Master data kecamatan di Indonesia. Berelasi ke kota/kabupaten.';
COMMENT ON COLUMN kecamatan.id IS 'Primary key UUID';
COMMENT ON COLUMN kecamatan.kota_id IS 'Foreign key ke tabel kota';
COMMENT ON COLUMN kecamatan.nama IS 'Nama kecamatan';
COMMENT ON COLUMN kecamatan.kode IS 'Kode BPS kecamatan (6 digit), nullable untuk auto-generate';
COMMENT ON COLUMN kecamatan.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kecamatan.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kecamatan.deleted_at IS 'Soft delete timestamp';
