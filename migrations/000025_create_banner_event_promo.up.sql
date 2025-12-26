-- =====================================================
-- TABEL: banner_event_promo
-- =====================================================
-- Banner untuk event dan promo
-- Bisa multiple aktif (carousel/slider)
-- Support scheduling dengan tanggal_mulai & tanggal_selesai
-- =====================================================

CREATE TABLE banner_event_promo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    gambar VARCHAR(255) NOT NULL,
    url_tujuan VARCHAR(255),
    urutan INT DEFAULT 0,
    is_active BOOLEAN DEFAULT false,
    tanggal_mulai TIMESTAMP,
    tanggal_selesai TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_banner_event_is_active ON banner_event_promo(is_active) 
    WHERE deleted_at IS NULL;
CREATE INDEX idx_banner_event_urutan ON banner_event_promo(urutan) 
    WHERE is_active = true AND deleted_at IS NULL;
CREATE INDEX idx_banner_event_tanggal ON banner_event_promo(tanggal_mulai, tanggal_selesai);

-- Trigger: Auto update timestamp
CREATE TRIGGER trg_banner_event_updated_at
    BEFORE UPDATE ON banner_event_promo
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE banner_event_promo IS 'Banner event dan promo. Bisa multiple aktif (carousel). Support scheduling.';
COMMENT ON COLUMN banner_event_promo.id IS 'Primary key UUID';
COMMENT ON COLUMN banner_event_promo.nama IS 'Nama banner untuk identifikasi di admin panel';
COMMENT ON COLUMN banner_event_promo.gambar IS 'Path atau URL gambar banner';
COMMENT ON COLUMN banner_event_promo.url_tujuan IS 'Link redirect ketika banner diklik (nullable)';
COMMENT ON COLUMN banner_event_promo.urutan IS 'Urutan tampil di carousel (ASC)';
COMMENT ON COLUMN banner_event_promo.is_active IS 'Status aktif. Bisa multiple true.';
COMMENT ON COLUMN banner_event_promo.tanggal_mulai IS 'Waktu mulai tampil. NULL = langsung aktif saat is_active=true';
COMMENT ON COLUMN banner_event_promo.tanggal_selesai IS 'Waktu selesai tampil. NULL = tampil selamanya';
COMMENT ON COLUMN banner_event_promo.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN banner_event_promo.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN banner_event_promo.deleted_at IS 'Soft delete timestamp';
