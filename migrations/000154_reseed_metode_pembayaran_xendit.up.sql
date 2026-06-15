-- =====================================================
-- RESEED: Metode Pembayaran (Xendit Channel Codes)
-- =====================================================
-- Bersihkan data lama yang duplikat/inkonsisten dan
-- re-seed dengan kode channel Xendit yang benar sesuai
-- dokumentasi resmi Xendit Payment Method API.
--
-- Total: 16 metode pembayaran
--   8 Bank Transfer / VA
--   6 E-Wallet
--   1 Kartu Kredit (Cards)
--   1 QRIS
-- =====================================================

-- Step 1: Nullify FK di pesanan_pembayaran (kolom nullable)
UPDATE pesanan_pembayaran
SET metode_pembayaran_id = NULL
WHERE metode_pembayaran_id IS NOT NULL;

-- Step 2: Hapus semua data metode pembayaran lama
DELETE FROM metode_pembayaran;

-- Step 3: Hapus semua data group lama
DELETE FROM metode_pembayaran_group;

-- Step 4: Re-seed group
INSERT INTO metode_pembayaran_group (id, nama, urutan, is_active, created_at, updated_at) VALUES
    (uuid_generate_v4(), 'Bank Transfer / VA', 1, true, NOW(), NOW()),
    (uuid_generate_v4(), 'E-Wallet',           2, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Kartu Kredit',        3, true, NOW(), NOW()),
    (uuid_generate_v4(), 'QRIS',               4, true, NOW(), NOW());

-- Step 5: Re-seed metode pembayaran dengan Xendit channel code resmi
-- =====================================================
-- BANK TRANSFER / VA
-- Referensi: https://docs.xendit.co/payment-methods/bank-transfer
-- =====================================================
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at) VALUES
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BCA',        'BCA_VIRTUAL_ACCOUNT',      'bca',      1, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'Mandiri',    'MANDIRI_VIRTUAL_ACCOUNT',  'mandiri',  2, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BNI',        'BNI_VIRTUAL_ACCOUNT',      'bni',      3, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BRI',        'BRI_VIRTUAL_ACCOUNT',      'bri',      4, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'Permata',    'PERMATA_VIRTUAL_ACCOUNT',  'permata',  5, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'CIMB Niaga', 'CIMB_VIRTUAL_ACCOUNT',     'cimb',     6, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BSI',        'BSI_VIRTUAL_ACCOUNT',      'bsi',      7, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BJB',        'BJB_VIRTUAL_ACCOUNT',      'bjb',      8, true, NOW(), NOW());

-- =====================================================
-- E-WALLET
-- Referensi: https://docs.xendit.co/payment-methods/e-wallet
-- =====================================================
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at) VALUES
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'GoPay',     'GOPAY',     'gopay',     1, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'OVO',       'OVO',       'ovo',       2, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'Dana',      'DANA',      'dana',      3, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'LinkAja',   'LINKAJA',   'linkaja',   4, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'ShopeePay', 'SHOPEEPAY', 'shopeepay', 5, true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'Akulaku',   'AKULAKU',   'akulaku',   6, true, NOW(), NOW());

-- =====================================================
-- KARTU KREDIT / DEBIT
-- Referensi: https://docs.xendit.co/payment-methods/cards
-- =====================================================
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at) VALUES
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'Kartu Kredit'), 'Kartu Kredit / Debit', 'CARDS', 'credit-card', 1, true, NOW(), NOW());

-- =====================================================
-- QRIS
-- Referensi: https://docs.xendit.co/payment-methods/qr-code
-- =====================================================
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at) VALUES
    (uuid_generate_v4(), (SELECT id FROM metode_pembayaran_group WHERE nama = 'QRIS'), 'QRIS', 'QRIS', 'qris', 1, true, NOW(), NOW());
