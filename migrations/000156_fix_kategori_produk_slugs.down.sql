-- migrations/000156_fix_kategori_produk_slugs.down.sql
-- Revert semua perubahan ke nilai semula (dari seed 000144)

-- Revert Bagian 1: slug_id

-- ATK
UPDATE kategori_produk
SET slug_id    = 'stationery',
    updated_at = NOW()
WHERE id = 'a1d27bec-e7f3-45e2-8b33-253998128d89'
  AND deleted_at IS NULL;

-- Mainan
UPDATE kategori_produk
SET slug_id    = 'toys',
    updated_at = NOW()
WHERE id = '9d45adfa-1ddb-4191-9a13-820dc2398667'
  AND deleted_at IS NULL;

-- Lainnya
UPDATE kategori_produk
SET slug_id    = 'uncategorized',
    updated_at = NOW()
WHERE id = '9cc2e714-28ca-4e45-a09a-17aedea3c95e'
  AND deleted_at IS NULL;

-- Fashion: tidak ada perubahan di up, tidak ada yang perlu direvert

-- Revert Bagian 2: nama_en dan slug_en

UPDATE kategori_produk SET nama_en = 'Aksesoris',        slug_en = 'aksesoris',        updated_at = NOW() WHERE id = '9cc2dd3c-49e8-4747-a957-0743e11e1f03' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Alat Rumah Tangga', slug_en = 'alat-rumah-tangga', updated_at = NOW() WHERE id = '9cc2dd0f-2acc-4867-a629-28b6e6ef8d60' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Buku',             slug_en = 'buku',             updated_at = NOW() WHERE id = '9cc2dd42-602f-429d-a723-12fd7fc68516' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Elektronik',       slug_en = 'elektronik',       updated_at = NOW() WHERE id = '9cc0e340-047e-471d-9a94-ace606b2c371' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Fashion & Aksesoris', slug_en = 'fashion-aksesoris', updated_at = NOW() WHERE id = '9cc2dd76-8c15-4087-810d-43033843d7a9' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Fashion & Tas',    slug_en = 'fashion-tas',      updated_at = NOW() WHERE id = '9cc2dd5a-5f45-4e81-b2c5-daf00ef9b4af' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Ibu & Anak',       slug_en = 'ibu-anak',         updated_at = NOW() WHERE id = '9cc2dced-553f-44e3-9158-3deb15e3a206' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Kosmetik',         slug_en = 'kosmetik',         updated_at = NOW() WHERE id = '9cc2dcf6-9f70-4669-9d0f-42fad62f1d76' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Kulkas',           slug_en = 'kulkas',           updated_at = NOW() WHERE id = '9cc2dd7c-ef99-42d8-8c54-668be2050fd1' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Lainnya',          slug_en = 'uncategorized',    updated_at = NOW() WHERE id = '9cc2e714-28ca-4e45-a09a-17aedea3c95e' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Mesin Cuci',       slug_en = 'mesin-cuci',       updated_at = NOW() WHERE id = '9cc2dd84-5037-4462-af48-0bb940f2509d' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Otomotif',         slug_en = 'otomotif',         updated_at = NOW() WHERE id = '9cc2dd00-5470-48fb-9c71-bc3654daaf13' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Sepatu',           slug_en = 'sepatu',           updated_at = NOW() WHERE id = '9cc2dd33-ed38-4172-b190-ce277301b7d8' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Tas',              slug_en = 'tas',              updated_at = NOW() WHERE id = '9cc2dd49-6eb5-4fa3-b7f0-0a8f7eece1a4' AND deleted_at IS NULL;
UPDATE kategori_produk SET nama_en = 'Unggulan',         slug_en = 'unggulan',         updated_at = NOW() WHERE id = '9d3fa7e2-5b04-4b82-96f9-6c1c05945546' AND deleted_at IS NULL;
