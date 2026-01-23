-- migrations/000090_seed_merek_produk.down.sql
-- Rollback: Hapus semua seed data merek_produk

DELETE FROM merek_produk 
WHERE slug IN (
    'samsung', 'lg', 'sony', 'philips', 'panasonic',
    'apple', 'xiaomi', 'lenovo', 'hp', 'dell',
    'asus', 'acer', 'canon', 'epson', 'brother', 'lainnya'
);
