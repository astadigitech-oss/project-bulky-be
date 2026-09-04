-- Force update is configured independently for Android and iOS.
-- Existing ALL records are retired: one shared build number cannot represent both platforms.
ALTER TABLE force_update_app
ADD COLUMN minimum_build_number BIGINT;

UPDATE force_update_app
SET is_active = false,
    deleted_at = COALESCE(deleted_at, NOW()),
    platform = 'ANDROID'
WHERE platform = 'ALL';

UPDATE force_update_app
SET is_active = false
WHERE minimum_build_number IS NULL;

ALTER TABLE force_update_app
ADD CONSTRAINT chk_force_update_platform
CHECK (platform IN ('ANDROID', 'IOS')),
ADD CONSTRAINT chk_force_update_minimum_build_number
CHECK (minimum_build_number IS NULL OR minimum_build_number > 0);

DROP TRIGGER IF EXISTS trg_single_active_force_update ON force_update_app;
CREATE TRIGGER trg_single_active_force_update
    AFTER INSERT OR UPDATE OF is_active, platform ON force_update_app
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_force_update();

COMMENT ON COLUMN force_update_app.kode_versi IS 'Versi tampilan untuk UI mobile (contoh: 2.4.1), bukan dasar keputusan update.';
COMMENT ON COLUMN force_update_app.minimum_build_number IS 'Build number minimum yang diizinkan memakai aplikasi. Dipakai untuk menentukan force update.';
