-- Seed data untuk tabel kondisi_produk
INSERT INTO kondisi_produk (id, nama_id, nama_en, slug, urutan, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Mediocre 60-80%', 'Mediocre 60-80%', 'mediocre-60-80', 1, true, NOW(), NOW()),
    (uuid_generate_v4(), '100% Brandnew', '100% Brandnew', 'brandnew', 2, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'New 90%-95%', 'New 90%-95%', '90-100', 3, true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second')
ON CONFLICT (slug) DO NOTHING;
