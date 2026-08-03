-- migrations/000165_alter_force_update_app_dual_bahasa_platform.up.sql

-- =====================================================
-- ALTER TABLE: force_update_app
-- =====================================================
-- 1. Tambah dukungan dual bahasa (informasi_update_en)
-- 2. Tambah kolom platform (ALL/ANDROID/IOS)
-- 3. Ubah trigger single-active agar dibatasi PER PLATFORM
--    (Android dan iOS masing-masing bisa punya 1 record aktif sendiri)
-- =====================================================

CREATE TYPE force_update_platform AS ENUM ('ALL', 'ANDROID', 'IOS');

ALTER TABLE force_update_app
ADD COLUMN informasi_update_en TEXT NOT NULL DEFAULT '',
ADD COLUMN platform force_update_platform NOT NULL DEFAULT 'ALL';

ALTER TABLE force_update_app
ALTER COLUMN informasi_update_en DROP DEFAULT,
ALTER COLUMN platform DROP DEFAULT;

CREATE INDEX idx_force_update_platform_is_active ON force_update_app(platform, is_active) WHERE deleted_at IS NULL;

-- Trigger: Ensure single active PER PLATFORM
CREATE OR REPLACE FUNCTION fn_ensure_single_active_force_update()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE force_update_app 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL AND platform = NEW.platform;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON COLUMN force_update_app.informasi_update_en IS 'Changelog / release notes dalam format HTML (Bahasa Inggris)';
COMMENT ON COLUMN force_update_app.platform IS 'Platform target: ALL (semua platform), ANDROID, atau IOS. Single active dibatasi per platform.';
