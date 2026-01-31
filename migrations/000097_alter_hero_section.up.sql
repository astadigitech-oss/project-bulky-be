-- 1. Drop old indexes first
DROP INDEX IF EXISTS idx_hero_section_is_active;
DROP INDEX IF EXISTS idx_hero_section_urutan;

-- 2. Drop old trigger and function
DROP TRIGGER IF EXISTS trg_single_active_hero ON hero_section;
DROP FUNCTION IF EXISTS fn_ensure_single_active_hero();

-- 3. Hapus field urutan (only if exists)
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'urutan'
    ) THEN
        ALTER TABLE hero_section DROP COLUMN urutan;
    END IF;
END $$;

-- 4. Rename is_active → is_default (only if is_active exists)
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'is_active'
    ) THEN
        ALTER TABLE hero_section RENAME COLUMN is_active TO is_default;
    END IF;
END $$;

-- 5. Add tanggal_mulai and tanggal_selesai if not exists
ALTER TABLE hero_section ADD COLUMN IF NOT EXISTS tanggal_mulai TIMESTAMP NULL;
ALTER TABLE hero_section ADD COLUMN IF NOT EXISTS tanggal_selesai TIMESTAMP NULL;

-- 6. Create new indexes
DROP INDEX IF EXISTS idx_hero_section_is_default;
CREATE INDEX idx_hero_section_is_default ON hero_section(is_default) 
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_hero_section_tanggal;
CREATE INDEX idx_hero_section_tanggal ON hero_section(tanggal_mulai, tanggal_selesai) 
    WHERE deleted_at IS NULL;

-- 7. Create new trigger function
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

-- 8. Update comments
COMMENT ON COLUMN hero_section.is_default IS 'Status default. Hanya 1 yang boleh true (enforced by trigger)';
COMMENT ON COLUMN hero_section.tanggal_mulai IS 'Tanggal mulai tampil (optional, untuk scheduled hero)';
COMMENT ON COLUMN hero_section.tanggal_selesai IS 'Tanggal selesai tampil (optional, untuk scheduled hero)';
