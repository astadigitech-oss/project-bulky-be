-- Hapus penanda ribbon "SALE" pada tabel produk.
ALTER TABLE produk
    DROP COLUMN IF EXISTS is_sale;
