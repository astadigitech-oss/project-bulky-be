-- =====================================================
-- REVISI: Kategori Produk - Kondisi Tambahan
-- =====================================================
-- Hapus field memiliki_kondisi_tambahan (redundant)
-- Ubah tipe_kondisi_tambahan jadi nullable enum uppercase
-- =====================================================

-- 1. Drop field memiliki_kondisi_tambahan
ALTER TABLE kategori_produk 
DROP COLUMN IF EXISTS memiliki_kondisi_tambahan;

-- 2. Update existing data: lowercase → uppercase
UPDATE kategori_produk 
SET tipe_kondisi_tambahan = UPPER(tipe_kondisi_tambahan)
WHERE tipe_kondisi_tambahan IS NOT NULL;

-- 3. Drop old constraint
ALTER TABLE kategori_produk 
DROP CONSTRAINT IF EXISTS kategori_produk_tipe_kondisi_tambahan_check;

-- 4. Add new constraint with uppercase enum
ALTER TABLE kategori_produk 
ADD CONSTRAINT kategori_produk_tipe_kondisi_tambahan_check 
CHECK (tipe_kondisi_tambahan IN ('TEKS', 'GAMBAR'));

-- 5. Update column type (optional - untuk clarity)
ALTER TABLE kategori_produk 
ALTER COLUMN tipe_kondisi_tambahan TYPE VARCHAR(10);

COMMENT ON COLUMN kategori_produk.tipe_kondisi_tambahan IS 
'Tipe kondisi tambahan: TEKS, GAMBAR, atau NULL jika tidak ada';
