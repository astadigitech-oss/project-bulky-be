-- =====================================================
-- ROLLBACK: Metode Pembayaran Group Seed
-- =====================================================

DELETE FROM metode_pembayaran_group 
WHERE nama IN ('Bank Transfer / VA', 'E-Wallet', 'Kartu Kredit', 'QRIS');
