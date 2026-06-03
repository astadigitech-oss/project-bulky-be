-- migrations/000149_create_keranjang_item.down.sql

DROP TRIGGER IF EXISTS update_keranjang_item_updated_at ON keranjang_item;
DROP TABLE IF EXISTS keranjang_item;
