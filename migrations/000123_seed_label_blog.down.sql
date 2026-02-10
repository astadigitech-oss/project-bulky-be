-- Rollback
DELETE FROM label_blog WHERE slug IN (
    'trending', 'popular', 'rekomendasi', 'terbaru', 
    'pilihan-editor', 'grosir', 'hemat', 'peluang-usaha'
);
