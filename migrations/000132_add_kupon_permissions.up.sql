-- Add Kupon permissions
INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Kupon', 'kupon:read', 'marketing', 'Melihat daftar kupon'),
    ('Manage Kupon', 'kupon:manage', 'marketing', 'CRUD kupon');

-- Grant Kupon permissions to Super Admin role
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('kupon:read', 'kupon:manage');

-- Grant Kupon permissions to Admin role (based on role_id)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Admin'
AND p.kode IN ('kupon:read', 'kupon:manage');

-- Grant Kupon permissions to Staff role (based on role_id)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Staff'
AND p.kode IN ('kupon:read', 'kupon:manage');