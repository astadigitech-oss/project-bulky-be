-- migrations/000142_buyer_phone_auth.down.sql

-- =====================================================
-- STEP 1: Kembalikan kolom email_verified_at
-- =====================================================

ALTER TABLE buyer
    ADD COLUMN email_verified_at TIMESTAMP;

-- Restore dari telepon_verified_at
UPDATE buyer
SET email_verified_at = telepon_verified_at::TIMESTAMP
WHERE telepon_verified_at IS NOT NULL;

-- =====================================================
-- STEP 2: Hapus telepon_verified_at
-- =====================================================

ALTER TABLE buyer
    DROP COLUMN IF EXISTS telepon_verified_at;

-- =====================================================
-- STEP 3: Kembalikan email, password, username ke NOT NULL
-- =====================================================

UPDATE buyer SET email = CONCAT('unknown-', SUBSTRING(id::text, 1, 8), '@placeholder.com') WHERE email IS NULL;
UPDATE buyer SET password = '$2a$10$placeholder_hash_for_migration_rollback' WHERE password IS NULL;
UPDATE buyer SET username = CONCAT('user-', SUBSTRING(id::text, 1, 8)) WHERE username IS NULL;

ALTER TABLE buyer
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password SET NOT NULL,
    ALTER COLUMN username SET NOT NULL;

-- =====================================================
-- STEP 4: Kembalikan telepon ke nullable
-- =====================================================

-- Hapus unique index baru
DROP INDEX IF EXISTS idx_buyer_telepon;

-- Hapus NOT NULL
ALTER TABLE buyer
    ALTER COLUMN telepon DROP NOT NULL;

-- Bersihkan placeholder data dari STEP 1 migration up
UPDATE buyer
SET telepon = NULL
WHERE telepon LIKE 'UNKNOWN-%';

-- Buat kembali partial index lama
CREATE INDEX idx_buyer_telepon
    ON buyer(telepon)
    WHERE telepon IS NOT NULL;

-- =====================================================
-- STEP 5: Restore comments
-- =====================================================

COMMENT ON COLUMN buyer.telepon IS 'Nomor telepon buyer (opsional)';
COMMENT ON COLUMN buyer.email IS 'Email unik untuk login dan notifikasi';
COMMENT ON COLUMN buyer.password IS 'Password yang sudah di-hash menggunakan bcrypt';
COMMENT ON COLUMN buyer.username IS 'Username unik untuk login';
COMMENT ON COLUMN buyer.is_verified IS 'Status verifikasi email';
