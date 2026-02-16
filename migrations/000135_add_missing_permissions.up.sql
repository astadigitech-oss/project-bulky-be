-- Add missing pesanan:delete permission
INSERT INTO permission (nama, kode, modul, deskripsi) 
VALUES ('Delete Pesanan', 'pesanan:delete', 'pesanan', 'Menghapus pesanan')
ON CONFLICT (kode) DO NOTHING;

-- Add ulasan:manage permission for consistency
INSERT INTO permission (nama, kode, modul, deskripsi) 
VALUES ('Manage Ulasan', 'ulasan:manage', 'ulasan', 'Mengelola ulasan (approve, reject, delete)')
ON CONFLICT (kode) DO NOTHING;

-- Assign new permissions to SUPER_ADMIN (all permissions)
INSERT INTO role_permission (role_id, permission_id)
SELECT 
    (SELECT id FROM role WHERE kode = 'SUPER_ADMIN'),
    id
FROM permission
WHERE kode IN ('pesanan:delete', 'ulasan:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign new permissions to ADMIN (all permissions except excluded ones)
INSERT INTO role_permission (role_id, permission_id)
SELECT 
    (SELECT id FROM role WHERE kode = 'ADMIN'),
    id
FROM permission
WHERE kode IN ('pesanan:delete', 'ulasan:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Note: STAFF tidak mendapat permission ini karena hanya read-only + specific actions
