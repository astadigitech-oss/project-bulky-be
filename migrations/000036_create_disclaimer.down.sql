DROP TRIGGER IF EXISTS trg_single_active_disclaimer ON disclaimer;
DROP TRIGGER IF EXISTS trg_generate_slug_disclaimer ON disclaimer;
DROP TRIGGER IF EXISTS trg_disclaimer_updated_at ON disclaimer;
DROP FUNCTION IF EXISTS fn_ensure_single_active_disclaimer();
DROP FUNCTION IF EXISTS generate_slug_disclaimer();
DROP TABLE IF EXISTS disclaimer;
