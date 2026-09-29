DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission
    WHERE kode IN ('auction_education_banner:read', 'auction_education_banner:manage')
);

DELETE FROM permission
WHERE kode IN ('auction_education_banner:read', 'auction_education_banner:manage');
