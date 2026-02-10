-- =====================================================
-- REVISI: Hero Section - Remove is_default Trigger on UPDATE
-- =====================================================
-- Menghapus trigger auto-set is_default pada UPDATE
-- is_default hanya bisa diubah via toggle endpoint
-- Keep trigger untuk INSERT (optional behavior)
-- =====================================================

-- Drop existing trigger
DROP TRIGGER IF EXISTS trg_hero_section_auto_sync ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync();

-- Create new trigger function - INSERT ONLY
CREATE OR REPLACE FUNCTION fn_hero_section_auto_sync_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- Only on INSERT: If tanggal_mulai AND tanggal_selesai are both set
    -- This hero becomes THE ONLY active hero (scheduled)
    IF TG_OP = 'INSERT' AND NEW.tanggal_mulai IS NOT NULL AND NEW.tanggal_selesai IS NOT NULL THEN
        NEW.is_default := true;
        
        -- Clear ALL other records: unset is_default AND clear their dates
        UPDATE hero_section 
        SET is_default = false,
            tanggal_mulai = NULL,
            tanggal_selesai = NULL,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND deleted_at IS NULL;
    
    -- Only on INSERT: If is_default = true (create as permanent default)
    -- This hero becomes THE ONLY active hero (permanent)
    ELSIF TG_OP = 'INSERT' AND NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        -- Clear tanggal for this record (permanent default)
        NEW.tanggal_mulai := NULL;
        NEW.tanggal_selesai := NULL;
        
        -- Clear ALL other records: unset is_default AND clear their dates
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

-- Create trigger - INSERT ONLY
CREATE TRIGGER trg_hero_section_auto_sync_insert
    BEFORE INSERT ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync_insert();

COMMENT ON TRIGGER trg_hero_section_auto_sync_insert ON hero_section IS 
'Auto-sync is_default on INSERT only. UPDATE must use toggle endpoint.';
