-- migrations/000173_add_deliveree_vehicle_permissions.up.sql
-- Permission untuk master data Deliveree Vehicle Type (SUPER_ADMIN only)

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Deliveree Vehicle', 'deliveree_vehicle:read', 'master', 'Melihat master data kendaraan Deliveree'),
    ('Manage Deliveree Vehicle', 'deliveree_vehicle:manage', 'master', 'Sync & kelola master data kendaraan Deliveree');

-- Grant HANYA ke Super Admin (bukan cross-join ke semua role)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('deliveree_vehicle:read', 'deliveree_vehicle:manage');
