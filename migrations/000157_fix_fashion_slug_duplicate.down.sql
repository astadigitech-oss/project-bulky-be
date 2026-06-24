-- migrations/000157_fix_fashion_slug_duplicate.down.sql

-- Revert row seed ke slug lama
UPDATE kategori_produk
SET slug       = 'fashion-1',
    slug_id    = 'fashion-1',
    slug_en    = 'fashion-1',
    updated_at = NOW()
WHERE id = '9cc2dd4f-cb59-47f8-8455-b4b229ab4e2c'
  AND deleted_at IS NULL;

-- Restore row lama
UPDATE kategori_produk
SET deleted_at = NULL,
    is_active  = true,
    slug       = 'fashion',
    updated_at = NOW()
WHERE id = '88b24214-20d2-4d64-ab32-4d703493ae8a';
