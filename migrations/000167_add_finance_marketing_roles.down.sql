-- migrations/000167_add_finance_marketing_roles.down.sql
-- Rollback: hapus role & permission baru yang dibuat migration ini.
-- Assignment role_permission TIDAK dikembalikan otomatis — jika perlu
-- memulihkan assignment lama, jalankan ulang migration 000044, 000103,
-- 000132, 000135 (urutan awal) ATAU restore backup DB.

-- Hapus assignment FINANCE & MARKETING
DELETE FROM role_permission
WHERE role_id IN (SELECT id FROM role WHERE kode IN ('FINANCE', 'MARKETING'));

-- Hapus role FINANCE & MARKETING
DELETE FROM role WHERE kode IN ('FINANCE', 'MARKETING');

-- Hapus permission baru diskon:*
-- (tipe_produk:read TIDAK dihapus karena sudah ada sejak 000043)
DELETE FROM permission WHERE kode IN ('diskon:read', 'diskon:manage');
