-- migrations/000181_add_backup_database_permissions.down.sql
-- Rollback permission modul Backup Database

DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN (
        'backup:read',
        'backup:create',
        'backup:download',
        'backup:delete'
    )
);

DELETE FROM permission
WHERE kode IN (
    'backup:read',
    'backup:create',
    'backup:download',
    'backup:delete'
);
