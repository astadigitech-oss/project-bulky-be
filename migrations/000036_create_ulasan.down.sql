-- migrations/000035_create_ulasan.down.sql

DROP TRIGGER IF EXISTS trg_set_ulasan_approved_at ON ulasan;
DROP TRIGGER IF EXISTS trg_validate_ulasan_order ON ulasan;
DROP TRIGGER IF EXISTS trg_ulasan_updated_at ON ulasan;
DROP FUNCTION IF EXISTS fn_set_ulasan_approved_at();
DROP FUNCTION IF EXISTS fn_validate_ulasan_order_completed();
DROP TABLE IF EXISTS ulasan;
