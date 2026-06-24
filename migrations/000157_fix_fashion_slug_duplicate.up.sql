-- migrations/000157_fix_fashion_slug_duplicate.up.sql
-- =====================================================
-- FIX: Duplikasi kategori "Fashion"
-- =====================================================
-- Kondisi di production:
--   88b24214 → slug='fashion',   slug_id=NULL, slug_en=NULL  (row lama, 0 produk)
--   9cc2dd4f → slug='fashion-1', slug_id='fashion-1'         (row seed 000144)
-- =====================================================
-- Catatan: kolom `slug` punya UNIQUE constraint penuh (bukan partial).
-- Trigger rewrite_slug_on_delete sudah di-drop di migration 000088.
-- Oleh karena itu slug row lama harus di-rename saat soft-delete
-- agar unique constraint bebas untuk row seed.
-- =====================================================

-- Step 1: Soft-delete row lama + rename slug-nya
UPDATE kategori_produk
SET deleted_at = NOW(),
    is_active  = false,
    slug       = 'fashion-deleted-88b24214',
    updated_at = NOW()
WHERE id = '88b24214-20d2-4d64-ab32-4d703493ae8a'
  AND deleted_at IS NULL;

-- Step 2: Update row seed ke slug yang benar
UPDATE kategori_produk
SET slug       = 'fashion',
    slug_id    = 'fashion',
    slug_en    = 'fashion',
    updated_at = NOW()
WHERE id = '9cc2dd4f-cb59-47f8-8455-b4b229ab4e2c'
  AND deleted_at IS NULL;
