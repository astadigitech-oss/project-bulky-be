-- =====================================================
-- Normalisasi timezone: timestamp without time zone → timestamptz
-- =====================================================
-- Latar belakang: standar baru adalah semua timestamp disimpan sebagai UTC
-- (DB server TimeZone=Etc/UTC, GORM NowFunc time.Now().UTC(), Dockerfile
-- TZ=UTC). Namun beberapa kolom lama ditulis sebagai wall-clock Asia/Jakarta
-- (WIB), sehingga saat dibaca sebagai UTC muncul pergeseran 7 jam.
--
-- Aturan shift per kolom:
--   * WIB-stored  → AT TIME ZONE 'Asia/Jakarta'  (interpretasi ulang: nilai
--     wall-clock dianggap WIB, jadikan timestamptz = -7 jam dari nilai asli)
--   * UTC-stored  → AT TIME ZONE 'UTC'           (nilai wall-clock sudah UTC)
--
-- ⚠️  SEBELUM EKSEKUSI: backup dulu (pg_dump), lalu jalankan SELECT verifikasi
--     hasil per tabel. Jangan dijalankan sembarangan di produksi.

-- -----------------------------------------------------
-- banner_event_promo
--   tanggal_mulai/tanggal_selesai = WIB (parseFlexibleDate, ParseInLocation
--   Asia/Jakarta di banner_event_promo_service.go). Data contoh: 2026-06-01
--   00:00:00 = 2026-05-31 17:00 UTC yang benar.
--   created_at/updated_at/deleted_at = UTC (NowFunc, semua data > Feb 2026).
-- -----------------------------------------------------
ALTER TABLE banner_event_promo
    ALTER COLUMN tanggal_mulai TYPE TIMESTAMPTZ USING tanggal_mulai AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN tanggal_selesai TYPE TIMESTAMPTZ USING tanggal_selesai AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN deleted_at TYPE TIMESTAMPTZ USING deleted_at AT TIME ZONE 'UTC';

-- -----------------------------------------------------
-- admin — MIXED:
--   * last_login_at = WIB (auth_service.go pakai time.Now(), Dockerfile lama
--     TZ=Asia/Jakarta) → shift -7 jam.
--   * created_at/updated_at/deleted_at = campuran: baris sebelum 12 Feb 2026
--     ditulis saat DSN TimeZone=Asia/Jakarta (WIB), sesudahnya UTC (NowFunc).
--     Pakai CASE dengan threshold 12 Feb 2026 (commit 37a0a46).
-- -----------------------------------------------------
ALTER TABLE admin
    ALTER COLUMN last_login_at TYPE TIMESTAMPTZ USING last_login_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING
        CASE WHEN created_at < '2026-02-12 00:00:00' THEN created_at AT TIME ZONE 'Asia/Jakarta'
             ELSE created_at AT TIME ZONE 'UTC' END,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING
        CASE WHEN updated_at < '2026-02-12 00:00:00' THEN updated_at AT TIME ZONE 'Asia/Jakarta'
             ELSE updated_at AT TIME ZONE 'UTC' END,
    ALTER COLUMN deleted_at TYPE TIMESTAMPTZ USING
        CASE WHEN deleted_at < '2026-02-12 00:00:00' THEN deleted_at AT TIME ZONE 'Asia/Jakarta'
             ELSE deleted_at AT TIME ZONE 'UTC' END;

-- -----------------------------------------------------
-- activity_log
--   created_at = WIB (auth_v2_service.go logActivity pakai time.Now()) → shift.
-- -----------------------------------------------------
ALTER TABLE activity_log
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'Asia/Jakarta';

-- -----------------------------------------------------
-- formulir_partai_besar_submission
--   email_sent_at = WIB (formulir_partai_besar_service.go pakai time.Now())
--     → shift. (data saat ini kosong, aman)
--   created_at = UTC (GORM autoCreateTime → NowFunc UTC) → tanpa shift.
-- -----------------------------------------------------
ALTER TABLE formulir_partai_besar_submission
    ALTER COLUMN email_sent_at TYPE TIMESTAMPTZ USING email_sent_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

-- Verifikasi cepat setelah migrasi (harus dibaca sebagai UTC):
--   SELECT id, tanggal_mulai, tanggal_selesai FROM banner_event_promo
--   ORDER BY created_at LIMIT 5;
--   2026-06-01 00:00:00 (input WIB) harus menjadi 2026-05-31 17:00:00+00.
