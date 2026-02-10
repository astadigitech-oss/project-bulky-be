-- Rollback
DELETE FROM kategori_video WHERE slug IN (
    'tutorial', 'review-produk', 'tips-usaha', 'promo-event', 
    'behind-the-scene', 'testimoni', 'tips-belanja'
);
