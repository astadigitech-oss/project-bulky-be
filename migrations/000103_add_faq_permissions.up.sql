-- Add FAQ permissions
INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View FAQ', 'faq:read', 'operasional', 'Melihat FAQ'),
    ('Manage FAQ', 'faq:manage', 'operasional', 'Mengelola FAQ (CRUD & reorder)');

-- Grant FAQ permissions to Super Admin role
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.nama = 'Super Admin'
AND p.kode IN ('faq:read', 'faq:manage');
