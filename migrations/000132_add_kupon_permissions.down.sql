-- Remove Kupon permissions from role_permission
DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('kupon:read', 'kupon:manage')
);

-- Remove Kupon permissions
DELETE FROM permission WHERE kode IN ('kupon:read', 'kupon:manage');
