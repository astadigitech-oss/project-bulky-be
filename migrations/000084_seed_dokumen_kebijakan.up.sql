-- =====================================================
-- SEED: Dokumen Kebijakan (7 Fixed Pages) - Bilingual
-- =====================================================
-- Content is in HTML format (sanitized by backend)
-- =====================================================

-- 1. Tentang Kami / About Us
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Tentang Kami',
    'About Us',
    'tentang-kami',
    'tentang-kami',
    'about-us',
    $$<h2>Siapa Kami</h2>
<p>Bulky.id adalah platform re-commerce grosir B2B terkemuka di Indonesia. Kami mengambil persediaan surplus, kelebihan stok, dan barang retur langsung dari platform e-commerce besar, penyedia logistik, dan perusahaan, memproses setiap barang melalui fasilitas gudang profesional kami, dan menyediakannya kepada reseller, pedagang, dan UKM dengan harga grosir yang transparan.</p>
<p>Didirikan pada tahun 2022 dan berkantor pusat di Jakarta, Bulky.id telah memproses lebih dari 700 metrik ton barang dan melayani ribuan reseller aktif di seluruh Indonesia. Pendekatan berbasis teknologi kami — menampilkan Sistem Manajemen Gudang (WMS) eksklusif, aplikasi mobile khusus, dan pelacakan inventaris real-time — memastikan transparansi, keandalan, dan efisiensi di setiap langkah.</p>
<h2>Visi</h2>
<p>Menjadi platform re-commerce grosir paling terpercaya dan transparan di Indonesia — memberdayakan generasi berikutnya reseller dan UKM Indonesia melalui akses andal ke inventaris surplus berkualitas dengan harga yang jujur dan kompetitif.</p>
<h2>Misi</h2>
<p>Bulky.id hadir untuk melayani reseller, pedagang, dan usaha kecil dengan kejujuran dan keandalan. Misi kami dijalankan melalui empat komitmen:</p>
<ol>
<li>Menyediakan reseller dan pedagang dengan akses yang andal dan konsisten ke inventaris surplus grosir berkualitas melalui platform berteknologi tinggi yang transparan.</li>
<li>Mendukung ekosistem UMKM Indonesia dengan menurunkan hambatan pengadaan grosir — memberikan kesempatan nyata bagi siapa saja untuk membangun bisnis resale yang layak.</li>
<li>Mendorong praktik komersial yang bertanggung jawab dengan memberikan barang surplus, kelebihan stok, dan barang retur kehidupan produktif kedua — mengurangi limbah komersial dan logistik.</li>
<li>Terus meningkatkan platform, operasi gudang, dan kualitas layanan kami untuk menghadirkan pengalaman pembelian grosir yang paling profesional, jujur, dan efisien di Indonesia.</li>
</ol>
<h2>Nilai-Nilai Kami</h2>
<ul>
<li><strong>Transparansi</strong> | Harga yang jelas, grade kondisi yang jujur, tanpa biaya tersembunyi. Apa yang Anda lihat adalah apa yang Anda dapatkan.</li>
<li><strong>Keandalan</strong> | Pasokan inventaris yang konsisten, operasi gudang profesional, dan jadwal yang berkomitmen — setiap pesanan, setiap saat.</li>
<li><strong>Aksesibilitas</strong> | Menyetarakan lapangan bermain bagi reseller dan UKM Indonesia. Inventaris grosir berkualitas seharusnya bukan hak istimewa pembeli besar saja.</li>
<li><strong>Keberlanjutan</strong> | Setiap pembelian memberi inventaris surplus kehidupan kedua. Bulky.id secara aktif mengurangi limbah komersial melalui re-commerce yang bertanggung jawab.</li>
<li><strong>Integritas</strong> | Praktik bisnis yang jujur, ketentuan yang jelas, dan komunikasi langsung — karena bisnis Anda bergantung pada mitra yang dapat Anda percaya.</li>
</ul>
<h2>Platform Kami</h2>
<h3>Yang Kami Jual</h3>
<ul>
<li><strong>Fashion &amp; Pakaian</strong> | 30% inventaris — kategori volume tertinggi</li>
<li><strong>Elektronik</strong> | 20% inventaris — nilai tertinggi; penilaian kondisi berbasis grade</li>
<li><strong>FMCG &amp; Perawatan Pribadi</strong> | 20% inventaris — item non-konsumtif diprioritaskan</li>
<li><strong>Rumah Tangga &amp; Umum</strong> | 20% inventaris — berkembang pesat melalui kemitraan merek</li>
<li><strong>Aksesori Otomotif</strong> | 10% inventaris</li>
</ul>
<h3>Gudang Kami</h3>
<p>Gudang utama kami terletak di Cibinong, Jawa Barat — memproses lebih dari 2.500 item per hari dengan kapasitas total melebihi 8.000 ton. Setiap item diterima, dinilai, diklasifikasikan, dan dikatalogkan melalui WMS eksklusif kami sebelum tersedia bagi pembeli di platform.</p>
<h3>Tingkatan Pembeli</h3>
<ul>
<li><strong>Bronze</strong> | Saat pendaftaran — akses katalog produk penuh</li>
<li><strong>Silver</strong> | Riwayat pembelian reguler — diskon tambahan, prioritas stok</li>
<li><strong>Gold</strong> | Pembeli volume tinggi — rabat, akses kontainer early-bird</li>
<li><strong>Platinum</strong> | Pembeli tier teratas — harga khusus, account manager dedicated, fulfillment prioritas</li>
</ul>
<h2>Teknologi Kami</h2>
<ul>
<li><strong>WMS Eksklusif</strong> — pelacakan inventaris real-time dan fulfillment pesanan dalam skala besar</li>
<li><strong>Aplikasi Mobile dan Platform Web</strong> — pemesanan yang mudah di semua perangkat</li>
<li><strong>Patungan Split-Payment</strong> — fitur split-payment grosir pertama di Indonesia untuk reseller yang menggabungkan sumber daya</li>
<li><strong>Live Inventory Feed</strong> — visibilitas stok real-time setiap saat</li>
</ul>
<p>Bulky.id dirancang khusus untuk industri re-commerce Indonesia — bukan marketplace generik.</p>$$,
    $$<h2>Who We Are</h2>
<p>Bulky.id is Indonesia's leading B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory directly from major e-commerce platforms, logistics providers, and corporations, process every item through our professional warehouse facility, and make it available to resellers, traders, and SMEs at transparent wholesale prices.</p>
<p>Founded in 2022 and headquartered in Jakarta, Bulky.id has processed over 700 metric tonnes of goods and serves thousands of active resellers across Indonesia. Our technology-first approach — featuring a proprietary Warehouse Management System (WMS), dedicated mobile app, and real-time inventory tracking — ensures transparency, reliability, and efficiency at every step.</p>
<h2>Vision</h2>
<p>To be Indonesia's most trusted and transparent wholesale re-commerce platform — empowering the next generation of Indonesian resellers and SMEs through reliable access to quality surplus inventory at honest, competitive prices.</p>
<h2>Mission</h2>
<p>Bulky.id exists to serve resellers, traders, and small businesses with honesty and reliability. Our mission is carried out through four commitments:</p>
<ol>
<li>To provide resellers and traders with reliable, consistent access to quality wholesale surplus inventory through a transparent, technology-enabled platform.</li>
<li>To support Indonesia's MSME ecosystem by lowering barriers to wholesale sourcing — giving anyone the real opportunity to build a viable resale business.</li>
<li>To promote responsible commercial practices by giving surplus, overstock, and returned goods a second productive life — reducing commercial and logistical waste.</li>
<li>To continuously improve our platform, warehouse operations, and service quality to deliver the most professional, honest, and efficient wholesale buying experience in Indonesia.</li>
</ol>
<h2>Our Values</h2>
<ul>
<li><strong>Transparency</strong> | Clear pricing, honest condition grades, no hidden fees. What you see is what you get.</li>
<li><strong>Reliability</strong> | Consistent inventory supply, professional warehouse operations, and committed timelines — every order, every time.</li>
<li><strong>Access</strong> | Leveling the playing field for Indonesian resellers and SMEs. Quality wholesale inventory should not be the privilege of large buyers only.</li>
<li><strong>Sustainability</strong> | Every purchase gives surplus inventory a second life. Bulky.id actively reduces commercial waste through responsible re-commerce.</li>
<li><strong>Integrity</strong> | Honest business practices, clear terms, and direct communication — because your business depends on a partner you can trust.</li>
</ul>
<h2>Our Platform</h2>
<h3>What We Sell</h3>
<ul>
<li><strong>Fashion &amp; Apparel</strong> | 30% of inventory — highest volume category</li>
<li><strong>Electronics</strong> | 20% of inventory — highest value; grade-based condition assessment</li>
<li><strong>FMCG &amp; Personal Care</strong> | 20% of inventory — non-consumptive items prioritised</li>
<li><strong>Household &amp; General</strong> | 20% of inventory — growing rapidly through brand partnerships</li>
<li><strong>Automotive Accessories</strong> | 10% of inventory</li>
</ul>
<h3>Our Warehouse</h3>
<p>Our main warehouse is located in Cibinong, West Java — processing over 2,500 items per day with total capacity exceeding 8,000 tonnes. Every item is received, assessed, graded, and catalogued through our proprietary WMS before being made available to buyers on the platform.</p>
<h3>Buyer Tiers</h3>
<ul>
<li><strong>Bronze</strong> | Upon registration — full product catalogue access</li>
<li><strong>Silver</strong> | Regular purchase history — additional discounts, stock priority</li>
<li><strong>Gold</strong> | High-volume buyer — rebates, early-bird container access</li>
<li><strong>Platinum</strong> | Top-tier buyer — custom pricing, dedicated account manager, priority fulfillment</li>
</ul>
<h2>Our Technology</h2>
<ul>
<li><strong>Proprietary WMS</strong> — real-time inventory tracking and order fulfillment at scale</li>
<li><strong>Mobile App and Web Platform</strong> — seamless ordering across all devices</li>
<li><strong>Patungan Split-Payment</strong> — Indonesia's first wholesale split-payment feature for resellers who pool resources</li>
<li><strong>Live Inventory Feed</strong> — real-time stock visibility at all times</li>
</ul>
<p>Bulky.id is purpose-built for the Indonesian re-commerce industry — not a generic marketplace.</p>$$,
    1, true, NOW(), NOW()
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 2. Cara Membeli / How to Buy
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Cara Membeli',
    'How to Buy',
    'cara-membeli',
    'cara-membeli',
    'how-to-buy',
    $$<p>Membeli di Bulky.id sangat mudah. Ikuti langkah-langkah di bawah ini untuk menempatkan pesanan pertama Anda. Untuk pertanyaan, hubungi tim kami via WhatsApp atau email.</p>
<h2>Alur Pembelian Utama</h2>
<h3>LANGKAH 1 — Buat Akun Anda</h3>
<p>Unduh aplikasi Bulky.id (App Store atau Google Play) atau kunjungi bulky.id. Klik Daftar dan lengkapi profil dengan nama, nomor telepon, email, dan informasi bisnis. Aktivasi akun biasanya dalam 1 x 24 jam.</p>
<h3>LANGKAH 2 — Telusuri Produk</h3>
<p>Setelah masuk, telusuri katalog berdasarkan kategori atau gunakan pencarian. Setiap listing menampilkan grade kondisi, jumlah, dan harga. Tinjau foto dan deskripsi kondisi dengan seksama.</p>
<h3>LANGKAH 3 — Tambahkan ke Keranjang</h3>
<p>Pilih produk dan jumlah lalu tambahkan ke keranjang. Anda dapat menambahkan beberapa item sebelum checkout.</p>
<h3>LANGKAH 4 — Tinjau Keranjang</h3>
<p>Buka keranjang, konfirmasi jumlah dan total, lalu lanjutkan ke checkout.</p>
<h3>LANGKAH 5 — Pilih Pengiriman atau Pengambilan</h3>
<p>Pilih: (a) Pengambilan Mandiri di Gudang Cibinong; atau (b) Pengiriman — masukkan alamat pengiriman dan pilih Deliveree atau Forwarder.ai untuk pengiriman domestik ke seluruh Indonesia. Biaya pengiriman dihitung saat checkout berdasarkan jarak, berat, dan dimensi.</p>
<h3>LANGKAH 6 — Pilih Metode Pembayaran</h3>
<p>Pilih <strong>Bayar Langsung</strong> untuk membayar penuh sendiri, atau <strong>Bayar Patungan dengan Teman</strong> untuk berbagi biaya dengan mitra. Lihat panduan Patungan di bawah.</p>
<h3>LANGKAH 7 — Selesaikan Pembayaran</h3>
<p>Pilih saluran pembayaran (Virtual Account, Kartu, E-Wallet, atau QRIS) dan selesaikan dalam batas waktu yang ditentukan. Pesanan yang tidak dibayar akan dibatalkan otomatis.</p>
<h3>LANGKAH 8 — Konfirmasi Pesanan</h3>
<p>Setelah pembayaran dikonfirmasi, Anda menerima konfirmasi melalui aplikasi dan WhatsApp. Pesanan kemudian diproses untuk pengambilan atau pengiriman.</p>
<h2>Patungan (Pembayaran Bersama) — Langkah demi Langkah</h2>
<p>Patungan memungkinkan dua atau lebih pembeli untuk membagi biaya satu pesanan. Dirancang untuk reseller yang menggabungkan sumber daya untuk membeli lot lebih besar bersama-sama.</p>
<h3>LANGKAH 1 — Pilih Item dan Lanjutkan ke Checkout</h3>
<p>Tambahkan item ke keranjang dan lanjutkan ke checkout seperti biasa.</p>
<h3>LANGKAH 2 — Pilih Metode Pengiriman</h3>
<p>Pilih opsi pengiriman atau pengambilan yang diinginkan.</p>
<h3>LANGKAH 3 — Pilih Bayar Patungan dengan Teman</h3>
<p>Di layar pembayaran, pilih <strong>Bayar Patungan dengan Teman</strong> sebagai pengganti <strong>Bayar Langsung</strong>.</p>
<h3>LANGKAH 4 — Undang Mitra Anda</h3>
<p>Klik <strong>Undang Temanmu</strong> dan masukkan alamat email Bulky.id terdaftar mitra Anda secara lengkap. Tap <strong>Buat Pesanan</strong> setelah semua mitra ditambahkan.</p>
<h3>LANGKAH 5 — Masukkan Bagian Anda</h3>
<p>Di layar Nominal, masukkan jumlah yang akan Anda bayar sesuai kesepakatan. Layar menampilkan bagian Anda, saldo yang tersisa, dan total pesanan termasuk pajak.</p>
<h3>LANGKAH 6 — Selesaikan Pembayaran Anda</h3>
<p>Pilih saluran pembayaran dan selesaikan bagian Anda. Pesan WhatsApp berisi link split-bill dikirim otomatis setelah pembayaran Anda diproses.</p>
<h3>LANGKAH 7 — Mitra Menyelesaikan Pembayarannya</h3>
<p>Mitra menerima link split-bill via WhatsApp, mengkliknya, memasukkan bagian mereka, memilih saluran pembayaran, dan menyelesaikan pembayaran.</p>
<h3>LANGKAH 8 — Konfirmasi Admin</h3>
<p>Setelah semua pihak membayar penuh, Bulky.id meninjau dan mengonfirmasi pesanan. Konfirmasi akhir dikirim via aplikasi dan WhatsApp.</p>
<p><strong>Penting:</strong> Semua pihak dalam transaksi Patungan bertanggung jawab bersama atas nilai pesanan penuh. Pastikan semua mitra berkomitmen sebelum memulai.</p>
<h2>Informasi Pengambilan Mandiri</h2>
<table>
<thead><tr><th>Informasi</th><th>Detail</th></tr></thead>
<tbody>
<tr><td><strong>Alamat Gudang</strong></td><td>Jl. Raya Mayor Oking No. 62a, Cibinong, Jawa Barat</td></tr>
<tr><td><strong>Jam Operasional</strong></td><td>Senin hingga Sabtu, 09:00 hingga 18:00 WIB</td></tr>
<tr><td><strong>Yang Dibawa</strong></td><td>Konfirmasi pesanan (dari aplikasi atau cetak) dan KTP/identitas resmi yang valid</td></tr>
<tr><td><strong>Penjadwalan</strong></td><td>Harap jadwalkan pengambilan terlebih dahulu via WhatsApp agar pesanan Anda sudah siap saat tiba.</td></tr>
</tbody>
</table>$$,
    $$<p>Buying at Bulky.id is easy. Follow the steps below to place your first order. For enquiries, contact our team via WhatsApp or email.</p>
<h2>Main Purchase Flow</h2>
<h3>STEP 1 — Create Your Account</h3>
<p>Download the Bulky.id app (App Store or Google Play) or visit bulky.id. Click Register and complete your profile with your name, phone number, email, and business information. Account activation is typically within 1 x 24 hours.</p>
<h3>STEP 2 — Browse Products</h3>
<p>Once logged in, browse the catalogue by category or use the search function. Each listing shows the condition grade, quantity, and price. Review photos and condition descriptions carefully.</p>
<h3>STEP 3 — Add to Cart</h3>
<p>Select the product and quantity, then add to cart. You can add multiple items before checking out.</p>
<h3>STEP 4 — Review Your Cart</h3>
<p>Open your cart, confirm quantities and total, then proceed to checkout.</p>
<h3>STEP 5 — Choose Delivery or Pickup</h3>
<p>Select: (a) Self-Pickup at the Cibinong Warehouse; or (b) Delivery — enter your delivery address and select Deliveree or Forwarder.ai for domestic delivery across Indonesia. Delivery costs are calculated at checkout based on distance, weight, and dimensions.</p>
<h3>STEP 6 — Choose Payment Method</h3>
<p>Select <strong>Pay Directly</strong> to pay in full yourself, or <strong>Pay Patungan with Friends</strong> to share the cost with partners. See the Patungan guide below.</p>
<h3>STEP 7 — Complete Payment</h3>
<p>Choose your payment channel (Virtual Account, Card, E-Wallet, or QRIS) and complete within the specified time limit. Unpaid orders will be automatically cancelled.</p>
<h3>STEP 8 — Order Confirmation</h3>
<p>Once payment is confirmed, you receive confirmation via the app and WhatsApp. The order is then processed for pickup or delivery.</p>
<h2>Patungan (Co-Payment) — Step by Step</h2>
<p>Patungan allows two or more buyers to split the cost of a single order. Designed for resellers pooling resources to purchase larger lots together.</p>
<h3>STEP 1 — Select Items and Proceed to Checkout</h3>
<p>Add items to cart and proceed to checkout as normal.</p>
<h3>STEP 2 — Choose Delivery Method</h3>
<p>Select your preferred delivery or pickup option.</p>
<h3>STEP 3 — Select Pay Patungan with Friends</h3>
<p>On the payment screen, select <strong>Pay Patungan with Friends</strong> instead of <strong>Pay Directly</strong>.</p>
<h3>STEP 4 — Invite Your Partners</h3>
<p>Click <strong>Invite Your Friends</strong> and enter each partner's registered Bulky.id email address in full. Tap <strong>Create Order</strong> once all partners have been added.</p>
<h3>STEP 5 — Enter Your Share</h3>
<p>On the Amount screen, enter the amount you will pay as agreed. The screen shows your share, the remaining balance, and the full order total including tax.</p>
<h3>STEP 6 — Complete Your Payment</h3>
<p>Choose your payment channel and complete your share. A WhatsApp message containing the split-bill link is sent automatically after your payment is processed.</p>
<h3>STEP 7 — Partners Complete Their Payment</h3>
<p>Each partner receives the split-bill link via WhatsApp, clicks it, enters their share, selects a payment channel, and completes payment.</p>
<h3>STEP 8 — Admin Confirmation</h3>
<p>Once all parties have paid in full, Bulky.id reviews and confirms the order. Final confirmation is sent via the app and WhatsApp.</p>
<p><strong>Important:</strong> All parties in a Patungan transaction are jointly and severally liable for the full order value. Ensure all partners are committed before starting.</p>
<h2>Self-Pickup Information</h2>
<table>
<thead><tr><th>Information</th><th>Detail</th></tr></thead>
<tbody>
<tr><td><strong>Warehouse Address</strong></td><td>Jl. Raya Mayor Oking No. 62a, Cibinong, West Java</td></tr>
<tr><td><strong>Operating Hours</strong></td><td>Monday to Saturday, 09:00 to 18:00 WIB</td></tr>
<tr><td><strong>What to Bring</strong></td><td>Order confirmation (from app or printed) and valid government-issued ID</td></tr>
<tr><td><strong>Scheduling</strong></td><td>Please schedule your pickup in advance via WhatsApp so your order is ready upon arrival.</td></tr>
</tbody>
</table>$$,
    2, true, NOW() + INTERVAL '1 second', NOW() + INTERVAL '1 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 3. Tentang Pembayaran / About Payment
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Tentang Pembayaran',
    'About Payment',
    'tentang-pembayaran',
    'tentang-pembayaran',
    'about-payment',
    $$<h2>Virtual Account (Transfer Bank)</h2>
<p>Pembayaran dari berbagai bank dapat dikenali dan diterima secara otomatis dan mudah hanya dengan satu Virtual Account, tanpa harus menggunakan rekening dari berbagai bank.</p>
<p><strong>Bank yang didukung:</strong> BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB</p>
<h3>Cara Bayar via Virtual Account BRI (BRIVA)</h3>
<ol>
<li>Masukkan kartu ATM kemudian masukkan PIN Anda (Pilih Bahasa)</li>
<li>Pilih Menu Lainnya</li>
<li>Pilih Pembayaran / Pembelian</li>
<li>Pilih BRIVA</li>
<li>Masukkan Nomor BRIVA (BRI Virtual Account) yang tertera di halaman pesanan</li>
</ol>
<h2>Kartu Kredit &amp; Debit</h2>
<p>Bulky.id menerima pembayaran kartu kredit dan debit dari jaringan berikut:</p>
<ul>
<li>Visa</li>
<li>Mastercard</li>
<li>JCB</li>
</ul>
<h2>E-Wallet</h2>
<p>Bayar dengan mudah menggunakan dompet digital Anda:</p>
<ul>
<li>OVO</li>
<li>DANA</li>
<li>GoPay</li>
<li>ShopeePay</li>
<li>LinkAja</li>
</ul>
<h2>QRIS</h2>
<p>Scan kode QRIS untuk pembayaran instan dari aplikasi mobile banking atau e-wallet manapun yang mendukung QRIS.</p>
<p><strong>Catatan:</strong> Tidak ada pembayaran tunai atau COD. Semua pembayaran diproses melalui Xendit.</p>$$,
    $$<h2>Virtual Accounts (Bank Transfer)</h2>
<p>Payments from various banks can be recognised and received automatically and easily through a single Virtual Account, without needing to use accounts from multiple banks.</p>
<p><strong>Supported banks:</strong> BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB</p>
<h3>How to Pay via BRI Virtual Account (BRIVA)</h3>
<ol>
<li>Insert your ATM card and enter your PIN (Select Language)</li>
<li>Select Other Menu</li>
<li>Select Payment / Purchase</li>
<li>Select BRIVA</li>
<li>Enter the BRIVA Number (BRI Virtual Account) shown on your order page</li>
</ol>
<h2>Credit &amp; Debit Cards</h2>
<p>Bulky.id accepts credit and debit card payments from the following networks:</p>
<ul>
<li>Visa</li>
<li>Mastercard</li>
<li>JCB</li>
</ul>
<h2>E-Wallets</h2>
<p>Pay conveniently using your digital wallet:</p>
<ul>
<li>OVO</li>
<li>DANA</li>
<li>GoPay</li>
<li>ShopeePay</li>
<li>LinkAja</li>
</ul>
<h2>QRIS</h2>
<p>Scan the QRIS code for instant payment from any mobile banking app or e-wallet that supports QRIS.</p>
<p><strong>Note:</strong> No cash or COD. All payments are processed via Xendit.</p>$$,
    3, true, NOW() + INTERVAL '2 second', NOW() + INTERVAL '2 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 4. Hubungi Kami / Contact Us
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Hubungi Kami',
    'Contact Us',
    'hubungi-kami',
    'hubungi-kami',
    'contact-us',
    $$<p>Kami siap membantu Anda. Hubungi kami melalui:</p>
<h2>Alamat Gudang</h2>
<p>Jl. Raya Mayor Oking Jaya Atmaja No.62a, Kel Cirimekar, Kec. Cibinong, Kabupaten Bogor, Jawa Barat 16918</p>
<h2>Telepon / WhatsApp</h2>
<p><a href="https://wa.me/62811833164">+62811 833 164</a></p>
<h2>Email</h2>
<p><a href="mailto:admin@bulky.id">admin@bulky.id</a></p>
<h2>Jam Operasional</h2>
<p>Senin – Sabtu: 09:00 – 18:00 WIB<br>Minggu: Tutup</p>
<h2>Media Sosial</h2>
<ul>
<li><strong>Instagram:</strong> <a href="https://www.instagram.com/bulky.id/" target="_blank">@bulky.id</a></li>
<li><strong>Facebook:</strong> <a href="https://www.facebook.com/liquid8wholesale" target="_blank">liquid8wholesale</a></li>
<li><strong>TikTok:</strong> <a href="https://www.tiktok.com/@bulky.id" target="_blank">@bulky.id</a></li>
</ul>$$,
    $$<p>We are here to help. Contact us through:</p>
<h2>Warehouse Address</h2>
<p>Jl. Raya Mayor Oking Jaya Atmaja No.62a, Kel Cirimekar, Kec. Cibinong, Kabupaten Bogor, West Java 16918</p>
<h2>Phone / WhatsApp</h2>
<p><a href="https://wa.me/62811833164">+62811 833 164</a></p>
<h2>Email</h2>
<p><a href="mailto:admin@bulky.id">admin@bulky.id</a></p>
<h2>Operating Hours</h2>
<p>Monday to Saturday: 09:00 to 18:00 WIB<br>Sunday: Closed</p>
<h2>Social Media</h2>
<ul>
<li><strong>Instagram:</strong> <a href="https://www.instagram.com/bulky.id/" target="_blank">@bulky.id</a></li>
<li><strong>Facebook:</strong> <a href="https://www.facebook.com/liquid8wholesale" target="_blank">liquid8wholesale</a></li>
<li><strong>TikTok:</strong> <a href="https://www.tiktok.com/@bulky.id" target="_blank">@bulky.id</a></li>
</ul>$$,
    4, true, NOW() + INTERVAL '3 second', NOW() + INTERVAL '3 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 5. FAQ / Sering Ditanyakan
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Sering Ditanyakan',
    'Frequently Asked Questions',
    'faq',
    'faq',
    'faq',
    $$<p>Tidak menemukan yang Anda cari? Hubungi kami via WhatsApp di <a href="https://wa.me/62811833164">0811 833 164</a> atau email <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
<h2>Tentang Bulky.id</h2>
<h3>Apa itu Bulky.id?</h3>
<p>Bulky.id adalah platform re-commerce grosir B2B Indonesia. Kami mengambil persediaan surplus, kelebihan stok, dan barang retur dari platform e-commerce, penyedia logistik, dan perusahaan, kemudian menyediakannya kepada reseller, pedagang, dan UKM dengan harga grosir. Kami berbasis di Jakarta.</p>
<h3>Untuk siapa Bulky.id?</h3>
<p>Bulky.id dirancang untuk pembeli bisnis — reseller, pedagang, UKM, dan distributor. Ini adalah platform B2B, bukan marketplace konsumen. Jika Anda ingin memulai atau mengembangkan bisnis resale, Bulky.id untuk Anda.</p>
<h3>Apakah ada minimum pembelian?</h3>
<p>Tidak ada nilai minimum pesanan yang tetap. Produk di Bulky.id dijual dalam unit grosir yang sudah dikemas — pembelian minimum adalah unit terkecil yang tersedia untuk setiap listing, yang bisa berupa palet, satu truk, satu kontainer, atau jumlah grosir lain yang telah ditentukan. Ukuran dan jumlah unit yang tersedia tercantum jelas di setiap listing produk.</p>
<h2>Tentang Produk</h2>
<h3>Jenis produk apa yang dijual Bulky.id?</h3>
<p>Lot grosir dalam lima kategori: Fashion dan Pakaian, Elektronik, FMCG dan Perawatan Pribadi, Rumah Tangga dan Barang Umum, serta Aksesori Otomotif. Bersumber sebagai surplus, kelebihan stok, dan barang retur dari mitra e-commerce dan logistik yang sah.</p>
<h3>Apa saja grade kondisi?</h3>
<p>Baru (tersegel pabrik), Seperti Baru (dibuka tapi belum digunakan atau hampir baru), Bekas (pernah digunakan, berfungsi), dan Salvage (rusak atau tidak berfungsi, untuk suku cadang atau daur ulang). Grade ditampilkan di setiap listing produk.</p>
<h3>Apakah produk asli?</h3>
<p>Ya. Bulky.id mengambil sumber langsung dari platform e-commerce berlisensi, operator logistik, dan perusahaan. Semua produk melalui proses skrining penerimaan kami. Kami tidak dengan sengaja mendaftarkan produk palsu. Sebagai barang re-commerce, produk dijual apa adanya tanpa garansi merek.</p>
<h3>Apa maksudnya dijual apa adanya (as-is)?</h3>
<p>Produk dijual dalam kondisi saat ini tanpa garansi tambahan dari Bulky.id. Unit individual dalam satu lot dapat bervariasi dalam kondisi, kelengkapan, dan spesifikasi — ini adalah sifat dari inventaris surplus dan barang retur. Harga mencerminkan hal ini.</p>
<h2>Tentang Pemesanan dan Pembayaran</h2>
<h3>Metode pembayaran apa yang diterima?</h3>
<p>Virtual Account (BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB), Kartu Kredit dan Debit (Visa, Mastercard, JCB), E-Wallet (OVO, DANA, GoPay, ShopeePay, LinkAja), dan QRIS. Semua pembayaran melalui Xendit. Tidak ada tunai atau COD.</p>
<h3>Apa itu Virtual Account?</h3>
<p>Nomor rekening bank unik yang dibuat untuk transaksi Anda. Transfer jumlah yang tepat ke nomor VA yang diberikan dan pembayaran direkonsiliasi secara otomatis. Aman, cepat diatur, dan banyak digunakan untuk transaksi B2B di Indonesia.</p>
<h3>Berapa lama saya harus membayar?</h3>
<p>15 menit untuk pembelian reguler, atau 30 menit untuk pesanan Patungan. Pesanan yang tidak dibayar dalam batas waktu ini akan dibatalkan secara otomatis.</p>
<h3>Apa itu Patungan?</h3>
<p>Fitur split-payment pembelian bersama dari Bulky.id. Dua atau lebih pembeli terdaftar membagi biaya satu pesanan — masing-masing membayar bagian yang disepakati melalui saluran pembayaran mereka sendiri. Dirancang untuk reseller yang menggabungkan sumber daya untuk membeli lot lebih besar.</p>
<h2>Tentang Pengiriman dan Pengambilan</h2>
<h3>Bisakah saya mengambil pesanan sendiri?</h3>
<p>Ya. Pengambilan mandiri di gudang Cibinong kami (Jl. Raya Mayor Oking No. 62a, Cibinong, Jawa Barat), Senin hingga Sabtu 09:00 hingga 18:00 WIB. Bawa konfirmasi pesanan dan ID yang valid. Jadwalkan terlebih dahulu via WhatsApp.</p>
<h3>Apakah ada layanan pengiriman?</h3>
<p>Ya. Bulky.id menawarkan pengiriman domestik ke seluruh Indonesia melalui Deliveree dan Forwarder.ai. Biaya pengiriman dihitung saat checkout berdasarkan jarak, berat, dan dimensi pesanan. Pengiriman mencakup seluruh kepulauan Indonesia.</p>
<h3>Bagaimana jika pesanan saya memiliki item lebih sedikit dari yang saya bayar?</h3>
<p>Periksa pesanan Anda saat diterima dan laporkan kekurangan jumlah dalam 24 jam dengan dokumentasi foto. Laporan setelah 24 jam atau tanpa dokumentasi tidak akan dipertimbangkan. Lihat Syarat dan Ketentuan lengkap kami untuk detailnya.</p>
<h2>Tentang Pengembalian dan Refund</h2>
<h3>Bisakah saya mengembalikan produk?</h3>
<p>Semua penjualan bersifat final. Pengembalian tidak diizinkan kecuali dalam keadaan yang sangat terbatas — kekurangan jumlah material atau terbuktinya kesalahan sortir dari Bulky.id — dilaporkan dalam 24 jam dengan dokumentasi. Harap tinjau Syarat dan Ketentuan lengkap kami sebelum membeli.</p>
<h3>Bagaimana jika saya tidak puas dengan kondisi produk?</h3>
<p>Variasi kondisi adalah hal yang melekat dalam inventaris surplus dan barang retur — inilah yang membuat penetapan harga menjadi mungkin. Ketidakpuasan dengan kondisi saja tidak memenuhi syarat untuk refund. Kami mendorong pembeli untuk meninjau grade kondisi dengan seksama dan membeli dengan ekspektasi yang realistis.</p>
<h3>Bagaimana refund diproses?</h3>
<p>Refund yang disetujui diproses ke metode pembayaran asli dalam 7 hingga 14 hari kerja.</p>$$,
    $$<p>Can't find what you're looking for? Contact us via WhatsApp at <a href="https://wa.me/62811833164">0811 833 164</a> or email <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
<h2>About Bulky.id</h2>
<h3>What is Bulky.id?</h3>
<p>Bulky.id is Indonesia's B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory from e-commerce platforms, logistics providers, and corporations, then make it available to resellers, traders, and SMEs at wholesale prices. We are based in Jakarta.</p>
<h3>Who is Bulky.id for?</h3>
<p>Bulky.id is designed for business buyers — resellers, traders, SMEs, and distributors. It is a B2B platform, not a consumer marketplace. If you want to start or grow a resale business, Bulky.id is for you.</p>
<h3>Is there a minimum order?</h3>
<p>There is no fixed minimum order value. Products on Bulky.id are sold in pre-packed bulk units — the minimum purchase is simply the smallest available unit for each listing, which may be a pallet, a truckload, a container load, or another pre-determined bulk quantity. The available unit size and quantity are clearly stated on each product listing.</p>
<h2>About Products</h2>
<h3>What types of products does Bulky.id sell?</h3>
<p>Wholesale lots across five categories: Fashion and Apparel, Electronics, FMCG and Personal Care, Household and General Merchandise, and Automotive Accessories. Sourced as surplus, overstock, and returned inventory from legitimate e-commerce and logistics partners.</p>
<h3>What are the condition grades?</h3>
<p>New (factory-sealed), Like New (opened but unused or near-new), Used (previously used, functional), and Salvage (damaged or non-functional, for parts or recycling). Grades are shown on each product listing.</p>
<h3>Are the products authentic?</h3>
<p>Yes. Bulky.id sources directly from licensed e-commerce platforms, logistics operators, and corporations. All products go through our intake screening process. We do not knowingly list counterfeit products. As re-commerce goods, products are sold as-is without brand warranties.</p>
<h3>What does as-is mean?</h3>
<p>Products are sold in their current condition with no additional warranty from Bulky.id. Individual units within a lot may vary in condition, completeness, and specification — this is the nature of surplus and returned inventory. Pricing reflects this.</p>
<h2>About Ordering and Payment</h2>
<h3>What payment methods are accepted?</h3>
<p>Virtual Account (BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB), Credit and Debit Cards (Visa, Mastercard, JCB), E-Wallets (OVO, DANA, GoPay, ShopeePay, LinkAja), and QRIS. All payments via Xendit. No cash or COD.</p>
<h3>What is a Virtual Account?</h3>
<p>A unique bank account number generated for your transaction. Transfer the exact amount to the VA number provided and payment is automatically reconciled. Safe, instant to set up, and widely used for B2B transactions in Indonesia.</p>
<h3>How long do I have to pay?</h3>
<p>15 minutes for regular purchases, or 30 minutes for Patungan split-payment orders. Orders not paid within this window are automatically cancelled.</p>
<h3>What is Patungan?</h3>
<p>Bulky.id's co-purchase split-payment feature. Two or more registered buyers split the cost of one order — each paying their agreed share via their own payment channel. Designed for resellers pooling resources to buy larger lots.</p>
<h2>About Delivery and Pickup</h2>
<h3>Can I pick up my order?</h3>
<p>Yes. Self-pickup at our Cibinong warehouse (Jl. Raya Mayor Oking No. 62a, Cibinong, West Java), Monday to Saturday 09:00 to 18:00 WIB. Bring order confirmation and valid ID. Schedule in advance via WhatsApp.</p>
<h3>Do you offer delivery?</h3>
<p>Yes. Bulky.id offers nationwide domestic delivery through Deliveree and Forwarder.ai. Delivery costs are calculated at checkout based on the distance, weight, and dimensions of your order. Delivery covers the entire Indonesian archipelago.</p>
<h3>What if my order has fewer items than I paid for?</h3>
<p>Inspect your order upon receipt and report any quantity shortfall within 24 hours with photographic documentation. Reports after 24 hours or without documentation will not be considered. See our full Terms and Conditions for details.</p>
<h2>About Returns and Refunds</h2>
<h3>Can I return products?</h3>
<p>All sales are final. Returns are not permitted except in narrowly defined circumstances — material quantity shortfall or confirmed mis-sort caused by Bulky.id — reported within 24 hours with documentation. Please review our full Terms and Conditions before purchasing.</p>
<h3>What if I am unsatisfied with product condition?</h3>
<p>Condition variance is inherent in surplus and returned inventory — this is what makes the pricing possible. Dissatisfaction with condition alone does not qualify for a refund. We encourage buyers to review condition grades carefully and purchase with realistic expectations.</p>
<h3>How are refunds processed?</h3>
<p>Approved refunds are processed to the original payment method within 7 to 14 business days.</p>$$,
    5, true, NOW() + INTERVAL '4 second', NOW() + INTERVAL '4 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 6. Syarat dan Ketentuan / Terms and Conditions
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Syarat dan Ketentuan',
    'Terms and Conditions',
    'syarat-ketentuan',
    'syarat-ketentuan',
    'terms-and-conditions',
    $$<p>Syarat dan Ketentuan Pembelian lengkap mengatur semua transaksi di Bulky.id. Ringkasan di bawah ini menyoroti poin-poin terpenting. Dengan membeli di Bulky.id, Anda setuju untuk terikat oleh Syarat dan Ketentuan lengkap — harap baca seluruhnya sebelum pembelian pertama Anda.</p>
<p>Syarat dan Ketentuan Pembeli lengkap (Versi 3.0, berlaku 1 Juni 2026) harus ditinjau dan diterima saat checkout. Di bawah ini hanya ringkasan dan tidak menggantikan dokumen lengkap.</p>
<h2>Prinsip Utama</h2>
<table>
<thead><tr><th>Prinsip</th><th>Penjelasan</th></tr></thead>
<tbody>
<tr><td><strong>Tujuan Platform</strong></td><td>Platform re-commerce grosir B2B untuk reseller, pedagang, dan UKM — bukan marketplace konsumen.</td></tr>
<tr><td><strong>Sifat Produk</strong></td><td>Semua produk adalah barang surplus, kelebihan stok, atau barang retur yang dijual SEBAGAIMANA ADANYA, DI LOKASI SAAT INI, DENGAN SEMUA CACAT. Harga mencerminkan hal ini.</td></tr>
<tr><td><strong>Tanggung Jawab Pembeli</strong></td><td>Pembeli secara mandiri mengevaluasi produk, memastikan kepatuhan regulasi untuk kegiatan resale mereka sendiri, dan mematuhi semua hukum Indonesia yang berlaku.</td></tr>
<tr><td><strong>Tidak Ada Garansi</strong></td><td>Bulky.id menolak semua garansi tersirat termasuk kelayakan jual dan kesesuaian untuk tujuan tertentu.</td></tr>
<tr><td><strong>Pengalihan Risiko</strong></td><td>Risiko kehilangan beralih ke Pembeli saat serah terima di gudang atau pengiriman ke pihak pengangkut — mana yang lebih awal.</td></tr>
<tr><td><strong>Semua Penjualan Final</strong></td><td>Pengembalian tidak diizinkan kecuali untuk kekurangan jumlah material atau kesalahan Bulky.id yang terkonfirmasi, dilaporkan dalam 24 jam dengan dokumentasi.</td></tr>
<tr><td><strong>Ganti Rugi</strong></td><td>Pembeli mengganti rugi Bulky.id atas klaim yang timbul dari kegiatan resale atau penggunaan produk mereka, termasuk klaim dari pelanggan mereka sendiri.</td></tr>
<tr><td><strong>Batas Tanggung Gugat</strong></td><td>Total tanggung gugat Bulky.id untuk setiap klaim dibatasi pada harga pembelian yang dibayarkan untuk produk tertentu yang dipermasalahkan.</td></tr>
<tr><td><strong>Patungan</strong></td><td>Semua pihak dalam pembelian bersama bertanggung jawab secara tanggung renteng atas pesanan penuh dan semua kewajiban Syarat &amp; Ketentuan.</td></tr>
<tr><td><strong>Hukum yang Berlaku</strong></td><td>Hukum Indonesia berlaku. Sengketa diselesaikan melalui negosiasi, kemudian mediasi, kemudian Pengadilan Negeri Jakarta Selatan.</td></tr>
<tr><td><strong>Persetujuan Elektronik</strong></td><td>Penerimaan saat checkout (scroll dan centang) merupakan perjanjian elektronik yang mengikat secara hukum berdasarkan UU ITE Indonesia.</td></tr>
</tbody>
</table>
<p><strong>SYARAT DAN KETENTUAN LENGKAP HARUS DIBACA SELURUHNYA SEBELUM MEMBELI. RINGKASAN INI BUKAN NASIHAT HUKUM.</strong></p>$$,
    $$<p>The full Buyer Terms and Conditions of Purchase govern all transactions on Bulky.id. The summary below highlights the most important points. By purchasing on Bulky.id, you agree to be bound by the full Terms and Conditions — please read them in their entirety before your first purchase.</p>
<p>The full Buyer Terms and Conditions (Version 3.0, effective 1 June 2026) must be reviewed and accepted at checkout. The below is a summary only and does not replace the full document.</p>
<h2>Key Principles</h2>
<table>
<thead><tr><th>Principle</th><th>Detail</th></tr></thead>
<tbody>
<tr><td><strong>Platform Purpose</strong></td><td>B2B wholesale re-commerce platform for resellers, traders, and SMEs — not a consumer marketplace.</td></tr>
<tr><td><strong>Product Nature</strong></td><td>All products are surplus, overstock, or returned goods sold on an AS-IS, WHERE-IS, WITH ALL FAULTS basis. Pricing reflects this.</td></tr>
<tr><td><strong>Buyer Responsibility</strong></td><td>Buyers independently evaluate products, ensure regulatory compliance for their own resale activities, and comply with all applicable Indonesian law.</td></tr>
<tr><td><strong>No Warranties</strong></td><td>Bulky.id disclaims all implied warranties including merchantability and fitness for purpose.</td></tr>
<tr><td><strong>Risk Transfer</strong></td><td>Risk of loss transfers to Buyer upon handover at the warehouse or delivery to the shipping carrier — whichever is earlier.</td></tr>
<tr><td><strong>All Sales Final</strong></td><td>Returns not permitted except for material quantity shortfalls or confirmed Bulky.id errors, reported within 24 hours with documentation.</td></tr>
<tr><td><strong>Indemnification</strong></td><td>Buyers indemnify Bulky.id against claims arising from their resale or use of products, including claims from their own customers.</td></tr>
<tr><td><strong>Liability Cap</strong></td><td>Bulky.id's total liability for any claim is capped at the purchase price paid for the specific products in question.</td></tr>
<tr><td><strong>Patungan</strong></td><td>All parties in a co-purchase are jointly and severally liable for the full order and all T&amp;C obligations.</td></tr>
<tr><td><strong>Governing Law</strong></td><td>Indonesian law governs. Disputes resolved via negotiation, then mediation, then the South Jakarta District Court.</td></tr>
<tr><td><strong>Electronic Consent</strong></td><td>Acceptance at checkout (scroll and tick) constitutes a legally binding electronic agreement under Indonesian ITE Law.</td></tr>
</tbody>
</table>
<p><strong>THE FULL TERMS AND CONDITIONS MUST BE READ IN THEIR ENTIRETY BEFORE PURCHASING. THIS SUMMARY IS NOT LEGAL ADVICE.</strong></p>$$,
    6, true, NOW() + INTERVAL '5 second', NOW() + INTERVAL '5 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 7. Kebijakan Privasi / Privacy Policy
INSERT INTO dokumen_kebijakan (id, judul, judul_en, slug, slug_id, slug_en, konten, konten_en, urutan, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Kebijakan Privasi',
    'Privacy Policy',
    'kebijakan-privasi',
    'kebijakan-privasi',
    'privacy-policy',
    $$<p>Bulky.id berkomitmen untuk melindungi privasi dan data pribadi semua pengguna platform kami. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, berbagi, dan melindungi informasi pribadi Anda saat Anda mengakses dan menggunakan situs web dan aplikasi mobile Bulky.id.</p>
<p>Kebijakan Privasi ini mematuhi Undang-Undang Indonesia No. 27 Tahun 2022 tentang Perlindungan Data Pribadi (UU PDP), Undang-Undang No. 11 Tahun 2008 tentang Informasi dan Transaksi Elektronik (UU ITE) sebagaimana telah diubah, dan Peraturan Pemerintah No. 71 Tahun 2019.</p>
<h2>1. Data Pribadi yang Kami Kumpulkan</h2>
<h3>1.1 Data yang Anda Berikan</h3>
<ul>
<li><strong>Pendaftaran akun:</strong> nama lengkap, email, telepon, nama bisnis, dan alamat bisnis</li>
<li><strong>Verifikasi identitas:</strong> KTP yang diterbitkan pemerintah jika diperlukan untuk verifikasi akun</li>
<li><strong>Data transaksi:</strong> riwayat pesanan, catatan pembayaran, alamat pengiriman</li>
<li><strong>Komunikasi:</strong> pesan yang dikirim ke layanan pelanggan melalui WhatsApp, email, atau chat dalam aplikasi</li>
</ul>
<h3>1.2 Data yang Dikumpulkan Secara Otomatis</h3>
<ul>
<li><strong>Perangkat dan akses:</strong> alamat IP, jenis perangkat, jenis dan versi browser, sistem operasi</li>
<li><strong>Data penggunaan:</strong> halaman yang dikunjungi, waktu yang dihabiskan, klik, dan pola navigasi</li>
<li><strong>Data log:</strong> log akses server, log error, dan log transaksi</li>
<li><strong>Data cookie dan pelacakan:</strong> sebagaimana dijelaskan dalam Bagian 4</li>
</ul>
<h3>1.3 Data dari Pihak Ketiga</h3>
<ul>
<li><strong>Data pembayaran:</strong> status transaksi dari Xendit — kami tidak menerima atau menyimpan nomor kartu lengkap atau kredensial bank</li>
<li><strong>Data pengiriman:</strong> informasi pelacakan pengiriman dari penyedia logistik</li>
</ul>
<h2>2. Cara Kami Menggunakan Data Anda</h2>
<table>
<tbody>
<tr><td><strong>Manajemen Akun</strong></td><td>Untuk membuat, memverifikasi, dan memelihara akun Bulky.id Anda</td></tr>
<tr><td><strong>Pemrosesan Transaksi</strong></td><td>Untuk memproses pesanan, pembayaran, dan pemenuhan Anda</td></tr>
<tr><td><strong>Layanan Pelanggan</strong></td><td>Untuk merespons pertanyaan, keluhan, dan permintaan dukungan</td></tr>
<tr><td><strong>Peningkatan Platform</strong></td><td>Untuk menganalisis pola penggunaan dan meningkatkan fitur dan pengalaman</td></tr>
<tr><td><strong>Keamanan</strong></td><td>Untuk mendeteksi, mencegah, dan merespons penipuan dan akses tidak sah</td></tr>
<tr><td><strong>Kepatuhan Hukum</strong></td><td>Untuk mematuhi hukum Indonesia yang berlaku, perintah pengadilan, atau persyaratan regulasi</td></tr>
<tr><td><strong>Komunikasi</strong></td><td>Konfirmasi pesanan, notifikasi transaksi, dan — jika disetujui — komunikasi promosi</td></tr>
<tr><td><strong>Analitik</strong></td><td>Analitik agregat dan anonim yang tidak dapat mengidentifikasi Anda secara individual</td></tr>
</tbody>
</table>
<h2>3. Cara Kami Berbagi Data Anda</h2>
<p>Kami tidak menjual data pribadi Anda. Kami hanya membagikannya sebagai berikut:</p>
<ul>
<li><strong>Penyedia Layanan:</strong> Xendit (pembayaran), hosting cloud, dan mitra logistik — terikat oleh kewajiban kerahasiaan</li>
<li><strong>Persyaratan Hukum:</strong> Otoritas pemerintah, penegak hukum, atau pengadilan jika diperlukan oleh hukum Indonesia</li>
<li><strong>Pengalihan Bisnis:</strong> Dalam merger atau akuisisi — Anda akan diberitahu terlebih dahulu</li>
<li><strong>Dengan Persetujuan Anda:</strong> Dalam keadaan lain, hanya dengan persetujuan eksplisit Anda sebelumnya</li>
</ul>
<h2>4. Cookie</h2>
<table>
<tbody>
<tr><td><strong>Cookie Esensial</strong></td><td>Diperlukan untuk fungsi platform. Tidak dapat dinonaktifkan.</td></tr>
<tr><td><strong>Cookie Fungsional</strong></td><td>Mengingat preferensi Anda untuk meningkatkan pengalaman Anda.</td></tr>
<tr><td><strong>Cookie Analitik</strong></td><td>Data anonim tentang interaksi platform untuk membantu kami berkembang.</td></tr>
<tr><td><strong>Cookie Iklan Pihak Ketiga</strong></td><td>Ditempatkan oleh mitra periklanan. Bulky.id tidak mengontrol ini. Kelola melalui pengaturan browser.</td></tr>
</tbody>
</table>
<p>Anda dapat menonaktifkan cookie di pengaturan browser. Menonaktifkan cookie esensial dapat memengaruhi fungsionalitas platform.</p>
<h2>5. Penyimpanan dan Keamanan Data</h2>
<ul>
<li>Data disimpan di server aman di Indonesia atau infrastruktur cloud yang patuh</li>
<li>Langkah-langkah keamanan teknis dan organisasional: enkripsi, kontrol akses, tinjauan keamanan reguler</li>
<li>Pemberitahuan pelanggaran data pribadi kepada Anda dan otoritas dalam 14 x 24 jam setelah penemuan, sesuai UU PDP</li>
<li>Data disimpan selama masa akun Anda ditambah minimal 5 tahun setelah penutupan untuk kepatuhan hukum</li>
</ul>
<h2>6. Hak-Hak Anda</h2>
<table>
<tbody>
<tr><td><strong>Hak atas Informasi</strong></td><td>Untuk diberitahu tentang cara data Anda diproses</td></tr>
<tr><td><strong>Hak Akses</strong></td><td>Untuk meminta salinan data yang kami miliki tentang Anda</td></tr>
<tr><td><strong>Hak Koreksi</strong></td><td>Untuk meminta koreksi data yang tidak akurat atau tidak lengkap</td></tr>
<tr><td><strong>Hak Penghapusan</strong></td><td>Untuk meminta penghapusan data Anda, tunduk pada kewajiban penyimpanan hukum</td></tr>
<tr><td><strong>Hak Keberatan</strong></td><td>Untuk menolak pemrosesan untuk tujuan pemasaran langsung</td></tr>
<tr><td><strong>Hak Portabilitas</strong></td><td>Untuk menerima data Anda dalam format terstruktur yang dapat dibaca mesin</td></tr>
<tr><td><strong>Cabut Persetujuan</strong></td><td>Untuk mencabut persetujuan yang sebelumnya diberikan — tidak memengaruhi pemrosesan sah sebelumnya</td></tr>
</tbody>
</table>
<p>Untuk menggunakan hak-hak ini: email <a href="mailto:admin@bulky.id">admin@bulky.id</a> dengan subjek <em>Permintaan Privasi Data</em>. Kami merespons dalam 14 hari kerja.</p>
<h2>7. Privasi Anak</h2>
<p>Bulky.id adalah platform B2B untuk orang dewasa. Kami tidak dengan sengaja mengumpulkan data dari individu di bawah 18 tahun. Jika Anda percaya seorang anak telah memberikan data, segera hubungi <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
<h2>8. Tautan Pihak Ketiga</h2>
<p>Platform kami mungkin berisi tautan ke situs web pihak ketiga. Kebijakan Privasi ini hanya berlaku untuk Bulky.id. Tinjau kebijakan privasi pihak ketiga secara independen.</p>
<h2>9. Perubahan pada Kebijakan Ini</h2>
<p>Kami dapat memperbarui Kebijakan Privasi ini. Perubahan material dikomunikasikan melalui platform atau email terdaftar setidaknya 7 hari sebelum berlaku. Tanggal Terakhir Diperbarui di bagian atas mencerminkan revisi terbaru.</p>
<h2>10. Kontak — Masalah Privasi</h2>
<table>
<tbody>
<tr><td><strong>Email</strong></td><td><a href="mailto:admin@bulky.id">admin@bulky.id</a> — subjek: Privasi Data</td></tr>
<tr><td><strong>Telepon</strong></td><td>0811 833 164</td></tr>
<tr><td><strong>Kantor</strong></td><td>Plaza Mutiara Lt. 16, Jl. Lingkar Mega Kuningan Kav. E 1.2 No. 1 &amp; 2, Jakarta Selatan 12950</td></tr>
<tr><td><strong>Jam Operasional</strong></td><td>Senin hingga Sabtu, 09:00 hingga 18:00 WIB</td></tr>
</tbody>
</table>$$,
    $$<p>Bulky.id is committed to protecting the privacy and personal data of all users of our platform. This Privacy Policy describes how we collect, use, store, share, and protect your personal information when you access and use the Bulky.id website and mobile application.</p>
<p>This Privacy Policy complies with Indonesian Law No. 27 of 2022 on Personal Data Protection (UU PDP), Law No. 11 of 2008 on Electronic Information and Transactions (UU ITE) as amended, and Government Regulation No. 71 of 2019.</p>
<h2>1. Personal Data We Collect</h2>
<h3>1.1 Data You Provide</h3>
<ul>
<li><strong>Account registration:</strong> full name, email, phone, business name, and business address</li>
<li><strong>Identity verification:</strong> government-issued ID where required for account verification</li>
<li><strong>Transaction data:</strong> order history, payment records, shipping addresses</li>
<li><strong>Communications:</strong> messages sent to customer service via WhatsApp, email, or in-app chat</li>
</ul>
<h3>1.2 Data Collected Automatically</h3>
<ul>
<li><strong>Device and access:</strong> IP address, device type, browser type and version, operating system</li>
<li><strong>Usage data:</strong> pages visited, time spent, clicks, and navigation patterns</li>
<li><strong>Log data:</strong> server access logs, error logs, and transaction logs</li>
<li><strong>Cookie and tracking data:</strong> as described in Section 4</li>
</ul>
<h3>1.3 Data from Third Parties</h3>
<ul>
<li><strong>Payment data:</strong> transaction status from Xendit — we do not receive or store full card numbers or bank credentials</li>
<li><strong>Delivery data:</strong> shipment tracking information from logistics providers</li>
</ul>
<h2>2. How We Use Your Data</h2>
<table>
<tbody>
<tr><td><strong>Account Management</strong></td><td>To create, verify, and maintain your Bulky.id account</td></tr>
<tr><td><strong>Transaction Processing</strong></td><td>To process your orders, payments, and fulfillment</td></tr>
<tr><td><strong>Customer Service</strong></td><td>To respond to inquiries, complaints, and support requests</td></tr>
<tr><td><strong>Platform Improvement</strong></td><td>To analyse usage patterns and improve features and experience</td></tr>
<tr><td><strong>Security</strong></td><td>To detect, prevent, and respond to fraud and unauthorised access</td></tr>
<tr><td><strong>Legal Compliance</strong></td><td>To comply with applicable Indonesian law, court orders, or regulatory requirements</td></tr>
<tr><td><strong>Communications</strong></td><td>Order confirmations, transaction notifications, and — where consented — promotional communications</td></tr>
<tr><td><strong>Analytics</strong></td><td>Aggregated, anonymised analytics that cannot identify you individually</td></tr>
</tbody>
</table>
<h2>3. How We Share Your Data</h2>
<p>We do not sell your personal data. We share it only as follows:</p>
<ul>
<li><strong>Service Providers:</strong> Xendit (payment), cloud hosting, and logistics partners — bound by confidentiality obligations</li>
<li><strong>Legal Requirements:</strong> Government authorities, law enforcement, or courts where required by Indonesian law</li>
<li><strong>Business Transfers:</strong> In a merger or acquisition — you will be notified in advance</li>
<li><strong>With Your Consent:</strong> In any other circumstance, only with your express prior consent</li>
</ul>
<h2>4. Cookies</h2>
<table>
<tbody>
<tr><td><strong>Essential Cookies</strong></td><td>Required for platform function. Cannot be disabled.</td></tr>
<tr><td><strong>Functional Cookies</strong></td><td>Remember your preferences to improve your experience.</td></tr>
<tr><td><strong>Analytics Cookies</strong></td><td>Anonymised data about platform interaction to help us improve.</td></tr>
<tr><td><strong>Third-Party Ad Cookies</strong></td><td>Placed by advertising partners. Bulky.id does not control these. Manage via browser settings.</td></tr>
</tbody>
</table>
<p>You can disable cookies in your browser settings. Disabling essential cookies may affect platform functionality.</p>
<h2>5. Data Storage and Security</h2>
<ul>
<li>Data stored on secure servers in Indonesia or compliant cloud infrastructure</li>
<li>Technical and organisational security measures: encryption, access controls, regular security reviews</li>
<li>Personal data breach notification to you and authorities within 14 x 24 hours of discovery, per UU PDP</li>
<li>Data retained for the duration of your account plus a minimum of 5 years after closure for legal compliance</li>
</ul>
<h2>6. Your Rights</h2>
<table>
<tbody>
<tr><td><strong>Right to Information</strong></td><td>To be informed about how your data is processed</td></tr>
<tr><td><strong>Right to Access</strong></td><td>To request a copy of data we hold about you</td></tr>
<tr><td><strong>Right to Correction</strong></td><td>To request correction of inaccurate or incomplete data</td></tr>
<tr><td><strong>Right to Deletion</strong></td><td>To request deletion of your data, subject to legal retention obligations</td></tr>
<tr><td><strong>Right to Object</strong></td><td>To object to processing for direct marketing purposes</td></tr>
<tr><td><strong>Right to Portability</strong></td><td>To receive your data in a structured, machine-readable format</td></tr>
<tr><td><strong>Withdraw Consent</strong></td><td>To withdraw consent previously given — does not affect prior lawful processing</td></tr>
</tbody>
</table>
<p>To exercise these rights: email <a href="mailto:admin@bulky.id">admin@bulky.id</a> with subject <em>Data Privacy Request</em>. We respond within 14 business days.</p>
<h2>7. Children's Privacy</h2>
<p>Bulky.id is a B2B platform for adults. We do not knowingly collect data from individuals under 18. If you believe a minor has provided data, contact <a href="mailto:admin@bulky.id">admin@bulky.id</a> immediately.</p>
<h2>8. Third-Party Links</h2>
<p>Our platform may contain links to third-party websites. This Privacy Policy applies only to Bulky.id. Review third-party privacy policies independently.</p>
<h2>9. Changes to This Policy</h2>
<p>We may update this Privacy Policy. Material changes are communicated via the platform or registered email at least 7 days before taking effect. The Last Updated date at the top reflects the most recent revision.</p>
<h2>10. Contact — Privacy Matters</h2>
<table>
<tbody>
<tr><td><strong>Email</strong></td><td><a href="mailto:admin@bulky.id">admin@bulky.id</a> — subject: Data Privacy</td></tr>
<tr><td><strong>Phone</strong></td><td>0811 833 164</td></tr>
<tr><td><strong>Office</strong></td><td>Plaza Mutiara Lt. 16, Jl. Lingkar Mega Kuningan Kav. E 1.2 No. 1 &amp; 2, Jakarta Selatan 12950</td></tr>
<tr><td><strong>Business Hours</strong></td><td>Monday to Saturday, 09:00 to 18:00 WIB</td></tr>
</tbody>
</table>$$,
    7, true, NOW() + INTERVAL '6 second', NOW() + INTERVAL '6 second'
)
ON CONFLICT (slug) DO UPDATE SET
    judul = EXCLUDED.judul,
    judul_en = EXCLUDED.judul_en,
    slug_id = EXCLUDED.slug_id,
    slug_en = EXCLUDED.slug_en,
    konten = EXCLUDED.konten,
    konten_en = EXCLUDED.konten_en,
    urutan = EXCLUDED.urutan,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

COMMENT ON TABLE dokumen_kebijakan IS
'7 fixed policy pages with bilingual HTML content (ID + EN). Content can be edited via admin panel.';
