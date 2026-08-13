-- migrations/000179_add_forwarder_mapping_permissions.up.sql
-- Permission untuk master data Forwarder Mapping (SUPER_ADMIN only)

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Forwarder Mapping', 'forwarder_mapping:read', 'shipping', 'Melihat master data mapping kota & kecamatan Forwarder'),
    ('Manage Forwarder Mapping', 'forwarder_mapping:manage', 'shipping', 'Sync & kelola master data mapping kota & kecamatan Forwarder')
ON CONFLICT (kode) DO NOTHING;

-- Grant HANYA ke Super Admin (bukan cross-join ke semua role)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('forwarder_mapping:read', 'forwarder_mapping:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;
