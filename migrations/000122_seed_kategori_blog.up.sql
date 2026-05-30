-- migrations/000122_seed_kategori_blog.up.sql
-- 8 kategori blog (bilingual)

INSERT INTO kategori_blog (id, nama_id, nama_en, slug, is_active, urutan) VALUES
(uuid_generate_v4(), 'Tips & Trik',              'Tips & Tricks',            'tips-trik',           true, 1),
(uuid_generate_v4(), 'Panduan Pembelian',         'Buying Guide',             'panduan-pembelian',   true, 2),
(uuid_generate_v4(), 'Logistik & Pengiriman',     'Logistics & Shipping',     'logistik-pengiriman', true, 3),
(uuid_generate_v4(), 'Pengetahuan Produk',        'Product Knowledge',        'pengetahuan-produk',  true, 4),
(uuid_generate_v4(), 'Bisnis Grosir',             'Wholesale Business',       'bisnis-grosir',       true, 5),
(uuid_generate_v4(), 'Berita & Update',           'News & Updates',           'berita-update',       true, 6),
(uuid_generate_v4(), 'Regulasi & Kebijakan',      'Regulations & Policy',     'regulasi-kebijakan',  true, 7),
(uuid_generate_v4(), 'Studi Kasus',               'Case Study',               'studi-kasus',         true, 8)
ON CONFLICT (slug) DO UPDATE SET
    nama_id    = EXCLUDED.nama_id,
    nama_en    = EXCLUDED.nama_en,
    is_active  = EXCLUDED.is_active,
    urutan     = EXCLUDED.urutan;
