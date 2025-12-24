DROP TRIGGER IF EXISTS trigger_first_alamat_buyer_default ON alamat_buyer;
DROP TRIGGER IF EXISTS trigger_single_default_alamat_buyer ON alamat_buyer;
DROP TRIGGER IF EXISTS update_alamat_buyer_updated_at ON alamat_buyer;
DROP FUNCTION IF EXISTS set_first_alamat_buyer_as_default();
DROP FUNCTION IF EXISTS ensure_single_default_alamat_buyer();
DROP TABLE IF EXISTS alamat_buyer;
