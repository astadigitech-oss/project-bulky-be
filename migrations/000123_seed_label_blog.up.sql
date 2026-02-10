-- =====================================================
-- SEED: label_blog
-- =====================================================

INSERT INTO label_blog (id, nama_id, nama_en, slug, urutan) VALUES
(uuid_generate_v4(), 'Trending', 'Trending', 'trending', 1),
(uuid_generate_v4(), 'Popular', 'Popular', 'popular', 2),
(uuid_generate_v4(), 'Rekomendasi', 'Recommended', 'rekomendasi', 3),
(uuid_generate_v4(), 'Terbaru', 'Latest', 'terbaru', 4),
(uuid_generate_v4(), 'Pilihan Editor', 'Editor''s Choice', 'pilihan-editor', 5),
(uuid_generate_v4(), 'Grosir', 'Wholesale', 'grosir', 6),
(uuid_generate_v4(), 'Hemat', 'Save Money', 'hemat', 7),
(uuid_generate_v4(), 'Peluang Usaha', 'Business Opportunity', 'peluang-usaha', 8);
