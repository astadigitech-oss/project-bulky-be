-- migrations/000179_add_forwarder_mapping_permissions.down.sql
DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('forwarder_mapping:read', 'forwarder_mapping:manage')
);

DELETE FROM permission WHERE kode IN ('forwarder_mapping:read', 'forwarder_mapping:manage');
