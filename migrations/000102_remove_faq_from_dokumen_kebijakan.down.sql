-- Re-insert FAQ record if rollback needed
INSERT INTO dokumen_kebijakan (judul, judul_en, slug, konten, konten_en, is_active)
VALUES (
    'Sering Ditanyakan',
    'Frequently Asked Questions',
    'faq',
    '[]',
    '[]',
    true
);
