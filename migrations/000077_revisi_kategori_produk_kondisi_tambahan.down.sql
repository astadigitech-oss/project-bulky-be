-- =====================================================
-- ROLLBACK: Kembalikan ke struktur lama
-- =====================================================

-- 1. Kembalikan ke lowercase
UPDATE kategori_produk 
SET tipe_kondisi_tambahan = LOWER(tipe_kondisi_tambahan)
WHERE tipe_kondisi_tambahan IS NOT NULL;

-- 2. Drop new constraint
ALTER TABLE kategori_produk 
DROP CONSTRAINT IF EXISTS kategori_produk_tipe_kondisi_tambahan_check;

-- 3. Add old constraint
ALTER TABLE kategori_produk 
ADD CONSTRAINT kategori_produk_tipe_kondisi_tambahan_check 
CHECK (tipe_kondisi_tambahan IN ('gambar', 'teks'));

-- 4. Add back memiliki_kondisi_tambahan
ALTER TABLE kategori_produk 
ADD COLUMN memiliki_kondisi_tambahan BOOLEAN DEFAULT false;

-- 5. Populate memiliki_kondisi_tambahan based on tipe
UPDATE kategori_produk 
SET memiliki_kondisi_tambahan = (tipe_kondisi_tambahan IS NOT NULL);

-- 6. Revert column type
ALTER TABLE kategori_produk 
ALTER COLUMN tipe_kondisi_tambahan TYPE VARCHAR(20);
