-- migrations/000137_fix_timestamp_to_timestamptz.down.sql
--
-- Rollback: kembalikan ke TIMESTAMP tanpa timezone + restore DEFAULT NOW()

-- =====================================================
-- TABEL: pesanan
-- =====================================================

ALTER TABLE pesanan
    ALTER COLUMN expired_at    TYPE TIMESTAMP USING expired_at AT TIME ZONE 'UTC',
    ALTER COLUMN paid_at       TYPE TIMESTAMP USING paid_at AT TIME ZONE 'UTC',
    ALTER COLUMN processed_at  TYPE TIMESTAMP USING processed_at AT TIME ZONE 'UTC',
    ALTER COLUMN ready_at      TYPE TIMESTAMP USING ready_at AT TIME ZONE 'UTC',
    ALTER COLUMN shipped_at    TYPE TIMESTAMP USING shipped_at AT TIME ZONE 'UTC',
    ALTER COLUMN completed_at  TYPE TIMESTAMP USING completed_at AT TIME ZONE 'UTC',
    ALTER COLUMN cancelled_at  TYPE TIMESTAMP USING cancelled_at AT TIME ZONE 'UTC';

ALTER TABLE pesanan
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'UTC';

-- =====================================================
-- TABEL: ulasan
-- =====================================================

ALTER TABLE ulasan
    ALTER COLUMN approved_at TYPE TIMESTAMP USING approved_at AT TIME ZONE 'UTC';

ALTER TABLE ulasan
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'UTC';

-- =====================================================
-- TABEL: kupon
-- =====================================================

ALTER TABLE kupon
    ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
    ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
