-- =====================================================
-- ROLLBACK: Revert to logo URL
-- =====================================================

-- 1. Reconstruct URLs from values
UPDATE metode_pembayaran 
SET logo_value = 'payment-methods/' || logo_value || '.png'
WHERE logo_value IS NOT NULL;

-- 2. Rename back
ALTER TABLE metode_pembayaran 
RENAME COLUMN logo_value TO logo;

-- 3. Update column type
ALTER TABLE metode_pembayaran 
ALTER COLUMN logo TYPE VARCHAR(255);
