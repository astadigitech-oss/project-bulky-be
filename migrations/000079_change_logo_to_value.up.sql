-- =====================================================
-- CHANGE: Logo URL → Logo Value
-- =====================================================
-- Ubah field 'logo' dari URL ke value identifier
-- Frontend akan match value dengan local assets
-- =====================================================

-- 1. Rename column
ALTER TABLE metode_pembayaran 
RENAME COLUMN logo TO logo_value;

-- 2. Update existing data (extract filename from URL)
-- Example: 'payment-methods/bca.png' → 'bca'
UPDATE metode_pembayaran 
SET logo_value = LOWER(
    REPLACE(
        SUBSTRING(logo_value FROM '[^/]+$'),  -- Get filename
        '.png', ''                             -- Remove extension
    )
)
WHERE logo_value IS NOT NULL;

-- 3. Update column type & comment
ALTER TABLE metode_pembayaran 
ALTER COLUMN logo_value TYPE VARCHAR(50);

COMMENT ON COLUMN metode_pembayaran.logo_value IS 
'Logo identifier for matching with frontend assets. Example: "bca" → /assets/payment-logos/bca.png';

-- =====================================================
-- Verification
-- =====================================================
-- Check updated values
-- SELECT nama, kode, logo_value FROM metode_pembayaran LIMIT 10;
