-- Drop triggers
DROP TRIGGER IF EXISTS trg_prevent_delete_default ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_first_alamat_default ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_single_default_alamat ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_alamat_buyer_updated_at ON alamat_buyer;

-- Drop functions
DROP FUNCTION IF EXISTS fn_prevent_delete_default_alamat();
DROP FUNCTION IF EXISTS fn_first_alamat_as_default();
DROP FUNCTION IF EXISTS fn_ensure_single_default_alamat();

-- Drop table
DROP TABLE IF EXISTS alamat_buyer;
