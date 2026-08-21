-- migrations/000180_add_wms_integration_permissions.down.sql

DELETE FROM role_permission WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('wms_integration:read', 'wms_integration:manage')
);
DELETE FROM permission WHERE kode IN ('wms_integration:read', 'wms_integration:manage');
