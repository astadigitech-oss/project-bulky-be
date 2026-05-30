-- migrations/000124_seed_kategori_video.up.sql
-- 6 kategori video (bilingual)

INSERT INTO kategori_video (id, nama_id, nama_en, slug, is_active, urutan) VALUES
(uuid_generate_v4(), 'Tutorial Platform',     'Platform Tutorial',      'tutorial-platform',    true, 1),
(uuid_generate_v4(), 'Review Produk',         'Product Review',         'review-produk',        true, 2),
(uuid_generate_v4(), 'Panduan Pengiriman',    'Shipping Guide',         'panduan-pengiriman',   true, 3),
(uuid_generate_v4(), 'FAQ & Tanya Jawab',     'FAQ',                    'faq',                  true, 4),
(uuid_generate_v4(), 'Testimoni Pembeli',     'Buyer Testimonials',     'testimoni',            true, 5),
(uuid_generate_v4(), 'Behind the Scenes',     'Behind the Scenes',      'behind-the-scenes',    true, 6)
ON CONFLICT (slug) DO UPDATE SET
    nama_id    = EXCLUDED.nama_id,
    nama_en    = EXCLUDED.nama_en,
    is_active  = EXCLUDED.is_active,
    urutan     = EXCLUDED.urutan;
