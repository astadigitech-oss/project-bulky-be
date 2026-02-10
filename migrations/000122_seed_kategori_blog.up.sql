-- =====================================================
-- SEED: kategori_blog
-- =====================================================

INSERT INTO kategori_blog (id, nama_id, nama_en, slug, is_active, urutan) VALUES
(uuid_generate_v4(), 'Tips & Trick', 'Tips & Tricks', 'tips-trick', true, 1),
(uuid_generate_v4(), 'Berita', 'News', 'berita', true, 2),
(uuid_generate_v4(), 'Produk', 'Products', 'produk', true, 3),
(uuid_generate_v4(), 'Tutorial', 'Tutorial', 'tutorial', true, 4),
(uuid_generate_v4(), 'Promo', 'Promotion', 'promo', true, 5),
(uuid_generate_v4(), 'Bisnis', 'Business', 'bisnis', true, 6),
(uuid_generate_v4(), 'Lifestyle', 'Lifestyle', 'lifestyle', true, 7);
