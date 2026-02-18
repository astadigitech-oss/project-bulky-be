-- migrations/000137_fix_timestamp_to_timestamptz.up.sql
--
-- Perbaikan tipe data timestamp di tabel pesanan, ulasan, dan kupon:
-- 1. Semua kolom yang masih TIMESTAMP → TIMESTAMPTZ (agar timezone-aware)
-- 2. Hapus DEFAULT NOW()/CURRENT_TIMESTAMP dari created_at/updated_at
--    supaya nilai dari FE (UTC+0) tersimpan apa adanya tanpa konversi oleh DB.
--    Pengisian nilai dilakukan di application layer.

-- =====================================================
-- TABEL: pesanan
-- =====================================================

-- Status timestamps (tidak ada default, nullable)
ALTER TABLE pesanan
    ALTER COLUMN expired_at    TYPE TIMESTAMPTZ USING expired_at AT TIME ZONE 'UTC',
    ALTER COLUMN paid_at       TYPE TIMESTAMPTZ USING paid_at AT TIME ZONE 'UTC',
    ALTER COLUMN processed_at  TYPE TIMESTAMPTZ USING processed_at AT TIME ZONE 'UTC',
    ALTER COLUMN ready_at      TYPE TIMESTAMPTZ USING ready_at AT TIME ZONE 'UTC',
    ALTER COLUMN shipped_at    TYPE TIMESTAMPTZ USING shipped_at AT TIME ZONE 'UTC',
    ALTER COLUMN completed_at  TYPE TIMESTAMPTZ USING completed_at AT TIME ZONE 'UTC',
    ALTER COLUMN cancelled_at  TYPE TIMESTAMPTZ USING cancelled_at AT TIME ZONE 'UTC';

-- Standard timestamps: ubah tipe + hapus default
ALTER TABLE pesanan
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at DROP DEFAULT,
    ALTER COLUMN deleted_at TYPE TIMESTAMPTZ USING deleted_at AT TIME ZONE 'UTC';

-- =====================================================
-- TABEL: ulasan
-- =====================================================

-- approved_at
ALTER TABLE ulasan
    ALTER COLUMN approved_at TYPE TIMESTAMPTZ USING approved_at AT TIME ZONE 'UTC';

-- Standard timestamps: ubah tipe + hapus default
ALTER TABLE ulasan
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at DROP DEFAULT,
    ALTER COLUMN deleted_at TYPE TIMESTAMPTZ USING deleted_at AT TIME ZONE 'UTC';

-- =====================================================
-- TABEL: kupon
-- =====================================================

-- Hapus default dari created_at dan updated_at (tipe sudah TIMESTAMPTZ sejak migration awal)
ALTER TABLE kupon
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN updated_at DROP DEFAULT;
