-- Seed data untuk tabel kondisi_produk
INSERT INTO kondisi_produk (id, nama_id, nama_en, slug, urutan, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Baru', 'New', 'baru', 1, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Bekas Seperti Baru', 'Like New', 'bekas-seperti-baru', 2, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Bekas Baik', 'Good Condition', 'bekas-baik', 3, true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'),
    (uuid_generate_v4(), 'Bekas Cukup Baik', 'Fair Condition', 'bekas-cukup-baik', 4, true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second'),
    (uuid_generate_v4(), 'Rusak', 'Damaged', 'rusak', 5, true, NOW() + INTERVAL '4 second', NOW() + INTERVAL '4 second')
ON CONFLICT (slug) DO NOTHING;
