-- migrations/000089_update_informasi_pickup_coordinates.down.sql

-- Rollback: hapus koordinat dan Google Maps URL
UPDATE informasi_pickup 
SET 
    latitude = NULL,
    longitude = NULL,
    google_maps_url = NULL
WHERE id = (SELECT id FROM informasi_pickup LIMIT 1);
