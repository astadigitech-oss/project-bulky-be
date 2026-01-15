-- Rollback seed data untuk tabel sumber_produk
DELETE FROM sumber_produk 
WHERE slug IN ('retur', 'reject', 'overstock', 'closeout', 'excess', 'liquidasi');
