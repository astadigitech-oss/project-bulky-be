INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Auction Education Banners', 'auction_education_banner:read', 'auction', 'Melihat banner edukasi lelang'),
    ('Manage Auction Education Banners', 'auction_education_banner:manage', 'auction', 'Mengelola banner edukasi lelang')
ON CONFLICT (kode) DO UPDATE SET
    nama = EXCLUDED.nama,
    modul = EXCLUDED.modul,
    deskripsi = EXCLUDED.deskripsi;

INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.kode IN ('SUPER_ADMIN', 'ADMIN', 'MARKETING')
  AND p.kode IN ('auction_education_banner:read', 'auction_education_banner:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;
