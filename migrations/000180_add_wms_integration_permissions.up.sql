-- migrations/000180_add_wms_integration_permissions.up.sql
-- Permission untuk integrasi WMS (OAuth token exchange, cek koneksi, dan
-- nantinya sync produk palet dari inventory WMS jadi cargo online).
-- SUPER_ADMIN only.

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View WMS Integration', 'wms_integration:read', 'integration', 'Melihat status integrasi WMS'),
    ('Manage WMS Integration', 'wms_integration:manage', 'integration', 'Cek koneksi & kelola sync data dari WMS')
ON CONFLICT (kode) DO NOTHING;

-- Grant HANYA ke Super Admin (bukan cross-join ke semua role)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('wms_integration:read', 'wms_integration:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;
