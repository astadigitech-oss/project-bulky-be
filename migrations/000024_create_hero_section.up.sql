-- =====================================================
-- TABEL: hero_section
-- =====================================================
-- Banner utama di homepage
-- Hanya 1 yang boleh aktif pada satu waktu
-- Support scheduling dengan tanggal_mulai & tanggal_selesai
-- =====================================================

CREATE TABLE hero_section (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    gambar VARCHAR(255) NOT NULL,
    urutan INT DEFAULT 0,
    is_active BOOLEAN DEFAULT false,
    tanggal_mulai TIMESTAMP,
    tanggal_selesai TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_hero_section_is_active ON hero_section(is_active) 
    WHERE deleted_at IS NULL;
CREATE INDEX idx_hero_section_tanggal ON hero_section(tanggal_mulai, tanggal_selesai);
CREATE INDEX idx_hero_section_urutan ON hero_section(urutan) 
    WHERE is_active = true AND deleted_at IS NULL;

-- Trigger: Auto update timestamp
CREATE TRIGGER trg_hero_section_updated_at
    BEFORE UPDATE ON hero_section
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- TRIGGER: Ensure Single Active Hero
-- =====================================================
-- Hanya 1 hero yang boleh aktif
-- Jika set hero baru sebagai aktif, yang lama otomatis di-unset

CREATE OR REPLACE FUNCTION fn_ensure_single_active_hero()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        -- Unset semua hero aktif lainnya
        UPDATE hero_section 
        SET is_active = false,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND is_active = true
          AND deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_hero
    AFTER INSERT OR UPDATE OF is_active ON hero_section
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_hero();

-- Comments
COMMENT ON TABLE hero_section IS 'Banner utama homepage. Hanya 1 yang boleh aktif. Support scheduling.';
COMMENT ON COLUMN hero_section.id IS 'Primary key UUID';
COMMENT ON COLUMN hero_section.nama IS 'Nama/judul hero untuk identifikasi di admin panel';
COMMENT ON COLUMN hero_section.gambar IS 'Path atau URL gambar hero';
COMMENT ON COLUMN hero_section.urutan IS 'Urutan tampil (untuk future carousel support)';
COMMENT ON COLUMN hero_section.is_active IS 'Status aktif. Hanya 1 yang boleh true (enforced by trigger)';
COMMENT ON COLUMN hero_section.tanggal_mulai IS 'Waktu mulai tampil. NULL = langsung aktif saat is_active=true';
COMMENT ON COLUMN hero_section.tanggal_selesai IS 'Waktu selesai tampil. NULL = tampil selamanya';
COMMENT ON COLUMN hero_section.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN hero_section.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN hero_section.deleted_at IS 'Soft delete timestamp';
