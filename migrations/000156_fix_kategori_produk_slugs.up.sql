-- migrations/000156_fix_kategori_produk_slugs.up.sql
-- =====================================================
-- FIX: Kategori Produk - Koreksi slug_id, nama_en, slug_en
-- =====================================================
-- Masalah yang diperbaiki:
-- 1. slug_id tidak sesuai nama_id: ATK, Mainan, Lainnya
-- 2. slug_id 'fashion-1' (artefak dedup) → 'fashion' untuk nama_id 'Fashion'
-- 3. nama_en masih dalam Bahasa Indonesia → diterjemahkan ke Inggris
-- 4. slug_en disesuaikan dengan nama_en baru
-- =====================================================
-- AMAN: semua UPDATE menggunakan PRIMARY KEY (UUID), tidak ada DDL,
--       tidak ada DROP/TRUNCATE, menggunakan id spesifik.
-- =====================================================

-- -----------------------------------------------
-- BAGIAN 1: Fix slug_id agar sesuai nama_id
-- -----------------------------------------------

-- ATK: slug_id 'stationery' → 'atk'
-- (slug_en tetap 'stationery' karena sudah sesuai nama_en = 'Stationery')
UPDATE kategori_produk
SET slug_id    = 'atk',
    updated_at = NOW()
WHERE id = 'a1d27bec-e7f3-45e2-8b33-253998128d89'
  AND deleted_at IS NULL;

-- Mainan: slug_id 'toys' → 'mainan'
-- (slug_en tetap 'toys' karena sudah sesuai nama_en = 'Toys')
UPDATE kategori_produk
SET slug_id    = 'mainan',
    updated_at = NOW()
WHERE id = '9d45adfa-1ddb-4191-9a13-820dc2398667'
  AND deleted_at IS NULL;

-- Lainnya: slug_id 'uncategorized' → 'lainnya'
UPDATE kategori_produk
SET slug_id    = 'lainnya',
    updated_at = NOW()
WHERE id = '9cc2e714-28ca-4e45-a09a-17aedea3c95e'
  AND deleted_at IS NULL;

-- Fashion: slug 'fashion-1' TIDAK diubah karena slug='fashion' sudah dipakai row lain di production

-- -----------------------------------------------
-- BAGIAN 2: Fix nama_en dan slug_en
-- -----------------------------------------------

-- Aksesoris → Accessories
UPDATE kategori_produk
SET nama_en    = 'Accessories',
    slug_en    = 'accessories',
    updated_at = NOW()
WHERE id = '9cc2dd3c-49e8-4747-a957-0743e11e1f03'
  AND deleted_at IS NULL;

-- Alat Rumah Tangga → Household Appliances
UPDATE kategori_produk
SET nama_en    = 'Household Appliances',
    slug_en    = 'household-appliances',
    updated_at = NOW()
WHERE id = '9cc2dd0f-2acc-4867-a629-28b6e6ef8d60'
  AND deleted_at IS NULL;

-- Buku → Books
UPDATE kategori_produk
SET nama_en    = 'Books',
    slug_en    = 'books',
    updated_at = NOW()
WHERE id = '9cc2dd42-602f-429d-a723-12fd7fc68516'
  AND deleted_at IS NULL;

-- Elektronik → Electronics
UPDATE kategori_produk
SET nama_en    = 'Electronics',
    slug_en    = 'electronics',
    updated_at = NOW()
WHERE id = '9cc0e340-047e-471d-9a94-ace606b2c371'
  AND deleted_at IS NULL;

-- Fashion & Aksesoris → Fashion & Accessories
UPDATE kategori_produk
SET nama_en    = 'Fashion & Accessories',
    slug_en    = 'fashion-accessories',
    updated_at = NOW()
WHERE id = '9cc2dd76-8c15-4087-810d-43033843d7a9'
  AND deleted_at IS NULL;

-- Fashion & Tas → Fashion & Bags
UPDATE kategori_produk
SET nama_en    = 'Fashion & Bags',
    slug_en    = 'fashion-bags',
    updated_at = NOW()
WHERE id = '9cc2dd5a-5f45-4e81-b2c5-daf00ef9b4af'
  AND deleted_at IS NULL;

-- Ibu & Anak → Mother & Baby
UPDATE kategori_produk
SET nama_en    = 'Mother & Baby',
    slug_en    = 'mother-and-baby',
    updated_at = NOW()
WHERE id = '9cc2dced-553f-44e3-9158-3deb15e3a206'
  AND deleted_at IS NULL;

-- Kosmetik → Cosmetics
UPDATE kategori_produk
SET nama_en    = 'Cosmetics',
    slug_en    = 'cosmetics',
    updated_at = NOW()
WHERE id = '9cc2dcf6-9f70-4669-9d0f-42fad62f1d76'
  AND deleted_at IS NULL;

-- Kulkas → Refrigerator
UPDATE kategori_produk
SET nama_en    = 'Refrigerator',
    slug_en    = 'refrigerator',
    updated_at = NOW()
WHERE id = '9cc2dd7c-ef99-42d8-8c54-668be2050fd1'
  AND deleted_at IS NULL;

-- Lainnya → Others
UPDATE kategori_produk
SET nama_en    = 'Others',
    slug_en    = 'others',
    updated_at = NOW()
WHERE id = '9cc2e714-28ca-4e45-a09a-17aedea3c95e'
  AND deleted_at IS NULL;

-- Mesin Cuci → Washing Machine
UPDATE kategori_produk
SET nama_en    = 'Washing Machine',
    slug_en    = 'washing-machine',
    updated_at = NOW()
WHERE id = '9cc2dd84-5037-4462-af48-0bb940f2509d'
  AND deleted_at IS NULL;

-- Otomotif → Automotive
UPDATE kategori_produk
SET nama_en    = 'Automotive',
    slug_en    = 'automotive',
    updated_at = NOW()
WHERE id = '9cc2dd00-5470-48fb-9c71-bc3654daaf13'
  AND deleted_at IS NULL;

-- Sepatu → Shoes
UPDATE kategori_produk
SET nama_en    = 'Shoes',
    slug_en    = 'shoes',
    updated_at = NOW()
WHERE id = '9cc2dd33-ed38-4172-b190-ce277301b7d8'
  AND deleted_at IS NULL;

-- Tas → Bags
UPDATE kategori_produk
SET nama_en    = 'Bags',
    slug_en    = 'bags',
    updated_at = NOW()
WHERE id = '9cc2dd49-6eb5-4fa3-b7f0-0a8f7eece1a4'
  AND deleted_at IS NULL;

-- Unggulan → Featured
UPDATE kategori_produk
SET nama_en    = 'Featured',
    slug_en    = 'featured',
    updated_at = NOW()
WHERE id = '9d3fa7e2-5b04-4b82-96f9-6c1c05945546'
  AND deleted_at IS NULL;
