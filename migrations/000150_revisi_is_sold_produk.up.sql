-- migrations/000083_revisi_is_sold_produk.up.sql

-- Tambah kolom is_sold
ALTER TABLE produk
    ADD COLUMN is_sold BOOLEAN NOT NULL DEFAULT false;

-- Index untuk query filter produk tersedia
CREATE INDEX idx_produk_is_sold ON produk(is_sold) WHERE deleted_at IS NULL;

-- Drop kolom quantity_terjual
ALTER TABLE produk
    DROP COLUMN IF EXISTS quantity_terjual;

-- Hapus permission produk:stock (tidak digunakan lagi)
DELETE FROM permission WHERE kode = 'produk:stock';

COMMENT ON COLUMN produk.quantity IS 'Jumlah item/unit di dalam satu palet (isi palet)';
COMMENT ON COLUMN produk.is_sold IS 'True jika palet sudah terjual/diklaim oleh pesanan aktif';
