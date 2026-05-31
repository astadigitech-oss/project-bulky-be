-- migrations/000101_seed_faq_data.up.sql
-- FAQ data diambil dari dokumen_kebijakan section 5 (migration 084)
-- 17 Q&A, bilingual (ID + EN), grouped by category

-- Clear old FAQ data before re-seeding
DELETE FROM faq;

INSERT INTO faq (id, question, question_en, answer, answer_en, urutan, is_active, created_at, updated_at)
VALUES

-- Tentang Bulky.id / About Bulky.id
(
    'b0fa0001-0000-4000-8000-000000000001',
    'Apa itu Bulky.id?',
    'What is Bulky.id?',
    'Bulky.id adalah platform re-commerce grosir B2B Indonesia. Kami mengambil persediaan surplus, kelebihan stok, dan barang retur dari platform e-commerce, penyedia logistik, dan perusahaan, kemudian menyediakannya kepada reseller, pedagang, dan UKM dengan harga grosir. Kami berbasis di Jakarta.',
    'Bulky.id is Indonesia''s B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory from e-commerce platforms, logistics providers, and corporations, then make it available to resellers, traders, and SMEs at wholesale prices. We are based in Jakarta.',
    1, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000002',
    'Untuk siapa Bulky.id?',
    'Who is Bulky.id for?',
    'Bulky.id dirancang untuk pembeli bisnis — reseller, pedagang, UKM, dan distributor. Ini adalah platform B2B, bukan marketplace konsumen. Jika Anda ingin memulai atau mengembangkan bisnis resale, Bulky.id untuk Anda.',
    'Bulky.id is designed for business buyers — resellers, traders, SMEs, and distributors. It is a B2B platform, not a consumer marketplace. If you want to start or grow a resale business, Bulky.id is for you.',
    2, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000003',
    'Apakah ada minimum pembelian?',
    'Is there a minimum order?',
    'Tidak ada nilai minimum pesanan yang tetap. Produk di Bulky.id dijual dalam unit grosir yang sudah dikemas — pembelian minimum adalah unit terkecil yang tersedia untuk setiap listing, yang bisa berupa palet, satu truk, satu kontainer, atau jumlah grosir lain yang telah ditentukan. Ukuran dan jumlah unit yang tersedia tercantum jelas di setiap listing produk.',
    'There is no fixed minimum order value. Products on Bulky.id are sold in pre-packed bulk units — the minimum purchase is simply the smallest available unit for each listing, which may be a pallet, a truckload, a container load, or another pre-determined bulk quantity. The available unit size and quantity are clearly stated on each product listing.',
    3, true, NOW(), NOW()
),

-- Tentang Produk / About Products
(
    'b0fa0001-0000-4000-8000-000000000004',
    'Jenis produk apa yang dijual Bulky.id?',
    'What types of products does Bulky.id sell?',
    'Lot grosir dalam lima kategori: Fashion dan Pakaian, Elektronik, FMCG dan Perawatan Pribadi, Rumah Tangga dan Barang Umum, serta Aksesori Otomotif. Bersumber sebagai surplus, kelebihan stok, dan barang retur dari mitra e-commerce dan logistik yang sah.',
    'Wholesale lots across five categories: Fashion and Apparel, Electronics, FMCG and Personal Care, Household and General Merchandise, and Automotive Accessories. Sourced as surplus, overstock, and returned inventory from legitimate e-commerce and logistics partners.',
    4, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000005',
    'Apa saja grade kondisi?',
    'What are the condition grades?',
    'Baru (tersegel pabrik), Seperti Baru (dibuka tapi belum digunakan atau hampir baru), Bekas (pernah digunakan, berfungsi), dan Salvage (rusak atau tidak berfungsi, untuk suku cadang atau daur ulang). Grade ditampilkan di setiap listing produk.',
    'New (factory-sealed), Like New (opened but unused or near-new), Used (previously used, functional), and Salvage (damaged or non-functional, for parts or recycling). Grades are shown on each product listing.',
    5, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000006',
    'Apakah produk asli?',
    'Are the products authentic?',
    'Ya. Bulky.id mengambil sumber langsung dari platform e-commerce berlisensi, operator logistik, dan perusahaan. Semua produk melalui proses skrining penerimaan kami. Kami tidak dengan sengaja mendaftarkan produk palsu. Sebagai barang re-commerce, produk dijual apa adanya tanpa garansi merek.',
    'Yes. Bulky.id sources directly from licensed e-commerce platforms, logistics operators, and corporations. All products go through our intake screening process. We do not knowingly list counterfeit products. As re-commerce goods, products are sold as-is without brand warranties.',
    6, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000007',
    'Apa maksudnya dijual apa adanya (as-is)?',
    'What does as-is mean?',
    'Produk dijual dalam kondisi saat ini tanpa garansi tambahan dari Bulky.id. Unit individual dalam satu lot dapat bervariasi dalam kondisi, kelengkapan, dan spesifikasi — ini adalah sifat dari inventaris surplus dan barang retur. Harga mencerminkan hal ini.',
    'Products are sold in their current condition with no additional warranty from Bulky.id. Individual units within a lot may vary in condition, completeness, and specification — this is the nature of surplus and returned inventory. Pricing reflects this.',
    7, true, NOW(), NOW()
),

-- Tentang Pemesanan dan Pembayaran / About Ordering and Payment
(
    'b0fa0001-0000-4000-8000-000000000008',
    'Metode pembayaran apa yang diterima?',
    'What payment methods are accepted?',
    'Virtual Account (BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB), Kartu Kredit dan Debit (Visa, Mastercard, JCB), E-Wallet (OVO, DANA, GoPay, ShopeePay, LinkAja), dan QRIS. Semua pembayaran melalui Xendit. Tidak ada tunai atau COD.',
    'Virtual Account (BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB), Credit and Debit Cards (Visa, Mastercard, JCB), E-Wallets (OVO, DANA, GoPay, ShopeePay, LinkAja), and QRIS. All payments via Xendit. No cash or COD.',
    8, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000009',
    'Apa itu Virtual Account?',
    'What is a Virtual Account?',
    'Nomor rekening bank unik yang dibuat untuk transaksi Anda. Transfer jumlah yang tepat ke nomor VA yang diberikan dan pembayaran direkonsiliasi secara otomatis. Aman, cepat diatur, dan banyak digunakan untuk transaksi B2B di Indonesia.',
    'A unique bank account number generated for your transaction. Transfer the exact amount to the VA number provided and payment is automatically reconciled. Safe, instant to set up, and widely used for B2B transactions in Indonesia.',
    9, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000010',
    'Berapa lama saya harus membayar?',
    'How long do I have to pay?',
    '15 menit untuk pembelian reguler, atau 30 menit untuk pesanan Patungan. Pesanan yang tidak dibayar dalam batas waktu ini akan dibatalkan secara otomatis.',
    '15 minutes for regular purchases, or 30 minutes for Patungan split-payment orders. Orders not paid within this window are automatically cancelled.',
    10, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000011',
    'Apa itu Patungan?',
    'What is Patungan?',
    'Fitur split-payment pembelian bersama dari Bulky.id. Dua atau lebih pembeli terdaftar membagi biaya satu pesanan — masing-masing membayar bagian yang disepakati melalui saluran pembayaran mereka sendiri. Dirancang untuk reseller yang menggabungkan sumber daya untuk membeli lot lebih besar.',
    'Bulky.id''s co-purchase split-payment feature. Two or more registered buyers split the cost of one order — each paying their agreed share via their own payment channel. Designed for resellers pooling resources to buy larger lots.',
    11, true, NOW(), NOW()
),

-- Tentang Pengiriman dan Pengambilan / About Delivery and Pickup
(
    'b0fa0001-0000-4000-8000-000000000012',
    'Bisakah saya mengambil pesanan sendiri?',
    'Can I pick up my order?',
    'Ya. Pengambilan mandiri di gudang Cibinong kami (Jl. Raya Mayor Oking No. 62a, Cibinong, Jawa Barat), Senin hingga Sabtu 09:00 hingga 18:00 WIB. Bawa konfirmasi pesanan dan ID yang valid. Jadwalkan terlebih dahulu via WhatsApp.',
    'Yes. Self-pickup at our Cibinong warehouse (Jl. Raya Mayor Oking No. 62a, Cibinong, West Java), Monday to Saturday 09:00 to 18:00 WIB. Bring order confirmation and valid ID. Schedule in advance via WhatsApp.',
    12, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000013',
    'Apakah ada layanan pengiriman?',
    'Do you offer delivery?',
    'Ya. Bulky.id menawarkan pengiriman domestik ke seluruh Indonesia melalui Deliveree dan Forwarder.ai. Biaya pengiriman dihitung saat checkout berdasarkan jarak, berat, dan dimensi pesanan. Pengiriman mencakup seluruh kepulauan Indonesia.',
    'Yes. Bulky.id offers nationwide domestic delivery through Deliveree and Forwarder.ai. Delivery costs are calculated at checkout based on the distance, weight, and dimensions of your order. Delivery covers the entire Indonesian archipelago.',
    13, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000014',
    'Bagaimana jika pesanan saya memiliki item lebih sedikit dari yang saya bayar?',
    'What if my order has fewer items than I paid for?',
    'Periksa pesanan Anda saat diterima dan laporkan kekurangan jumlah dalam 24 jam dengan dokumentasi foto. Laporan setelah 24 jam atau tanpa dokumentasi tidak akan dipertimbangkan. Lihat Syarat dan Ketentuan lengkap kami untuk detailnya.',
    'Inspect your order upon receipt and report any quantity shortfall within 24 hours with photographic documentation. Reports after 24 hours or without documentation will not be considered. See our full Terms and Conditions for details.',
    14, true, NOW(), NOW()
),

-- Tentang Pengembalian dan Refund / About Returns and Refunds
(
    'b0fa0001-0000-4000-8000-000000000015',
    'Bisakah saya mengembalikan produk?',
    'Can I return products?',
    'Semua penjualan bersifat final. Pengembalian tidak diizinkan kecuali dalam keadaan yang sangat terbatas — kekurangan jumlah material atau terbuktinya kesalahan sortir dari Bulky.id — dilaporkan dalam 24 jam dengan dokumentasi. Harap tinjau Syarat dan Ketentuan lengkap kami sebelum membeli.',
    'All sales are final. Returns are not permitted except in narrowly defined circumstances — material quantity shortfall or confirmed mis-sort caused by Bulky.id — reported within 24 hours with documentation. Please review our full Terms and Conditions before purchasing.',
    15, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000016',
    'Bagaimana jika saya tidak puas dengan kondisi produk?',
    'What if I am unsatisfied with product condition?',
    'Variasi kondisi adalah hal yang melekat dalam inventaris surplus dan barang retur — inilah yang membuat penetapan harga menjadi mungkin. Ketidakpuasan dengan kondisi saja tidak memenuhi syarat untuk refund. Kami mendorong pembeli untuk meninjau grade kondisi dengan seksama dan membeli dengan ekspektasi yang realistis.',
    'Condition variance is inherent in surplus and returned inventory — this is what makes the pricing possible. Dissatisfaction with condition alone does not qualify for a refund. We encourage buyers to review condition grades carefully and purchase with realistic expectations.',
    16, true, NOW(), NOW()
),
(
    'b0fa0001-0000-4000-8000-000000000017',
    'Bagaimana refund diproses?',
    'How are refunds processed?',
    'Refund yang disetujui diproses ke metode pembayaran asli dalam 7 hingga 14 hari kerja.',
    'Approved refunds are processed to the original payment method within 7 to 14 business days.',
    17, true, NOW(), NOW()
)

ON CONFLICT (id) DO UPDATE SET
    question    = EXCLUDED.question,
    question_en = EXCLUDED.question_en,
    answer      = EXCLUDED.answer,
    answer_en   = EXCLUDED.answer_en,
    urutan      = EXCLUDED.urutan,
    is_active   = EXCLUDED.is_active,
    updated_at  = NOW();
