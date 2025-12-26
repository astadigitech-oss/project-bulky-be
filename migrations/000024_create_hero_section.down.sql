DROP TRIGGER IF EXISTS trg_single_active_hero ON hero_section;
DROP TRIGGER IF EXISTS trg_hero_section_updated_at ON hero_section;
DROP FUNCTION IF EXISTS fn_ensure_single_active_hero();
DROP TABLE IF EXISTS hero_section;
