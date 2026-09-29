INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('Lihat Syarat dan Ketentuan Lelang', 'syarat_ketentuan_lelang:read', 'auction', 'Melihat syarat dan ketentuan lelang'),
    ('Kelola Syarat dan Ketentuan Lelang', 'syarat_ketentuan_lelang:manage', 'auction', 'Mengubah syarat dan ketentuan lelang')
ON CONFLICT (kode) DO UPDATE SET
    nama = EXCLUDED.nama,
    modul = EXCLUDED.modul,
    deskripsi = EXCLUDED.deskripsi;

INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r
CROSS JOIN permission p
WHERE r.kode IN ('SUPER_ADMIN', 'ADMIN')
  AND p.kode IN ('syarat_ketentuan_lelang:read', 'syarat_ketentuan_lelang:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;
