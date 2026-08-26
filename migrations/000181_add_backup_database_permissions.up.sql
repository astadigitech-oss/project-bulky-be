-- migrations/000181_add_backup_database_permissions.up.sql
-- Permission untuk modul Backup Database di Admin Panel.
-- Hanya diberikan (grant) ke role Super Admin.

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Database Backups', 'backup:read', 'backup', 'Melihat daftar riwayat file backup database'),
    ('Create Database Backup', 'backup:create', 'backup', 'Memicu proses backup database manual'),
    ('Download Database Backup', 'backup:download', 'backup', 'Mengunduh file dump backup database'),
    ('Delete Database Backup', 'backup:delete', 'backup', 'Menghapus file backup database')
ON CONFLICT (kode) DO NOTHING;

-- Grant HANYA ke Super Admin
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('backup:read', 'backup:create', 'backup:download', 'backup:delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;
