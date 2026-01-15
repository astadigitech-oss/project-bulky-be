-- Seed data untuk tabel sumber_produk
INSERT INTO sumber_produk (id, nama_id, nama_en, slug, is_active, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Retur', 'Return', 'retur', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Reject', 'Reject', 'reject', true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'),
    (uuid_generate_v4(), 'Overstock', 'Overstock', 'overstock', true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'),
    (uuid_generate_v4(), 'Closeout', 'Closeout', 'closeout', true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second'),
    (uuid_generate_v4(), 'Excess', 'Excess', 'excess', true, NOW() + INTERVAL '4 second', NOW() + INTERVAL '4 second'),
    (uuid_generate_v4(), 'Liquidasi', 'Liquidation', 'liquidasi', true, NOW() + INTERVAL '5 second', NOW() + INTERVAL '5 second')
ON CONFLICT (slug) DO NOTHING;
