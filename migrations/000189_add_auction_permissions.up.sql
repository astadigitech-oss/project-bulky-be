-- ============================================================
-- PRD 01 / 04A — Permission modul Lelang (Auction)
-- Permission baru: auction:read dan auction:manage.
-- Role operasional (Super Admin & Admin) diberi keduanya.
-- ============================================================

INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Auction', 'auction:read', 'auction', 'Melihat batch lelang, bid, winner dan metrik'),
    ('Manage Auction', 'auction:manage', 'auction', 'Membuat/edit/publish batch, memilih winner dan operasi manual')
ON CONFLICT (kode) DO NOTHING;

-- Grant ke Super Admin & Admin (role operasional)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.kode IN ('SUPER_ADMIN', 'ADMIN')
AND p.kode IN ('auction:read', 'auction:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;
