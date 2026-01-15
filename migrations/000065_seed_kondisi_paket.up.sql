-- Seed data untuk tabel kondisi_paket
INSERT INTO kondisi_paket (id, nama_id, nama_en, slug, urutan, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Baik', 'Good', 'baik', 1, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Rusak Ringan', 'Slightly Damaged', 'rusak-ringan', 2, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Rusak Sedang', 'Moderately Damaged', 'rusak-sedang', 3, true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'),
    (uuid_generate_v4(), 'Rusak Berat', 'Heavily Damaged', 'rusak-berat', 4, true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second')
ON CONFLICT (slug) DO NOTHING;
