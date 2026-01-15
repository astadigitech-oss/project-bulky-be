ALTER TABLE kondisi_produk DROP COLUMN IF EXISTS nama_en;
ALTER TABLE kondisi_produk RENAME COLUMN nama_id TO nama;
