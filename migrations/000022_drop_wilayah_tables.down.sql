-- =====================================================
-- RECREATE WILAYAH TABLES (jika perlu rollback)
-- =====================================================
-- Note: Ini hanya struktur kosong, data perlu di-seed ulang
-- =====================================================

-- Provinsi
CREATE TABLE IF NOT EXISTS provinsi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_provinsi_nama ON provinsi(nama);
CREATE INDEX IF NOT EXISTS idx_provinsi_kode ON provinsi(kode) WHERE kode IS NOT NULL;

CREATE TRIGGER update_provinsi_updated_at
    BEFORE UPDATE ON provinsi
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Kota
CREATE TABLE IF NOT EXISTS kota (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provinsi_id UUID NOT NULL REFERENCES provinsi(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_kota_provinsi_id ON kota(provinsi_id);
CREATE INDEX IF NOT EXISTS idx_kota_nama ON kota(nama);
CREATE INDEX IF NOT EXISTS idx_kota_kode ON kota(kode) WHERE kode IS NOT NULL;

CREATE TRIGGER update_kota_updated_at
    BEFORE UPDATE ON kota
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Kecamatan
CREATE TABLE IF NOT EXISTS kecamatan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kota_id UUID NOT NULL REFERENCES kota(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_kecamatan_kota_id ON kecamatan(kota_id);
CREATE INDEX IF NOT EXISTS idx_kecamatan_nama ON kecamatan(nama);
CREATE INDEX IF NOT EXISTS idx_kecamatan_kode ON kecamatan(kode) WHERE kode IS NOT NULL;

CREATE TRIGGER update_kecamatan_updated_at
    BEFORE UPDATE ON kecamatan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Kelurahan
CREATE TABLE IF NOT EXISTS kelurahan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kecamatan_id UUID NOT NULL REFERENCES kecamatan(id) ON DELETE CASCADE,
    nama VARCHAR(100) NOT NULL,
    kode VARCHAR(15) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_kelurahan_kecamatan_id ON kelurahan(kecamatan_id);
CREATE INDEX IF NOT EXISTS idx_kelurahan_nama ON kelurahan(nama);
CREATE INDEX IF NOT EXISTS idx_kelurahan_kode ON kelurahan(kode) WHERE kode IS NOT NULL;

CREATE TRIGGER update_kelurahan_updated_at
    BEFORE UPDATE ON kelurahan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Alamat Buyer (struktur lama dengan FK ke kelurahan)
CREATE TABLE IF NOT EXISTS alamat_buyer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE CASCADE,
    kelurahan_id UUID NOT NULL REFERENCES kelurahan(id),
    label VARCHAR(50) NOT NULL,
    nama_penerima VARCHAR(100) NOT NULL,
    telepon_penerima VARCHAR(20) NOT NULL,
    kode_pos VARCHAR(10) NOT NULL,
    alamat_lengkap TEXT NOT NULL,
    catatan TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alamat_buyer_buyer_id ON alamat_buyer(buyer_id);
CREATE INDEX IF NOT EXISTS idx_alamat_buyer_kelurahan_id ON alamat_buyer(kelurahan_id);
