ALTER TABLE sumber_produk DROP COLUMN IF EXISTS nama_en;
ALTER TABLE sumber_produk RENAME COLUMN nama_id TO nama;
