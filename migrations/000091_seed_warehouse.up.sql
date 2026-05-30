-- migrations/000091_seed_warehouse.up.sql
-- Seed data warehouse Cibinong

-- Remove old Cilodong record if exists
DELETE FROM warehouse WHERE slug = 'gudang-cilodong';

INSERT INTO warehouse (nama, slug, alamat, kota, kode_pos, telepon, latitude, longitude, is_active, created_at, updated_at)
VALUES (
    'Warehouse Cibinong',
    'warehouse-cibinong',
    'Jl. Raya Mayor Oking Jaya Atmaja No.62a, Cirimekar, Kec. Cibinong, Kabupaten Bogor, Jawa Barat 16918',
    'Kabupaten Bogor',
    '16918',
    '62811833164',
    -6.46958024,
    106.85984984,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (slug) DO UPDATE SET
    nama       = EXCLUDED.nama,
    alamat     = EXCLUDED.alamat,
    kota       = EXCLUDED.kota,
    kode_pos   = EXCLUDED.kode_pos,
    telepon    = EXCLUDED.telepon,
    latitude   = EXCLUDED.latitude,
    longitude  = EXCLUDED.longitude,
    is_active  = EXCLUDED.is_active,
    updated_at = NOW();

COMMENT ON TABLE warehouse IS 'Master data warehouse/gudang untuk manajemen stok produk';
