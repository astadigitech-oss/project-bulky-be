-- =====================================================
-- SEED: Kategori Produk (22 kategori)
-- =====================================================
-- icon_url di-null-kan, bisa diisi via admin panel
-- slug_en = slug (sama untuk semua, mayoritas sudah bahasa Inggris/internasional)
-- Khusus: Stationery (nama_id=ATK), Toys (nama_id=Mainan)
-- =====================================================

INSERT INTO kategori_produk (id, nama_id, nama_en, slug, slug_id, slug_en, icon_url, is_active, created_at, updated_at)
VALUES
    ('9cc2dd3c-49e8-4747-a957-0743e11e1f03', 'Aksesoris',        'Aksesoris',        'aksesoris',        'aksesoris',        'aksesoris',          NULL, true, NOW(), NOW()),
    ('9cc2dd0f-2acc-4867-a629-28b6e6ef8d60', 'Alat Rumah Tangga','Alat Rumah Tangga','alat-rumah-tangga','alat-rumah-tangga','alat-rumah-tangga',  NULL, true, NOW(), NOW()),
    ('9cc2dd42-602f-429d-a723-12fd7fc68516', 'Buku',             'Buku',             'buku',             'buku',             'buku',               NULL, true, NOW(), NOW()),
    ('9cc0e340-047e-471d-9a94-ace606b2c371', 'Elektronik',       'Elektronik',       'elektronik',       'elektronik',       'elektronik',         NULL, true, NOW(), NOW()),
    ('9cc2dd4f-cb59-47f8-8455-b4b229ab4e2c', 'Fashion',          'Fashion',          'fashion-1',        'fashion-1',        'fashion-1',          NULL, true, NOW(), NOW()),
    ('9cc2dd76-8c15-4087-810d-43033843d7a9', 'Fashion & Aksesoris','Fashion & Aksesoris','fashion-aksesoris','fashion-aksesoris','fashion-aksesoris',NULL, true, NOW(), NOW()),
    ('9cc2dd5a-5f45-4e81-b2c5-daf00ef9b4af', 'Fashion & Tas',    'Fashion & Tas',    'fashion-tas',      'fashion-tas',      'fashion-tas',        NULL, true, NOW(), NOW()),
    ('9cc2dd15-c450-4200-ad3d-16f4c6352ab7', 'FMCG',             'FMCG',             'fmcg',             'fmcg',             'fmcg',               NULL, true, NOW(), NOW()),
    ('9cc2dced-553f-44e3-9158-3deb15e3a206', 'Ibu & Anak',       'Ibu & Anak',       'ibu-anak',         'ibu-anak',         'ibu-anak',           NULL, true, NOW(), NOW()),
    ('9cc2dcf6-9f70-4669-9d0f-42fad62f1d76', 'Kosmetik',         'Kosmetik',         'kosmetik',         'kosmetik',         'kosmetik',           NULL, true, NOW(), NOW()),
    ('9cc2dd7c-ef99-42d8-8c54-668be2050fd1', 'Kulkas',           'Kulkas',           'kulkas',           'kulkas',           'kulkas',             NULL, true, NOW(), NOW()),
    ('9cc2e714-28ca-4e45-a09a-17aedea3c95e', 'Lainnya',          'Lainnya',          'uncategorized',    'uncategorized',    'uncategorized',      NULL, true, NOW(), NOW()),
    ('9cc2dd84-5037-4462-af48-0bb940f2509d', 'Mesin Cuci',       'Mesin Cuci',       'mesin-cuci',       'mesin-cuci',       'mesin-cuci',         NULL, true, NOW(), NOW()),
    ('9cc2dd00-5470-48fb-9c71-bc3654daaf13', 'Otomotif',         'Otomotif',         'otomotif',         'otomotif',         'otomotif',           NULL, true, NOW(), NOW()),
    ('9cc2dd2d-9050-4367-9620-9e3359b1664b', 'Redknot',          'Redknot',          'redknot',          'redknot',          'redknot',            NULL, true, NOW(), NOW()),
    ('9cc2dd33-ed38-4172-b190-ce277301b7d8', 'Sepatu',           'Sepatu',           'sepatu',           'sepatu',           'sepatu',             NULL, true, NOW(), NOW()),
    ('a1d27bec-e7f3-45e2-8b33-253998128d89', 'ATK',              'Stationery',       'stationery',       'stationery',       'stationery',         NULL, true, NOW(), NOW()),
    ('9cc2dd49-6eb5-4fa3-b7f0-0a8f7eece1a4', 'Tas',              'Tas',              'tas',              'tas',              'tas',                NULL, true, NOW(), NOW()),
    ('9cc2dd1f-7da1-4124-9c3e-171ff73614ec', 'Tools',            'Tools',            'tools',            'tools',            'tools',              NULL, true, NOW(), NOW()),
    ('9d45adfa-1ddb-4191-9a13-820dc2398667', 'Mainan',           'Toys',             'toys',             'toys',             'toys',               NULL, true, NOW(), NOW()),
    ('9cc2dd8a-378b-4771-ae0c-83429104276b', 'TV',               'TV',               'tv',               'tv',               'tv',                 NULL, true, NOW(), NOW()),
    ('9d3fa7e2-5b04-4b82-96f9-6c1c05945546', 'Unggulan',         'Unggulan',         'unggulan',         'unggulan',         'unggulan',           NULL, true, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    nama_id    = EXCLUDED.nama_id,
    nama_en    = EXCLUDED.nama_en,
    slug       = EXCLUDED.slug,
    slug_id    = EXCLUDED.slug_id,
    slug_en    = EXCLUDED.slug_en,
    updated_at = NOW();
