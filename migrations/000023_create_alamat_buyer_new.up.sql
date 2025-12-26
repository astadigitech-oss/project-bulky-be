-- =====================================================
-- TABEL: alamat_buyer (Struktur Baru)
-- =====================================================
-- Menyimpan alamat pengiriman buyer
-- Data wilayah di-fetch dari Google Maps API oleh frontend
-- 
-- Business Rules:
-- 1. Satu buyer bisa punya banyak alamat
-- 2. Hanya boleh ada 1 alamat default per buyer
-- 3. Alamat pertama otomatis jadi default
-- 4. Set default baru akan unset yang lama
-- 5. Tidak bisa hapus default jika masih ada alamat lain
--
-- Lihat: 12-REVISI-ALAMAT-BUYER.md
-- =====================================================

CREATE TABLE alamat_buyer (
    -- Primary Key
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign Key ke Buyer
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE CASCADE,
    
    -- Identitas Penerima
    label VARCHAR(50) NOT NULL,
    nama_penerima VARCHAR(100) NOT NULL,
    telepon_penerima VARCHAR(20) NOT NULL,
    
    -- Wilayah (TEXT dari Google Maps API)
    provinsi VARCHAR(100) NOT NULL,
    kota VARCHAR(100) NOT NULL,
    kecamatan VARCHAR(100),
    kelurahan VARCHAR(100),
    kode_pos VARCHAR(10),
    
    -- Alamat Detail
    alamat_lengkap TEXT NOT NULL,
    catatan TEXT,
    
    -- Koordinat (dari Google Maps API)
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Google Maps Reference
    google_place_id VARCHAR(255),
    
    -- Status Default
    is_default BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Query by buyer
CREATE INDEX idx_alamat_buyer_buyer_id 
    ON alamat_buyer(buyer_id);

-- Query default address
CREATE INDEX idx_alamat_buyer_default 
    ON alamat_buyer(buyer_id, is_default) 
    WHERE deleted_at IS NULL;

-- Filter by wilayah (untuk reporting)
CREATE INDEX idx_alamat_buyer_provinsi 
    ON alamat_buyer(provinsi);

CREATE INDEX idx_alamat_buyer_kota 
    ON alamat_buyer(kota);

-- Geo queries (kalkulasi jarak)
CREATE INDEX idx_alamat_buyer_coordinates 
    ON alamat_buyer(latitude, longitude) 
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- =====================================================
-- TRIGGER 1: Auto Update Timestamp
-- =====================================================

CREATE TRIGGER trg_alamat_buyer_updated_at
    BEFORE UPDATE ON alamat_buyer
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- TRIGGER 2: Ensure Single Default Per Buyer
-- =====================================================
-- Ketika set alamat sebagai default, unset yang lain

CREATE OR REPLACE FUNCTION fn_ensure_single_default_alamat()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        -- Unset semua default lain milik buyer yang sama
        UPDATE alamat_buyer 
        SET is_default = false,
            updated_at = NOW()
        WHERE buyer_id = NEW.buyer_id 
          AND id != NEW.id 
          AND is_default = true
          AND deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_default_alamat
    AFTER INSERT OR UPDATE OF is_default ON alamat_buyer
    FOR EACH ROW
    WHEN (NEW.is_default = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_default_alamat();

-- =====================================================
-- TRIGGER 3: First Address Auto Default
-- =====================================================
-- Alamat pertama otomatis jadi default

CREATE OR REPLACE FUNCTION fn_first_alamat_as_default()
RETURNS TRIGGER AS $$
DECLARE
    v_count INTEGER;
BEGIN
    -- Hitung alamat aktif milik buyer
    SELECT COUNT(*) INTO v_count 
    FROM alamat_buyer 
    WHERE buyer_id = NEW.buyer_id 
      AND deleted_at IS NULL;
    
    -- Jika ini alamat pertama, set sebagai default
    IF v_count = 0 THEN
        NEW.is_default := true;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_first_alamat_default
    BEFORE INSERT ON alamat_buyer
    FOR EACH ROW
    EXECUTE FUNCTION fn_first_alamat_as_default();

-- =====================================================
-- TRIGGER 4: Prevent Delete Default (if has others)
-- =====================================================
-- Tidak boleh hapus alamat default jika masih ada alamat lain

CREATE OR REPLACE FUNCTION fn_prevent_delete_default_alamat()
RETURNS TRIGGER AS $$
DECLARE
    v_other_count INTEGER;
BEGIN
    -- Cek apakah ini soft delete pada alamat default
    IF OLD.is_default = true 
       AND OLD.deleted_at IS NULL 
       AND NEW.deleted_at IS NOT NULL THEN
        
        -- Hitung alamat lain yang masih aktif
        SELECT COUNT(*) INTO v_other_count 
        FROM alamat_buyer 
        WHERE buyer_id = OLD.buyer_id 
          AND id != OLD.id
          AND deleted_at IS NULL;
        
        -- Jika masih ada alamat lain, tolak
        IF v_other_count > 0 THEN
            RAISE EXCEPTION 'Tidak dapat menghapus alamat default. Set alamat lain sebagai default terlebih dahulu.';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_delete_default
    BEFORE UPDATE ON alamat_buyer
    FOR EACH ROW
    EXECUTE FUNCTION fn_prevent_delete_default_alamat();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE alamat_buyer IS 
'Alamat pengiriman buyer. Data wilayah dari Google Maps API. 1 buyer bisa punya banyak alamat, tapi hanya 1 default.';

-- Kolom identitas
COMMENT ON COLUMN alamat_buyer.id IS 'Primary key UUID';
COMMENT ON COLUMN alamat_buyer.buyer_id IS 'FK ke tabel buyer';
COMMENT ON COLUMN alamat_buyer.label IS 'Label alamat: Rumah, Kantor, Toko, dll';
COMMENT ON COLUMN alamat_buyer.nama_penerima IS 'Nama penerima paket';
COMMENT ON COLUMN alamat_buyer.telepon_penerima IS 'Telepon penerima untuk kurir';

-- Kolom wilayah
COMMENT ON COLUMN alamat_buyer.provinsi IS 'Provinsi dari Google Maps (administrative_area_level_1)';
COMMENT ON COLUMN alamat_buyer.kota IS 'Kota/Kabupaten dari Google Maps (administrative_area_level_2)';
COMMENT ON COLUMN alamat_buyer.kecamatan IS 'Kecamatan dari Google Maps (administrative_area_level_3)';
COMMENT ON COLUMN alamat_buyer.kelurahan IS 'Kelurahan dari Google Maps (administrative_area_level_4/sublocality)';
COMMENT ON COLUMN alamat_buyer.kode_pos IS 'Kode pos dari Google Maps (postal_code)';

-- Kolom alamat detail
COMMENT ON COLUMN alamat_buyer.alamat_lengkap IS 'Alamat lengkap: jalan, nomor, RT/RW';
COMMENT ON COLUMN alamat_buyer.catatan IS 'Catatan untuk kurir: patokan, warna rumah, dll';

-- Kolom koordinat
COMMENT ON COLUMN alamat_buyer.latitude IS 'Latitude dari Google Maps untuk map pin & kalkulasi jarak';
COMMENT ON COLUMN alamat_buyer.longitude IS 'Longitude dari Google Maps untuk map pin & kalkulasi jarak';
COMMENT ON COLUMN alamat_buyer.google_place_id IS 'Google Place ID untuk reference';

-- Kolom status
COMMENT ON COLUMN alamat_buyer.is_default IS 'Alamat default (hanya 1 per buyer, di-enforce oleh trigger)';
COMMENT ON COLUMN alamat_buyer.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN alamat_buyer.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN alamat_buyer.deleted_at IS 'Soft delete timestamp';
