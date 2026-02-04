-- Remove FAQ permissions from role_permission
DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('faq:read', 'faq:manage')
);

-- Remove FAQ permissions
DELETE FROM permission WHERE kode IN ('faq:read', 'faq:manage');
