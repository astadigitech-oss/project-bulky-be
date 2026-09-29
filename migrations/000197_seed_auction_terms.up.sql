INSERT INTO dokumen_kebijakan (
    judul, judul_en, slug, slug_id, slug_en, konten, konten_en,
    is_active, created_at, updated_at
)
VALUES (
    'Syarat dan Ketentuan Lelang',
    'Auction Terms and Conditions',
    'syarat-ketentuan-lelang',
    'syarat-ketentuan-lelang',
    'auction-terms-and-conditions',
    $$<h2>Syarat dan Ketentuan Lelang</h2>
<ol>
<li>Setiap bid yang dikirim bersifat final dan tidak dapat diubah atau dibatalkan.</li>
<li>Buyer wajib memastikan nominal bid dan estimasi biaya yang ditampilkan sudah dipahami sebelum mengirim bid.</li>
<li>Penetapan pemenang dilakukan oleh Bulky.id setelah periode lelang berakhir.</li>
<li>Buyer yang terpilih wajib mengikuti instruksi pembayaran dan pengambilan atau pengiriman yang diberikan oleh Bulky.id.</li>
<li>Barang lelang dijual sesuai kondisi dan informasi batch yang ditampilkan pada halaman lelang.</li>
</ol>$$,
    $$<h2>Auction Terms and Conditions</h2>
<ol>
<li>Every submitted bid is final and cannot be changed or cancelled.</li>
<li>The buyer must review and understand the bid amount and estimated costs before submitting.</li>
<li>Bulky.id selects the winner after the auction period ends.</li>
<li>The selected buyer must follow Bulky.id payment and fulfilment instructions.</li>
<li>Auction items are sold according to the condition and batch information shown on the auction page.</li>
</ol>$$,
    true, NOW(), NOW()
)
ON CONFLICT (slug) WHERE deleted_at IS NULL DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    is_active = true,
    updated_at = NOW();
