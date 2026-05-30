-- Seed data untuk tabel kondisi_paket
INSERT INTO kondisi_paket (id, nama_id, nama_en, slug, urutan, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Sedang', 'Sedang', 'sedang', 1, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Bagus', 'Bagus', 'bagus', 2, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second')
ON CONFLICT (slug) DO NOTHING;
