-- Remove role_permission assignments
DELETE FROM role_permission 
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('pesanan:delete', 'ulasan:manage')
);

-- Remove added permissions
DELETE FROM permission WHERE kode IN ('pesanan:delete', 'ulasan:manage');
