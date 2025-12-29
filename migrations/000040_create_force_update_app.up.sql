-- migrations/000040_create_force_update_app.up.sql

-- =====================================================
-- TABEL: force_update_app
-- =====================================================
-- Kontrol versi aplikasi mobile
-- Memaksa user update ketika ada major version atau critical fix
-- Hanya 1 yang bisa aktif
-- =====================================================

CREATE TYPE update_type AS ENUM ('OPTIONAL', 'MANDATORY');

CREATE TABLE force_update_app (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode_versi VARCHAR(20) NOT NULL,
    update_type update_type NOT NULL,
    informasi_update TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_force_update_is_active ON force_update_app(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_force_update_kode_versi ON force_update_app(kode_versi);

CREATE TRIGGER trg_force_update_updated_at
    BEFORE UPDATE ON force_update_app
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger: Ensure single active
CREATE OR REPLACE FUNCTION fn_ensure_single_active_force_update()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE force_update_app 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_force_update
    AFTER INSERT OR UPDATE OF is_active ON force_update_app
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_force_update();

COMMENT ON TABLE force_update_app IS 'Kontrol versi aplikasi mobile. OPTIONAL = bisa skip, MANDATORY = wajib update.';
COMMENT ON COLUMN force_update_app.kode_versi IS 'Versi terbaru aplikasi (format: X.Y.Z)';
COMMENT ON COLUMN force_update_app.informasi_update IS 'Changelog / release notes dalam format HTML';
