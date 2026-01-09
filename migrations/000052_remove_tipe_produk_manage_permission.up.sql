-- migrations/000052_remove_tipe_produk_manage_permission.up.sql
-- Remove tipe_produk:manage permission as tipe produk is now read-only

-- Remove permission from role_permission junction table first
DELETE FROM role_permission 
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode = 'tipe_produk:manage'
);

-- Remove the permission itself
DELETE FROM permission WHERE kode = 'tipe_produk:manage';

-- Update description for tipe_produk:read to clarify it's read-only
UPDATE permission 
SET deskripsi = 'Melihat tipe produk (read-only)' 
WHERE kode = 'tipe_produk:read';
