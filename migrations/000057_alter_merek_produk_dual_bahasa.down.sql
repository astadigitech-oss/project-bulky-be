-- Rollback: Remove nama_en and rename back
ALTER TABLE merek_produk DROP COLUMN IF EXISTS nama_en;
ALTER TABLE merek_produk RENAME COLUMN nama_id TO nama;
