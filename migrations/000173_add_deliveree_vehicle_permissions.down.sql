DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('deliveree_vehicle:read', 'deliveree_vehicle:manage')
);

DELETE FROM permission WHERE kode IN ('deliveree_vehicle:read', 'deliveree_vehicle:manage');
