-- Refactor: Consolidate informasi_pickup into warehouse and link jadwal_gudang to warehouse

-- Step 1: Add jam_operasional to warehouse (telepon already exists)
ALTER TABLE warehouse
ADD COLUMN IF NOT EXISTS jam_operasional VARCHAR(100);

-- Step 2: Copy data from informasi_pickup to warehouse (if exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'informasi_pickup') THEN
        UPDATE warehouse w
        SET 
            jam_operasional = ip.jam_operasional,
            telepon = ip.nomor_whatsapp
        FROM informasi_pickup ip
        WHERE w.is_active = true;
    END IF;
END $$;

-- Step 3: Add warehouse_id to jadwal_gudang
ALTER TABLE jadwal_gudang
ADD COLUMN IF NOT EXISTS warehouse_id UUID;

-- Step 4: Copy warehouse ID to jadwal_gudang
UPDATE jadwal_gudang jg
SET warehouse_id = (
    SELECT id 
    FROM warehouse 
    WHERE is_active = true 
    ORDER BY created_at ASC 
    LIMIT 1
)
WHERE warehouse_id IS NULL;

-- Step 5: Make warehouse_id NOT NULL
ALTER TABLE jadwal_gudang
ALTER COLUMN warehouse_id SET NOT NULL;

-- Step 6: Add foreign key constraint
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jadwal_gudang_warehouse'
    ) THEN
        ALTER TABLE jadwal_gudang
        ADD CONSTRAINT fk_jadwal_gudang_warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Step 7: Drop old foreign key and column
ALTER TABLE jadwal_gudang
DROP CONSTRAINT IF EXISTS fk_jadwal_gudang_informasi_pickup;

ALTER TABLE jadwal_gudang
DROP COLUMN IF EXISTS informasi_pickup_id;

-- Step 8: Drop informasi_pickup table
DROP TABLE IF EXISTS informasi_pickup CASCADE;

-- Add comments
COMMENT ON COLUMN warehouse.jam_operasional IS 'Jam operasional dalam format text (e.g., "Senin - Sabtu, 09.00 - 18.00 WIB")';
COMMENT ON COLUMN warehouse.telepon IS 'Nomor telepon/WhatsApp untuk kontak (format: 62xxx)';
COMMENT ON COLUMN jadwal_gudang.warehouse_id IS 'Reference to warehouse table';
