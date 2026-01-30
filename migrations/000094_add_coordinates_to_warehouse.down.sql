-- Remove latitude and longitude columns from warehouse table
ALTER TABLE warehouse
DROP COLUMN IF EXISTS latitude,
DROP COLUMN IF EXISTS longitude;
