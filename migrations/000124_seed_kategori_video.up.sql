-- =====================================================
-- SEED: kategori_video
-- =====================================================

INSERT INTO kategori_video (id, nama_id, nama_en, slug, is_active, urutan) VALUES
(uuid_generate_v4(), 'Tutorial', 'Tutorial', 'tutorial', true, 1),
(uuid_generate_v4(), 'Review Produk', 'Product Review', 'review-produk', true, 2),
(uuid_generate_v4(), 'Tips Usaha', 'Business Tips', 'tips-usaha', true, 3),
(uuid_generate_v4(), 'Promo & Event', 'Promo & Event', 'promo-event', true, 4),
(uuid_generate_v4(), 'Behind The Scene', 'Behind The Scene', 'behind-the-scene', true, 5),
(uuid_generate_v4(), 'Testimoni', 'Testimonial', 'testimoni', true, 6),
(uuid_generate_v4(), 'Tips Belanja', 'Shopping Tips', 'tips-belanja', true, 7);
