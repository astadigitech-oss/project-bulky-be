-- migrations/000091_seed_warehouse.up.sql
-- Seed data untuk tabel warehouse (single location - sama dengan informasi_pickup)

INSERT INTO warehouse (nama, slug, alamat, kota, kode_pos, telepon, is_active, created_at, updated_at)
VALUES (
    'Gudang Cilodong',
    'gudang-cilodong',
    'Jl. Cilodong Raya No.89, Cilodong, Kec. Cilodong, Kota Depok, Jawa Barat 16414',
    'Depok',
    '16414',
    '62811833164',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (slug) DO NOTHING;

COMMENT ON TABLE warehouse IS 'Master data warehouse/gudang untuk manajemen stok produk';
