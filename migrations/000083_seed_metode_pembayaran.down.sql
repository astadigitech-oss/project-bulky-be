-- =====================================================
-- ROLLBACK: Metode Pembayaran Seed
-- =====================================================

DELETE FROM metode_pembayaran 
WHERE kode IN (
    'BCA', 'MANDIRI', 'BNI', 'BRI', 'PERMATA', 'CIMB', 'BSI', 'BJB',
    'LINKAJA', 'DANA', 'ID_OVO', 'ID_SHOPEEPAY',
    'CREDIT_CARD',
    'ID_QRIS'
);
