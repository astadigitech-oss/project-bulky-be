-- migrations/000090_seed_merek_produk.up.sql
-- Seed data untuk tabel merek_produk

INSERT INTO merek_produk (nama_id, nama_en, slug, logo_url, is_active, created_at, updated_at)
VALUES
    ('Samsung', 'Samsung', 'samsung', NULL, true, NOW(), NOW()),
    ('LG', 'LG', 'lg', NULL, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    ('Sony', 'Sony', 'sony', NULL, true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'),
    ('Philips', 'Philips', 'philips', NULL, true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second'),
    ('Panasonic', 'Panasonic', 'panasonic', NULL, true, NOW() + INTERVAL '4 second', NOW() + INTERVAL '4 second'),
    ('Apple', 'Apple', 'apple', NULL, true, NOW() + INTERVAL '5 second', NOW() + INTERVAL '5 second'),
    ('Xiaomi', 'Xiaomi', 'xiaomi', NULL, true, NOW() + INTERVAL '6 second', NOW() + INTERVAL '6 second'),
    ('Lenovo', 'Lenovo', 'lenovo', NULL, true, NOW() + INTERVAL '7 second', NOW() + INTERVAL '7 second'),
    ('HP', 'HP', 'hp', NULL, true, NOW() + INTERVAL '8 second', NOW() + INTERVAL '8 second'),
    ('Dell', 'Dell', 'dell', NULL, true, NOW() + INTERVAL '9 second', NOW() + INTERVAL '9 second'),
    ('Asus', 'Asus', 'asus', NULL, true, NOW() + INTERVAL '10 second', NOW() + INTERVAL '10 second'),
    ('Acer', 'Acer', 'acer', NULL, true, NOW() + INTERVAL '11 second', NOW() + INTERVAL '11 second'),
    ('Canon', 'Canon', 'canon', NULL, true, NOW() + INTERVAL '12 second', NOW() + INTERVAL '12 second'),
    ('Epson', 'Epson', 'epson', NULL, true, NOW() + INTERVAL '13 second', NOW() + INTERVAL '13 second'),
    ('Brother', 'Brother', 'brother', NULL, true, NOW() + INTERVAL '14 second', NOW() + INTERVAL '14 second'),
    ('Lainnya', 'Others', 'lainnya', NULL, true, NOW() + INTERVAL '15 second', NOW() + INTERVAL '15 second')
ON CONFLICT (slug) DO NOTHING;

COMMENT ON TABLE merek_produk IS 'Master data merek produk untuk B2B recommerce';
