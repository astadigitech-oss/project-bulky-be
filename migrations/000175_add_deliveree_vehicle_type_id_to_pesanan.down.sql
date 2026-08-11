-- Hapus kolom deliveree_vehicle_type_id.
ALTER TABLE pesanan
    DROP COLUMN IF EXISTS deliveree_vehicle_type_id;
