-- =====================================================
-- REVISI: Tambah nama_en ke tabel produk
-- =====================================================
-- Dual language support untuk nama produk
-- =====================================================

ALTER TABLE produk 
ADD COLUMN nama_en VARCHAR(255) NOT NULL DEFAULT '';

-- Rename "nama" to "nama_id" for clarity
ALTER TABLE produk 
RENAME COLUMN nama TO nama_id;

-- Index untuk search
CREATE INDEX idx_produk_nama_en ON produk(nama_en) 
    WHERE deleted_at IS NULL;

-- Comment
COMMENT ON COLUMN produk.nama_en IS 'Nama produk dalam bahasa Inggris (opsional)';
COMMENT ON COLUMN produk.nama_id IS 'Nama produk dalam bahasa Indonesia';