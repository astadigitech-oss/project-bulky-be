-- =====================================================
-- SEED: Metode Pembayaran Group
-- =====================================================
-- Insert payment method groups
-- Based on Xendit payment categories
-- =====================================================

-- Bank Transfer / VA
INSERT INTO metode_pembayaran_group (id, nama, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Bank Transfer / VA',
    1,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (nama) DO UPDATE SET
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- E-Wallet
INSERT INTO metode_pembayaran_group (id, nama, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'E-Wallet',
    2,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (nama) DO UPDATE SET
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Kartu Kredit
INSERT INTO metode_pembayaran_group (id, nama, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Kartu Kredit',
    3,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (nama) DO UPDATE SET
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- QRIS
INSERT INTO metode_pembayaran_group (id, nama, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'QRIS',
    4,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (nama) DO UPDATE SET
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Add comment
COMMENT ON TABLE metode_pembayaran_group IS 'Payment method groups based on Xendit categories';
