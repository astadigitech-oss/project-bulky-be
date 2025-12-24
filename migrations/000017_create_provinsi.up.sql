CREATE TABLE provinsi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_provinsi_nama ON provinsi(nama);
CREATE INDEX idx_provinsi_kode ON provinsi(kode) WHERE kode IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_provinsi_updated_at
    BEFORE UPDATE ON provinsi
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE provinsi IS 'Master data provinsi di Indonesia. Kode mengikuti standar BPS.';
COMMENT ON COLUMN provinsi.id IS 'Primary key UUID';
COMMENT ON COLUMN provinsi.nama IS 'Nama provinsi';
COMMENT ON COLUMN provinsi.kode IS 'Kode BPS provinsi (2 digit), nullable untuk auto-generate';
COMMENT ON COLUMN provinsi.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN provinsi.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN provinsi.deleted_at IS 'Soft delete timestamp';
