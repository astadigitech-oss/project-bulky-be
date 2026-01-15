-- Seed data untuk tabel banner_event_promo
INSERT INTO banner_event_promo (id, nama, gambar_url_id, gambar_url_en, url_tujuan, urutan, is_active, tanggal_mulai, tanggal_selesai, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Promo Tahun Baru', 'banner-promo/promo-tahun-baru-id.jpg', 'banner-promo/promo-tahun-baru-en.jpg', '/promo/tahun-baru', 1, true, NOW(), NOW() + INTERVAL '30 days', NOW(), NOW()),
    (uuid_generate_v4(), 'Flash Sale Weekend', 'banner-promo/flash-sale-id.jpg', 'banner-promo/flash-sale-en.jpg', '/promo/flash-sale', 2, true, NOW(), NOW() + INTERVAL '7 days', NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Diskon Elektronik', 'banner-promo/diskon-elektronik-id.jpg', 'banner-promo/diskon-elektronik-en.jpg', '/kategori/elektronik', 3, true, NOW(), NOW() + INTERVAL '14 days', NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second')
ON CONFLICT DO NOTHING;
