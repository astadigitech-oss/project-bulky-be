ALTER TABLE kategori_produk DROP COLUMN IF EXISTS nama_en;
ALTER TABLE kategori_produk RENAME COLUMN nama_id TO nama;
