-- migrations/000089_update_informasi_pickup_coordinates.up.sql

-- Update koordinat GPS dan Google Maps URL untuk informasi pickup
UPDATE informasi_pickup 
SET 
    latitude = -6.4364499,
    longitude = 106.8479814,
    google_maps_url = 'https://www.google.com/maps/place/Jl.+Cilodong+Raya+No.89,+Cilodong,+Kec.+Cilodong,+Kota+Depok,+Jawa+Barat+16414/@-6.4365015,106.84757,20.89z/data=!4m6!3m5!1s0x2e69ea44081a3b35:0xae0cff609af501a6!8m2!3d-6.4364499!4d106.8479814!16s%2Fg%2F11fmgb30c7?entry=ttu&g_ep=EgoyMDI2MDEyMC4wIKXMDSoKLDEwMDc5MjA2N0gBUAM%3D'
WHERE id = (SELECT id FROM informasi_pickup LIMIT 1);

COMMENT ON COLUMN informasi_pickup.latitude IS 'Koordinat latitude lokasi warehouse';
COMMENT ON COLUMN informasi_pickup.longitude IS 'Koordinat longitude lokasi warehouse';
COMMENT ON COLUMN informasi_pickup.google_maps_url IS 'URL Google Maps untuk navigasi ke warehouse';
