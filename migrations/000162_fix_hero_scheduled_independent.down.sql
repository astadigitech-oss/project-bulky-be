-- Rollback: restore trigger ke versi 000118
DROP TRIGGER IF EXISTS trg_hero_section_auto_sync_insert ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync_insert();

CREATE OR REPLACE FUNCTION fn_hero_section_auto_sync_insert()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.tanggal_mulai IS NOT NULL AND NEW.tanggal_selesai IS NOT NULL THEN
        NEW.is_default := true;
        UPDATE hero_section
        SET is_default = false,
            tanggal_mulai = NULL,
            tanggal_selesai = NULL,
            updated_at = NOW()
        WHERE id != NEW.id
          AND deleted_at IS NULL;
    ELSIF TG_OP = 'INSERT' AND NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        NEW.tanggal_mulai := NULL;
        NEW.tanggal_selesai := NULL;
        UPDATE hero_section
        SET is_default = false,
            tanggal_mulai = NULL,
            tanggal_selesai = NULL,
            updated_at = NOW()
        WHERE id != NEW.id
          AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_hero_section_auto_sync_insert
    BEFORE INSERT ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync_insert();
