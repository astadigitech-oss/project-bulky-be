-- Add latitude and longitude columns to warehouse table
ALTER TABLE warehouse
ADD COLUMN latitude DECIMAL(10,8) NULL,
ADD COLUMN longitude DECIMAL(11,8) NULL;

-- Update existing warehouse with coordinates from informasi_pickup
UPDATE warehouse w
SET 
    latitude = ip.latitude,
    longitude = ip.longitude
FROM informasi_pickup ip
WHERE w.slug = 'gudang-cilodong';

-- Add comments
COMMENT ON COLUMN warehouse.latitude IS 'Koordinat latitude lokasi warehouse';
COMMENT ON COLUMN warehouse.longitude IS 'Koordinat longitude lokasi warehouse';
