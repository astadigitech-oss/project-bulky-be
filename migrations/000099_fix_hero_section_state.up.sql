-- Emergency fix for inconsistent state
-- This migration will check current state and fix it

-- Check and fix hero_section table
DO $$ 
DECLARE
    has_urutan BOOLEAN;
    has_is_active BOOLEAN;
    has_is_default BOOLEAN;
    has_tanggal_mulai BOOLEAN;
    has_tanggal_selesai BOOLEAN;
BEGIN
    -- Check what columns exist
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'urutan'
    ) INTO has_urutan;
    
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'is_active'
    ) INTO has_is_active;
    
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'is_default'
    ) INTO has_is_default;
    
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'tanggal_mulai'
    ) INTO has_tanggal_mulai;
    
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'tanggal_selesai'
    ) INTO has_tanggal_selesai;
    
    -- Drop urutan if exists
    IF has_urutan THEN
        ALTER TABLE hero_section DROP COLUMN urutan;
        RAISE NOTICE 'Dropped urutan column';
    END IF;
    
    -- Rename is_active to is_default if needed
    IF has_is_active AND NOT has_is_default THEN
        ALTER TABLE hero_section RENAME COLUMN is_active TO is_default;
        RAISE NOTICE 'Renamed is_active to is_default';
    END IF;
    
    -- Add tanggal columns if not exist
    IF NOT has_tanggal_mulai THEN
        ALTER TABLE hero_section ADD COLUMN tanggal_mulai TIMESTAMP NULL;
        RAISE NOTICE 'Added tanggal_mulai column';
    END IF;
    
    IF NOT has_tanggal_selesai THEN
        ALTER TABLE hero_section ADD COLUMN tanggal_selesai TIMESTAMP NULL;
        RAISE NOTICE 'Added tanggal_selesai column';
    END IF;
END $$;

-- Clean up old indexes and triggers
DROP INDEX IF EXISTS idx_hero_section_is_active;
DROP INDEX IF EXISTS idx_hero_section_urutan;
DROP TRIGGER IF EXISTS trg_single_active_hero ON hero_section;
DROP FUNCTION IF EXISTS fn_ensure_single_active_hero();

-- Create new indexes
DROP INDEX IF EXISTS idx_hero_section_is_default;
CREATE INDEX idx_hero_section_is_default ON hero_section(is_default) 
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_hero_section_tanggal;
CREATE INDEX idx_hero_section_tanggal ON hero_section(tanggal_mulai, tanggal_selesai) 
    WHERE deleted_at IS NULL;

-- Create new trigger
CREATE OR REPLACE FUNCTION fn_ensure_single_default_hero()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        UPDATE hero_section 
        SET is_default = false,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND is_default = true
          AND deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_single_default_hero ON hero_section;
CREATE TRIGGER trg_single_default_hero
    AFTER INSERT OR UPDATE OF is_default ON hero_section
    FOR EACH ROW
    WHEN (NEW.is_default = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_default_hero();

-- Update comments
COMMENT ON COLUMN hero_section.is_default IS 'Status default. Hanya 1 yang boleh true (enforced by trigger)';
COMMENT ON COLUMN hero_section.tanggal_mulai IS 'Tanggal mulai tampil (optional, untuk scheduled hero)';
COMMENT ON COLUMN hero_section.tanggal_selesai IS 'Tanggal selesai tampil (optional, untuk scheduled hero)';
