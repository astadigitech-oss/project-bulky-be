-- Hapus penanda threshold manual.
ALTER TABLE deliveree_vehicle_type
    DROP COLUMN IF EXISTS threshold_is_manual;
