-- migrations/000052_remove_tipe_produk_manage_permission.down.sql
-- Restore tipe_produk:manage permission if needed

-- Restore the permission
INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('Manage Tipe Produk', 'tipe_produk:manage', 'master', 'CRUD tipe produk');

-- Restore original description for tipe_produk:read
UPDATE permission 
SET deskripsi = 'Melihat tipe produk' 
WHERE kode = 'tipe_produk:read';
