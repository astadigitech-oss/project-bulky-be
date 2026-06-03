-- migrations/000148_create_keranjang.down.sql

DROP TRIGGER IF EXISTS update_keranjang_updated_at ON keranjang;
DROP TABLE IF EXISTS keranjang;
