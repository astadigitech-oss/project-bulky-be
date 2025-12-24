CREATE TABLE diskon_kategori (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori_id UUID NOT NULL REFERENCES kategori_produk(id) ON DELETE CASCADE,
    persentase_diskon DECIMAL(5,2) NOT NULL CHECK (persentase_diskon >= 0 AND persentase_diskon <= 100),
    nominal_diskon DECIMAL(15,2) DEFAULT 0,
    tanggal_mulai DATE,
    tanggal_selesai DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    CONSTRAINT check_tanggal CHECK (tanggal_selesai IS NULL OR tanggal_mulai IS NULL OR tanggal_selesai >= tanggal_mulai)
);

CREATE INDEX idx_diskon_kategori_kategori_id ON diskon_kategori(kategori_id);
CREATE INDEX idx_diskon_kategori_tanggal ON diskon_kategori(tanggal_mulai, tanggal_selesai);
CREATE INDEX idx_diskon_kategori_is_active ON diskon_kategori(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER update_diskon_kategori_updated_at
    BEFORE UPDATE ON diskon_kategori
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- Table & Column Comments
COMMENT ON TABLE diskon_kategori IS 'Menyimpan diskon otomatis per kategori produk';
COMMENT ON COLUMN diskon_kategori.id IS 'Primary key UUID';
COMMENT ON COLUMN diskon_kategori.kategori_id IS 'FK ke kategori_produk';
COMMENT ON COLUMN diskon_kategori.persentase_diskon IS 'Persentase diskon (0-100)';
COMMENT ON COLUMN diskon_kategori.nominal_diskon IS 'Nominal diskon fixed (opsional)';
COMMENT ON COLUMN diskon_kategori.tanggal_mulai IS 'Tanggal mulai berlaku';
COMMENT ON COLUMN diskon_kategori.tanggal_selesai IS 'Tanggal selesai berlaku';
COMMENT ON COLUMN diskon_kategori.is_active IS 'Status aktif';
COMMENT ON COLUMN diskon_kategori.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN diskon_kategori.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN diskon_kategori.deleted_at IS 'Soft delete timestamp';
