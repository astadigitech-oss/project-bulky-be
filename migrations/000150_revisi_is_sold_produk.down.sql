-- migrations/000083_revisi_is_sold_produk.down.sql

DROP INDEX IF EXISTS idx_produk_is_sold;

ALTER TABLE produk
    DROP COLUMN IF EXISTS is_sold;

ALTER TABLE produk
    ADD COLUMN quantity_terjual INTEGER DEFAULT 0 CHECK (quantity_terjual >= 0);

-- Restore permission produk:stock
INSERT INTO permission (nama, kode, modul, deskripsi)
VALUES ('Manage Stock', 'produk:stock', 'produk', 'Mengatur stok produk')
ON CONFLICT (kode) DO NOTHING;
