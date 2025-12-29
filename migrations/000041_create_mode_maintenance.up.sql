-- migrations/000041_create_mode_maintenance.up.sql

-- =====================================================
-- TABEL: mode_maintenance
-- =====================================================
-- Mode maintenance global
-- Ketika aktif, SELURUH app menampilkan halaman maintenance
-- Hanya 1 yang bisa aktif
-- =====================================================

CREATE TYPE maintenance_type AS ENUM ('BUG', 'ERROR', 'BIG_UPDATE', 'OTHER');

CREATE TABLE mode_maintenance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul VARCHAR(100) NOT NULL,
    tipe_maintenance maintenance_type NOT NULL,
    deskripsi TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_maintenance_is_active ON mode_maintenance(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_maintenance_updated_at
    BEFORE UPDATE ON mode_maintenance
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger: Ensure single active
CREATE OR REPLACE FUNCTION fn_ensure_single_active_maintenance()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE mode_maintenance 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_maintenance
    AFTER INSERT OR UPDATE OF is_active ON mode_maintenance
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_maintenance();

COMMENT ON TABLE mode_maintenance IS 'Mode maintenance global. Ketika aktif, seluruh app menampilkan halaman maintenance.';
