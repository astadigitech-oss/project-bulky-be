DROP TRIGGER IF EXISTS trg_single_active_ppn ON ppn;
DROP TRIGGER IF EXISTS trg_ppn_updated_at ON ppn;
DROP FUNCTION IF EXISTS fn_ensure_single_active_ppn();
DROP TABLE IF EXISTS ppn;
