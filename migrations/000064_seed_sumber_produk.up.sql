-- Seed data untuk tabel sumber_produk
INSERT INTO sumber_produk (id, nama_id, nama_en, slug, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Supplier Lokal', 'Local Supplier', 'supplier-lokal', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Impor', 'Import', 'impor', true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Lelang', 'Auction', 'lelang', true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'),
    (uuid_generate_v4(), 'Overstock', 'Overstock', 'overstock', true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second'),
    (uuid_generate_v4(), 'Retur', 'Return', 'retur', true, NOW() + INTERVAL '4 second', NOW() + INTERVAL '4 second'),
    (uuid_generate_v4(), 'Liquidasi', 'Liquidation', 'liquidasi', true, NOW() + INTERVAL '5 second', NOW() + INTERVAL '5 second'),
    (uuid_generate_v4(), 'Buyback', 'Buyback', 'buyback', true, NOW() + INTERVAL '6 second', NOW() + INTERVAL '6 second'),
    (uuid_generate_v4(), 'Trade-In', 'Trade-In', 'trade-in', true, NOW() + INTERVAL '7 second', NOW() + INTERVAL '7 second'),
    (uuid_generate_v4(), 'Ex-Display', 'Ex-Display', 'ex-display', true, NOW() + INTERVAL '8 second', NOW() + INTERVAL '8 second')
ON CONFLICT (slug) DO NOTHING;
