-- Rollback: Convert VARCHAR back to JSONB
-- Note: This will lose slug data, only preserves IDs

-- Step 1: Create temporary function to convert string to JSONB
CREATE OR REPLACE FUNCTION string_to_jsonb_array(s text) RETURNS jsonb AS $$
BEGIN
    IF s IS NULL OR s = '' THEN
        RETURN NULL;
    END IF;
    
    RETURN (
        SELECT jsonb_agg(jsonb_build_object('id', trim(id), 'slug', ''))
        FROM unnest(string_to_array(s, ',')) AS id
        WHERE trim(id) != ''
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Step 2: Alter column type using the function
ALTER TABLE banner_event_promo 
ALTER COLUMN tujuan TYPE JSONB USING string_to_jsonb_array(tujuan);

-- Step 3: Drop temporary function
DROP FUNCTION string_to_jsonb_array(text);

-- Step 4: Recreate index
CREATE INDEX idx_banner_tujuan ON banner_event_promo 
USING GIN (tujuan) WHERE deleted_at IS NULL;
