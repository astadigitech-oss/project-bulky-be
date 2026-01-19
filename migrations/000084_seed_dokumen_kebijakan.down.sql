-- =====================================================
-- ROLLBACK: Dokumen Kebijakan Seed
-- =====================================================

DELETE FROM dokumen_kebijakan 
WHERE slug IN (
    'tentang-kami',
    'cara-membeli',
    'tentang-pembayaran',
    'hubungi-kami',
    'faq',
    'syarat-ketentuan',
    'kebijakan-privasi'
);
