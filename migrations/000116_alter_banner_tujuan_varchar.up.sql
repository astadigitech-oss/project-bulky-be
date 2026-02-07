-- =====================================================
-- REVISI: banner_event_promo.tujuan
-- =====================================================
-- Mengubah tujuan dari JSONB ke VARCHAR (comma-separated)
-- Format: "uuid-1,uuid-2,uuid-3"
-- =====================================================

-- Step 1: Drop JSONB index if exists
DROP INDEX IF EXISTS idx_banner_tujuan;

-- Step 2: Create temporary function to convert JSONB to string
CREATE OR REPLACE FUNCTION jsonb_array_to_string(j jsonb) RETURNS text AS $$
BEGIN
    IF j IS NULL THEN
        RETURN NULL;
    END IF;
    
    IF jsonb_typeof(j) = 'array' THEN
        RETURN (
            SELECT string_agg(elem->>'id', ',')
            FROM jsonb_array_elements(j) AS elem
        );
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Step 3: Alter column type using the function
ALTER TABLE banner_event_promo 
ALTER COLUMN tujuan TYPE VARCHAR(1000) USING jsonb_array_to_string(tujuan);

-- Step 4: Drop temporary function
DROP FUNCTION jsonb_array_to_string(jsonb);

-- Step 5: Update comment
COMMENT ON COLUMN banner_event_promo.tujuan IS 'Comma-separated kategori_produk IDs. NULL = banner tanpa redirect. Example: "uuid-1,uuid-2,uuid-3"';
