-- 000127_create_produk_merek.up.sql
-- Create pivot table for many-to-many relationship between produk and merek_produk

-- 1. Create pivot table
CREATE TABLE produk_merek (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    produk_id UUID NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    merek_id UUID NOT NULL REFERENCES merek_produk(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Unique constraint: no duplicate merek per produk
    CONSTRAINT uq_produk_merek UNIQUE (produk_id, merek_id)
);

-- Create indexes for better query performance
CREATE INDEX idx_produk_merek_produk_id ON produk_merek(produk_id);
CREATE INDEX idx_produk_merek_merek_id ON produk_merek(merek_id);

-- 2. Migrate existing data from produk.merek_id to produk_merek
INSERT INTO produk_merek (produk_id, merek_id)
SELECT id, merek_id 
FROM produk 
WHERE merek_id IS NOT NULL AND deleted_at IS NULL;

-- 3. Add comment for documentation
COMMENT ON TABLE produk_merek IS 'Pivot table for many-to-many relationship between produk and merek_produk';
COMMENT ON COLUMN produk_merek.produk_id IS 'Foreign key to produk table';
COMMENT ON COLUMN produk_merek.merek_id IS 'Foreign key to merek_produk table';
