-- migrations/000167_add_finance_marketing_roles.up.sql
-- ============================================================
-- 1. Role baru: FINANCE & MARKETING
-- ============================================================
INSERT INTO role (nama, kode, deskripsi) VALUES
    ('Finance', 'FINANCE', 'Akses transaksi pesanan, pembayaran & laporan keuangan'),
    ('Marketing', 'MARKETING', 'Akses konten promosi: hero section, banner, blog, video, kupon, diskon')
ON CONFLICT (kode) DO NOTHING;

-- ============================================================
-- 2. Permission gap: diskon:* (dipakai routes.go tapi belum pernah di-seed)
-- ============================================================
INSERT INTO permission (nama, kode, modul, deskripsi) VALUES
    ('View Diskon Kategori', 'diskon:read', 'marketing', 'Melihat diskon kategori produk'),
    ('Manage Diskon Kategori', 'diskon:manage', 'marketing', 'CRUD diskon kategori produk')
ON CONFLICT (kode) DO NOTHING;

-- ============================================================
-- 3. Hapus permission usang yang tidak dipakai route mana pun:
--    pesanan:update, ulasan:approve, ulasan:delete
--    (produk:stock sudah dihapus di migration 000150)
-- ============================================================
DELETE FROM permission WHERE kode IN ('pesanan:update', 'ulasan:approve', 'ulasan:delete');

-- ============================================================
-- 4. Reset assignment SEMUA role (deterministik, whitelist eksplisit).
--    Ini menimpa assignment lama dari migration 000044/000103/000132/000135.
-- ============================================================
DELETE FROM role_permission;

-- --- SUPER_ADMIN: semua permission (all access) ---
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r CROSS JOIN permission p
WHERE r.kode = 'SUPER_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- --- ADMIN: produk, relasi & master data ---
-- Tidak punya: admin:*, role:manage, system:* (kecuali tipe_produk:read),
-- activity_log, pembayaran (FINANCE), diskon/kupon/marketing (MARKETING),
-- operasional:manage (STAFF), pesanan:delete & buyer:manage (SuperAdmin).
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r CROSS JOIN permission p
WHERE r.kode = 'ADMIN'
  AND p.kode IN (
      'dashboard:read',
      'kategori:read', 'kategori:manage',
      'brand:read', 'brand:manage',
      'kondisi:read', 'kondisi:manage',
      'tipe_produk:read',
      'produk:read', 'produk:create', 'produk:update', 'produk:delete',
      'ulasan:read', 'ulasan:manage',
      'pesanan:read'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- --- FINANCE: transaksi, pembayaran, laporan ---
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r CROSS JOIN permission p
WHERE r.kode = 'FINANCE'
  AND p.kode IN (
      'dashboard:read',
      'pesanan:read', 'pesanan:update_status',
      'pembayaran:read', 'pembayaran:manage',
      'buyer:read',
      'kupon:read', 'diskon:read',
      'ulasan:read',
      'activity_log:read'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- --- MARKETING: konten & promosi ---
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r CROSS JOIN permission p
WHERE r.kode = 'MARKETING'
  AND p.kode IN (
      'dashboard:read',
      'marketing:read', 'marketing:manage',
      'kupon:read', 'kupon:manage',
      'diskon:read', 'diskon:manage',
      'ulasan:read', 'ulasan:manage',
      'produk:read', 'kategori:read', 'brand:read', 'kondisi:read',
      'tipe_produk:read'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- --- STAFF: read-only + update status pesanan ---
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r CROSS JOIN permission p
WHERE r.kode = 'STAFF'
  AND (
      p.kode IN (
          'dashboard:read',
          'kategori:read', 'brand:read', 'kondisi:read', 'tipe_produk:read',
          'produk:read', 'ulasan:read', 'pesanan:read', 'buyer:read',
          'operasional:read', 'pembayaran:read', 'system:read', 'faq:read',
          'marketing:read', 'kupon:read', 'activity_log:read'
      )
      OR p.kode = 'pesanan:update_status'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================
-- Catatan: route panel/tipe-produk diubah dari system:read
-- → tipe_produk:read di internal/routes/routes.go
-- ============================================================
