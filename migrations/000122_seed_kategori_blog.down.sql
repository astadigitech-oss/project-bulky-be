-- Rollback
DELETE FROM kategori_blog WHERE slug IN (
    'tips-trick', 'berita', 'produk', 'tutorial', 'promo', 'bisnis', 'lifestyle'
);
