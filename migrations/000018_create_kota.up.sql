CREATE TABLE kota (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provinsi_id UUID NOT NULL REFERENCES provinsi(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_kota_provinsi_id ON kota(provinsi_id);
CREATE INDEX idx_kota_nama ON kota(nama);
CREATE INDEX idx_kota_kode ON kota(kode) WHERE kode IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_kota_updated_at
    BEFORE UPDATE ON kota
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE kota IS 'Master data kabupaten/kota di Indonesia. Berelasi ke provinsi.';
COMMENT ON COLUMN kota.id IS 'Primary key UUID';
COMMENT ON COLUMN kota.provinsi_id IS 'Foreign key ke tabel provinsi';
COMMENT ON COLUMN kota.nama IS 'Nama kabupaten/kota';
COMMENT ON COLUMN kota.kode IS 'Kode BPS kabupaten/kota (4 digit), nullable untuk auto-generate';
COMMENT ON COLUMN kota.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kota.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kota.deleted_at IS 'Soft delete timestamp';
