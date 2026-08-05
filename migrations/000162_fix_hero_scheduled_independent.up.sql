-- =====================================================
-- FIX: hero_section trigger - Banner bertanggal berdiri sendiri
-- Banner bertanggal (scheduled) tidak dipengaruhi oleh is_default
-- Banner is_default hanya mempengaruhi permanent defaults (tanpa tanggal)
-- =====================================================

DROP TRIGGER IF EXISTS trg_hero_section_auto_sync_insert ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync_insert();

CREATE OR REPLACE FUNCTION fn_hero_section_auto_sync_insert()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Jika banner punya tanggal: ini scheduled banner, pastikan is_default = false
        -- Tidak menyentuh banner lain sama sekali
        IF NEW.tanggal_mulai IS NOT NULL AND NEW.tanggal_selesai IS NOT NULL THEN
            NEW.is_default := false;

        -- Jika is_default = true tanpa tanggal: permanent default
        -- Hanya clear tanggal banner ini sendiri, lalu unset is_default
        -- dari permanent defaults lain (yang juga tidak punya tanggal)
        -- Banner bertanggal (scheduled) TIDAK disentuh sama sekali
        ELSIF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
            NEW.tanggal_mulai := NULL;
            NEW.tanggal_selesai := NULL;

            UPDATE hero_section
            SET is_default = false,
                updated_at = NOW()
            WHERE id != NEW.id
              AND is_default = true
              AND tanggal_mulai IS NULL
              AND tanggal_selesai IS NULL
              AND deleted_at IS NULL;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_hero_section_auto_sync_insert
    BEFORE INSERT ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync_insert();

COMMENT ON TRIGGER trg_hero_section_auto_sync_insert ON hero_section IS
'INSERT only: scheduled banners (punya tanggal) berdiri sendiri dengan is_default=false. Permanent default hanya unset permanent defaults lain, tidak menyentuh scheduled banners.';
