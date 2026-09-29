DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission
    WHERE kode IN ('syarat_ketentuan_lelang:read', 'syarat_ketentuan_lelang:manage')
);

DELETE FROM permission
WHERE kode IN ('syarat_ketentuan_lelang:read', 'syarat_ketentuan_lelang:manage');
