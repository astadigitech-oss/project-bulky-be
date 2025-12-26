DROP TRIGGER IF EXISTS trg_generate_order_code ON pesanan;
DROP TRIGGER IF EXISTS trg_pesanan_updated_at ON pesanan;
DROP FUNCTION IF EXISTS generate_order_code();
DROP TABLE IF EXISTS pesanan;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_type;
DROP TYPE IF EXISTS delivery_type;
