-- ============================================================
-- PRD 01 / 04A — Rollback permission modul Lelang
-- ============================================================

DELETE FROM role_permission
WHERE permission_id IN (
    SELECT id FROM permission WHERE kode IN ('auction:read', 'auction:manage')
);

DELETE FROM permission WHERE kode IN ('auction:read', 'auction:manage');
