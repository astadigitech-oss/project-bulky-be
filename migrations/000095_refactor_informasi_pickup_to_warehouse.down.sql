-- Rollback: Restore informasi_pickup table and revert jadwal_gudang

-- Step 1: Recreate informasi_pickup table
CREATE TABLE IF NOT EXISTS informasi_pickup (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alamat TEXT NOT NULL,
    jam_operasional VARCHAR(100) NOT NULL,
    nomor_whatsapp VARCHAR(20) NOT NULL,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    google_maps_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Step 2: Copy data from warehouse to informasi_pickup (use telepon for nomor_whatsapp)
INSERT INTO informasi_pickup (jam_operasional, nomor_whatsapp, alamat, latitude, longitude)
SELECT 
    COALESCE(jam_operasional, 'Senin - Sabtu, 09.00 - 18.00 WIB'),
    COALESCE(telepon, '62811833164'),
    COALESCE(alamat, 'Alamat belum diisi'),
    latitude,
    longitude
FROM warehouse
WHERE is_active = true
LIMIT 1;

-- Step 3: Add informasi_pickup_id back to jadwal_gudang
ALTER TABLE jadwal_gudang
ADD COLUMN IF NOT EXISTS informasi_pickup_id UUID;

-- Step 4: Set informasi_pickup_id
UPDATE jadwal_gudang jg
SET informasi_pickup_id = (SELECT id FROM informasi_pickup LIMIT 1)
WHERE informasi_pickup_id IS NULL;

-- Step 5: Make informasi_pickup_id NOT NULL
ALTER TABLE jadwal_gudang
ALTER COLUMN informasi_pickup_id SET NOT NULL;

-- Step 6: Add foreign key
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jadwal_gudang_informasi_pickup'
    ) THEN
        ALTER TABLE jadwal_gudang
        ADD CONSTRAINT fk_jadwal_gudang_informasi_pickup
        FOREIGN KEY (informasi_pickup_id) REFERENCES informasi_pickup(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Step 7: Drop warehouse_id
ALTER TABLE jadwal_gudang
DROP CONSTRAINT IF EXISTS fk_jadwal_gudang_warehouse;

ALTER TABLE jadwal_gudang
DROP COLUMN IF EXISTS warehouse_id;

-- Step 8: Remove jam_operasional from warehouse (telepon stays)
ALTER TABLE warehouse
DROP COLUMN IF EXISTS jam_operasional;
