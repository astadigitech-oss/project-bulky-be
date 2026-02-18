-- Alter tanggal_kedaluarsa from DATE to TIMESTAMPTZ
-- Existing DATE values (stored as date-only, interpreted at UTC midnight) are cast to
-- TIMESTAMPTZ at UTC midnight so no data is lost.
ALTER TABLE kupon
    ALTER COLUMN tanggal_kedaluarsa TYPE TIMESTAMPTZ
        USING tanggal_kedaluarsa::TIMESTAMPTZ;

-- Update comment to reflect new type
COMMENT ON COLUMN kupon.tanggal_kedaluarsa IS 'Waktu berakhirnya kupon (timestamptz, UTC+0)';
