-- =====================================================
-- SEED: Metode Pembayaran
-- =====================================================
-- Insert payment methods with Xendit channel codes
-- Total: 14 payment methods
-- - 8 Bank Transfer / VA
-- - 4 E-Wallet
-- - 1 Credit Card
-- - 1 QRIS
-- =====================================================

-- =====================================================
-- BANK TRANSFER / VA (8 methods)
-- =====================================================

-- BCA
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'BCA',
    'BCA',
    'bca',
    1,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'BCA');

-- MANDIRI
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'MANDIRI',
    'MANDIRI',
    'mandiri',
    2,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'MANDIRI');

-- BNI
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'BNI',
    'BNI',
    'bni',
    3,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'BNI');

-- BRI
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'BRI',
    'BRI',
    'bri',
    4,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'BRI');

-- PERMATA
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'PERMATA',
    'PERMATA',
    'permata',
    5,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'PERMATA');

-- CIMB NIAGA
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'CIMB NIAGA',
    'CIMB',
    'cimb',
    6,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'CIMB');

-- BSI
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'BSI',
    'BSI',
    'bsi',
    7,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'BSI');

-- BJB
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'),
    'BJB',
    'BJB',
    'bjb',
    8,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'BJB');

-- =====================================================
-- E-WALLET (4 methods)
-- =====================================================

-- LINK AJA
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'),
    'LINK AJA',
    'LINKAJA',
    'linkaja',
    1,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'LINKAJA');

-- DANA
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'),
    'DANA',
    'DANA',
    'dana',
    2,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'DANA');

-- OVO
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'),
    'OVO',
    'ID_OVO',
    'ovo',
    3,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'ID_OVO');

-- ShopeePay
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'),
    'ShopeePay',
    'ID_SHOPEEPAY',
    'shopeepay',
    4,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'ID_SHOPEEPAY');

-- =====================================================
-- KARTU KREDIT (1 method)
-- =====================================================

-- KARTU KREDIT
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'Kartu Kredit'),
    'KARTU KREDIT',
    'CREDIT_CARD',
    'credit-card',
    1,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'CREDIT_CARD');

-- =====================================================
-- QRIS (1 method)
-- =====================================================

-- QRIS
INSERT INTO metode_pembayaran (id, group_id, nama, kode, logo_value, urutan, is_active, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM metode_pembayaran_group WHERE nama = 'QRIS'),
    'QRIS',
    'ID_QRIS',
    'qris',
    1,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM metode_pembayaran WHERE kode = 'ID_QRIS');

-- Add comment
COMMENT ON TABLE metode_pembayaran IS 'Payment methods with Xendit channel codes for integration';
