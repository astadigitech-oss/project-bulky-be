-- Rollback: kembalikan kolom ke timestamp without time zone (wall-clock UTC).
-- Catatan: nilai yang tersimpan sekarang adalah UTC, sehingga untuk mengembalikan
-- ke semantik lama (WIB-stored) perlu konversi balik AT TIME ZONE 'UTC'
-- (memotong offset) supaya wall-clock-nya = UTC asli, BUKAN AT TIME ZONE 'Asia/Jakarta'.
-- ⚠️ Rollback TIDAK mengembalikan data lama ke WIB-stored — hanya memotong offset UTC.

ALTER TABLE banner_event_promo
    ALTER COLUMN tanggal_mulai TYPE TIMESTAMP WITHOUT TIME ZONE USING tanggal_mulai AT TIME ZONE 'UTC',
    ALTER COLUMN tanggal_selesai TYPE TIMESTAMP WITHOUT TIME ZONE USING tanggal_selesai AT TIME ZONE 'UTC',
    ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITHOUT TIME ZONE USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN deleted_at TYPE TIMESTAMP WITHOUT TIME ZONE USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE admin
    ALTER COLUMN last_login_at TYPE TIMESTAMP WITHOUT TIME ZONE USING last_login_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITHOUT TIME ZONE USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN deleted_at TYPE TIMESTAMP WITHOUT TIME ZONE USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE activity_log
    ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE USING created_at AT TIME ZONE 'UTC';

ALTER TABLE formulir_partai_besar_submission
    ALTER COLUMN email_sent_at TYPE TIMESTAMP WITHOUT TIME ZONE USING email_sent_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE USING created_at AT TIME ZONE 'UTC';
