-- migrations/000123_seed_label_blog.up.sql
-- 12 label/tag blog (bilingual)

INSERT INTO label_blog (id, nama_id, nama_en, slug, urutan) VALUES
(uuid_generate_v4(), 'Pemula',              'Beginner',           'pemula',           1),
(uuid_generate_v4(), 'Lanjutan',            'Advanced',           'lanjutan',         2),
(uuid_generate_v4(), 'Paletbox',            'Paletbox',           'paletbox',         3),
(uuid_generate_v4(), 'Pengiriman Darat',    'Land Freight',       'pengiriman-darat', 4),
(uuid_generate_v4(), 'Pengiriman Laut',     'Sea Cargo',          'pengiriman-laut',  5),
(uuid_generate_v4(), 'Pembayaran',          'Payment',            'pembayaran',       6),
(uuid_generate_v4(), 'Kemasan & Packing',   'Packaging',          'kemasan-packing',  7),
(uuid_generate_v4(), 'Efisiensi Biaya',     'Cost Efficiency',    'efisiensi-biaya',  8),
(uuid_generate_v4(), 'Pergudangan',         'Warehousing',        'pergudangan',      9),
(uuid_generate_v4(), 'B2B',                 'B2B',                'b2b',              10),
(uuid_generate_v4(), 'Impor & Ekspor',      'Import & Export',    'impor-ekspor',     11),
(uuid_generate_v4(), 'Regulasi Bea Cukai',  'Customs Regulation', 'bea-cukai',        12)
ON CONFLICT (slug) DO UPDATE SET
    nama_id   = EXCLUDED.nama_id,
    nama_en   = EXCLUDED.nama_en,
    urutan    = EXCLUDED.urutan;
