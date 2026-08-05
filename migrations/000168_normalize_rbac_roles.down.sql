-- migrations/000168_normalize_rbac_roles.down.sql
-- Rollback migration 000168.
-- CATATAN: tidak bisa mengembalikan assignment role_permission lama secara
-- otomatis (state sebelum 000167/000168 tidak disimpan di sini). Jika perlu
-- memulihkan, jalankan ulang migration 000044/000103/000132/000135 ATAU
-- restore backup DB.

DELETE FROM role_permission;

DELETE FROM role WHERE kode IN ('FINANCE', 'MARKETING');

DELETE FROM permission WHERE kode IN ('diskon:read', 'diskon:manage');
