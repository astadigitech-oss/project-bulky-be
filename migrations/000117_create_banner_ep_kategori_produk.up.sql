-- =====================================================
-- TABEL: banner_ep_kategori_produk (Pivot Table)
-- =====================================================
-- Relasi one-to-many antara banner_event_promo dan kategori_produk
-- Menggantikan field tujuan (VARCHAR) di banner_event_promo
-- =====================================================

-- Step 1: Create pivot table
CREATE TABLE banner_ep_kategori_produk (
    banner_id UUID NOT NULL,
    kategori_produk_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Composite Primary Key
    PRIMARY KEY (banner_id, kategori_produk_id),
    
    -- Foreign Keys
    CONSTRAINT fk_banner_ep_banner
        FOREIGN KEY (banner_id) 
        REFERENCES banner_event_promo(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_banner_ep_kategori
        FOREIGN KEY (kategori_produk_id) 
        REFERENCES kategori_produk(id)
        ON DELETE CASCADE
);

-- Step 2: Indexes for query performance
CREATE INDEX idx_banner_ep_banner_id ON banner_ep_kategori_produk(banner_id);
CREATE INDEX idx_banner_ep_kategori_id ON banner_ep_kategori_produk(kategori_produk_id);

-- Step 3: Migrate existing data (if tujuan column exists with comma-separated values)
-- Skip if tujuan column doesn't exist or is empty
DO $$
DECLARE
    banner_record RECORD;
    kategori_id TEXT;
    kategori_ids TEXT[];
BEGIN
    -- Check if tujuan column exists
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'banner_event_promo' AND column_name = 'tujuan'
    ) THEN
        FOR banner_record IN 
            SELECT id, tujuan FROM banner_event_promo 
            WHERE tujuan IS NOT NULL AND tujuan != '' AND deleted_at IS NULL
        LOOP
            -- Split comma-separated string
            kategori_ids := string_to_array(banner_record.tujuan, ',');
            
            -- Insert each kategori
            FOREACH kategori_id IN ARRAY kategori_ids
            LOOP
                BEGIN
                    INSERT INTO banner_ep_kategori_produk (banner_id, kategori_produk_id)
                    VALUES (banner_record.id, TRIM(kategori_id)::UUID)
                    ON CONFLICT DO NOTHING;
                EXCEPTION WHEN OTHERS THEN
                    -- Skip invalid UUIDs
                    NULL;
                END;
            END LOOP;
        END LOOP;
    END IF;
END $$;

-- Step 4: Drop tujuan column from banner_event_promo
ALTER TABLE banner_event_promo DROP COLUMN IF EXISTS tujuan;

-- Step 5: Comments
COMMENT ON TABLE banner_ep_kategori_produk IS 'Pivot table untuk relasi banner_event_promo ke kategori_produk';
COMMENT ON COLUMN banner_ep_kategori_produk.banner_id IS 'FK ke banner_event_promo.id';
COMMENT ON COLUMN banner_ep_kategori_produk.kategori_produk_id IS 'FK ke kategori_produk.id';
