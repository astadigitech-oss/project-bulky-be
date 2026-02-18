-- Revert tanggal_kedaluarsa from TIMESTAMPTZ back to DATE
-- Truncates the time portion; only the UTC date is kept.
ALTER TABLE kupon
    ALTER COLUMN tanggal_kedaluarsa TYPE DATE
        USING tanggal_kedaluarsa::DATE;

COMMENT ON COLUMN kupon.tanggal_kedaluarsa IS 'Tanggal berakhirnya kupon';
