ALTER TABLE auction_batches
    ADD COLUMN IF NOT EXISTS slug_id VARCHAR(320),
    ADD COLUMN IF NOT EXISTS slug_en VARCHAR(320);

UPDATE auction_batches
SET slug_id = CONCAT(
    COALESCE(
        NULLIF(regexp_replace(lower(nama_id), '[^a-z0-9]+', '-', 'g'), ''),
        'auction'
    ),
    '-',
    substring(id::text FROM 1 FOR 8)
)
WHERE slug_id IS NULL OR slug_id = '';

ALTER TABLE auction_batches
    ALTER COLUMN slug_id SET NOT NULL;

UPDATE auction_batches
SET slug_en = CONCAT(
    COALESCE(
        NULLIF(regexp_replace(lower(COALESCE(NULLIF(nama_en, ''), nama_id)), '[^a-z0-9]+', '-', 'g'), ''),
        'auction'
    ),
    '-',
    substring(id::text FROM 1 FOR 8)
)
WHERE slug_en IS NULL OR slug_en = '';

ALTER TABLE auction_batches
    ALTER COLUMN slug_en SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_auction_batches_slug_id
    ON auction_batches(slug_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_auction_batches_slug_en
    ON auction_batches(slug_en);
