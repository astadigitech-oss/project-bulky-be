DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission
    WHERE kode IN ('warehouse_visit:read', 'warehouse_visit:manage')
);

DELETE FROM permission
WHERE kode IN ('warehouse_visit:read', 'warehouse_visit:manage');

DELETE FROM role
WHERE kode = 'ADMIN_GUDANG';
