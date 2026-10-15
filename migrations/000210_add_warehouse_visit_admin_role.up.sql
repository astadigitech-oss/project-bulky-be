INSERT INTO role (nama, kode, deskripsi)
VALUES (
    'Admin Gudang',
    'ADMIN_GUDANG',
    'Mengelola jadwal dan reservasi kunjungan gudang'
)
ON CONFLICT (kode) DO NOTHING;

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('Lihat Kunjungan Gudang', 'warehouse_visit:read', 'warehouse_visit', 'Melihat jadwal dan reservasi kunjungan gudang'),
    ('Kelola Kunjungan Gudang', 'warehouse_visit:manage', 'warehouse_visit', 'Mengatur jadwal, memindahkan atau membatalkan reservasi, dan mencatat kehadiran')
ON CONFLICT (kode) DO NOTHING;

-- Keep this dedicated role limited to the warehouse-visit permissions below.
DELETE FROM role_permission
WHERE role_id = (SELECT id FROM role WHERE kode = 'ADMIN_GUDANG');

-- Keep existing administrator roles able to manage this feature. The new
-- ADMIN_GUDANG role receives only these two permissions.
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE (r.kode IN ('SUPER_ADMIN', 'ADMIN', 'ADMIN_GUDANG')
       AND p.kode = 'warehouse_visit:manage')
   OR (r.kode IN ('SUPER_ADMIN', 'ADMIN', 'ADMIN_GUDANG', 'STAFF')
       AND p.kode = 'warehouse_visit:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;
