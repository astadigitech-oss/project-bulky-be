DROP TRIGGER IF EXISTS trg_calculate_subtotal ON pesanan_item;
DROP TRIGGER IF EXISTS trg_pesanan_item_updated_at ON pesanan_item;
DROP FUNCTION IF EXISTS calculate_pesanan_item_subtotal();
DROP TABLE IF EXISTS pesanan_item;
