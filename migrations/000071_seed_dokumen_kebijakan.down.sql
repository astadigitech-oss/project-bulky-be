-- Remove seeded dokumen kebijakan data
DELETE FROM dokumen_kebijakan WHERE id IN (
    '9ceb8800-b705-41a3-b123-48c7108ddd18',
    '9ceb8800-b9ae-423c-9d55-dbf86c8c62d7',
    '9ceb8800-ba65-4467-b80b-1f6bed5bf9ba',
    '9ceb8800-bb37-40ea-a625-b95ef4835902',
    '9d3fa30e-e585-4c7a-b788-05ed97a360ec'
);
