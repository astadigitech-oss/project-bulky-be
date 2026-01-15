-- Seed data untuk tabel hero_section
INSERT INTO hero_section (id, nama, gambar_url_id, gambar_url_en, urutan, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Hero Banner 1', 'hero-section/hero-1-id.jpg', 'hero-section/hero-1-en.jpg', 1, false, NOW(), NOW()),
    (uuid_generate_v4(), 'Hero Banner 2', 'hero-section/hero-2-id.jpg', 'hero-section/hero-2-en.jpg', 2, false, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Hero Banner 3', 'hero-section/hero-3-id.jpg', 'hero-section/hero-3-en.jpg', 3, false, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second')
ON CONFLICT DO NOTHING;
