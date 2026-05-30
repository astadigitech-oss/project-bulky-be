-- Rollback seed data untuk tabel sumber_produk
DELETE FROM sumber_produk
WHERE slug IN ('supplier-lokal', 'impor', 'lelang', 'overstock', 'retur', 'liquidasi', 'buyback', 'trade-in', 'ex-display');
