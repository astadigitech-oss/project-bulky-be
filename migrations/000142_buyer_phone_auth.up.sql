-- migrations/000142_buyer_phone_auth.up.sql

-- =====================================================
-- STEP 1: Handle existing data sebelum ALTER
-- =====================================================

-- Pastikan tidak ada buyer dengan telepon NULL sebelum set NOT NULL.
-- Jika ada, set placeholder dulu (akan perlu diisi manual via admin panel).
UPDATE buyer
SET telepon = CONCAT('UNKNOWN-', SUBSTRING(id::text, 1, 8))
WHERE telepon IS NULL
  AND deleted_at IS NULL;

-- =====================================================
-- STEP 2: Ubah kolom telepon jadi NOT NULL + UNIQUE
-- =====================================================

-- Hapus partial index lama
DROP INDEX IF EXISTS idx_buyer_telepon;

-- Set NOT NULL
ALTER TABLE buyer
    ALTER COLUMN telepon SET NOT NULL;

-- Tambah unique constraint (partial — exclude soft deleted)
CREATE UNIQUE INDEX idx_buyer_telepon
    ON buyer(telepon)
    WHERE deleted_at IS NULL;

-- =====================================================
-- STEP 3: Email, password, username jadi nullable
-- (backward compat — data lama tidak hilang)
-- =====================================================

ALTER TABLE buyer
    ALTER COLUMN email DROP NOT NULL,
    ALTER COLUMN password DROP NOT NULL,
    ALTER COLUMN username DROP NOT NULL;

-- =====================================================
-- STEP 4: Tambah kolom telepon_verified_at
-- =====================================================

ALTER TABLE buyer
    ADD COLUMN telepon_verified_at TIMESTAMPTZ;

-- Buyer yang sudah is_verified = true, anggap telepon sudah verified
-- (migrasi semantics dari email_verified ke telepon_verified)
UPDATE buyer
SET telepon_verified_at = email_verified_at
WHERE is_verified = true
  AND email_verified_at IS NOT NULL;

-- =====================================================
-- STEP 5: Hapus kolom email_verified_at
-- =====================================================

ALTER TABLE buyer
    DROP COLUMN IF EXISTS email_verified_at;

-- =====================================================
-- STEP 6: Update comments
-- =====================================================

COMMENT ON COLUMN buyer.telepon IS 'Nomor telepon buyer format E.164 (+62xxx). Identifier utama untuk login via WA OTP. NOT NULL dan UNIQUE per buyer aktif.';
COMMENT ON COLUMN buyer.email IS 'Email buyer (opsional). Dipertahankan untuk backward compatibility data lama.';
COMMENT ON COLUMN buyer.password IS 'Password buyer (opsional). Dipertahankan untuk backward compatibility data lama.';
COMMENT ON COLUMN buyer.username IS 'Username buyer (opsional). Dipertahankan untuk backward compatibility data lama.';
COMMENT ON COLUMN buyer.is_verified IS 'Status verifikasi nomor telepon via OTP WA. true = nomor sudah terverifikasi.';
COMMENT ON COLUMN buyer.telepon_verified_at IS 'Timestamp pertama kali nomor telepon berhasil diverifikasi via OTP WA.';
