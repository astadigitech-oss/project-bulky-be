-- Fix missing seed data from migration ordering issues
-- Applies: dokumen_kebijakan English content, warehouse coordinates, disclaimer slug_en

-- ============================================================
-- 1. Dokumen Kebijakan — judul_en & konten_en (actual English)
-- ============================================================
UPDATE dokumen_kebijakan SET judul_en = 'About Us', konten_en = $$<h2>Who We Are</h2>
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
</ul>$$ WHERE slug_id = 'tentang-kami';

UPDATE dokumen_kebijakan SET judul_en = 'How to Buy', konten_en = $$<p>Buying at Bulky.id is easy. Follow the steps below to place your first order. For enquiries, contact our team via WhatsApp or email.</p>
<h2>Main Purchase Flow</h2>
<h3>STEP 1 — Create Your Account</h3>
<p>Download the Bulky.id app or visit bulky.id. Click Register and complete your profile. Account activation is typically within 1 x 24 hours.</p>
<h3>STEP 2 — Browse Products</h3>
<p>Once logged in, browse the catalogue by category or use the search function. Each listing shows the condition grade, quantity, and price.</p>
<h3>STEP 3 — Add to Cart</h3>
<p>Select the product and quantity, then add to cart.</p>
<h3>STEP 4 — Choose Delivery or Pickup</h3>
<p>Select self-pickup at Cibinong Warehouse, or delivery via Deliveree or Forwarder.ai. Delivery costs are calculated at checkout.</p>
<h3>STEP 5 — Choose Payment Method</h3>
<p>Select Pay Directly or Pay Patungan with Friends to share the cost with partners.</p>
<h3>STEP 6 — Complete Payment</h3>
<p>Choose your payment channel and complete within the specified time limit. Unpaid orders will be automatically cancelled.</p>
<h2>Self-Pickup Information</h2>
<p>Warehouse: Jl. Raya Mayor Oking No. 62a, Cibinong, West Java. Monday to Saturday, 09:00 to 18:00 WIB.</p>$$ WHERE slug_id = 'cara-membeli';

UPDATE dokumen_kebijakan SET judul_en = 'About Payment', konten_en = $$<h2>Virtual Accounts (Bank Transfer)</h2>
<p>Supported banks: BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB.</p>
<h2>Credit &amp; Debit Cards</h2>
<p>Visa, Mastercard, JCB.</p>
<h2>E-Wallets</h2>
<p>OVO, DANA, GoPay, ShopeePay, LinkAja.</p>
<h2>QRIS</h2>
<p>Scan the QRIS code for instant payment from any mobile banking app or e-wallet that supports QRIS.</p>
<p>Note: No cash or COD. All payments are processed via Xendit.</p>$$ WHERE slug_id = 'tentang-pembayaran';

UPDATE dokumen_kebijakan SET judul_en = 'Contact Us', konten_en = $$<p>We are here to help. Contact us through:</p>
<h2>Warehouse Address</h2>
<p>Jl. Raya Mayor Oking Jaya Atmaja No.62a, Kel Cirimekar, Kec. Cibinong, Kabupaten Bogor, West Java 16918</p>
<h2>Phone / WhatsApp</h2>
<p><a href="https://wa.me/62811833164">+62811 833 164</a></p>
<h2>Email</h2>
<p><a href="mailto:admin@bulky.id">admin@bulky.id</a></p>
<h2>Operating Hours</h2>
<p>Monday to Saturday: 09:00 to 18:00 WIB. Sunday: Closed.</p>$$ WHERE slug_id = 'hubungi-kami';

UPDATE dokumen_kebijakan SET judul_en = 'Frequently Asked Questions', konten_en = $$<p>Can't find what you're looking for? Contact us via WhatsApp at <a href="https://wa.me/62811833164">0811 833 164</a> or email <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
<h2>About Bulky.id</h2>
<h3>What is Bulky.id?</h3>
<p>Bulky.id is Indonesia's B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory from e-commerce platforms, logistics providers, and corporations, then make it available to resellers, traders, and SMEs at wholesale prices.</p>
<h3>Is there a minimum order?</h3>
<p>There is no fixed minimum order value. Products are sold in pre-packed bulk units. The minimum purchase is the smallest available unit for each listing.</p>
<h2>About Products</h2>
<h3>What are the condition grades?</h3>
<p>New (factory-sealed), Like New (opened but unused), Used (previously used, functional), and Salvage (damaged, for parts or recycling).</p>
<h2>About Ordering and Payment</h2>
<h3>What payment methods are accepted?</h3>
<p>Virtual Account, Credit/Debit Cards, E-Wallets (OVO, DANA, GoPay, ShopeePay, LinkAja), and QRIS via Xendit.</p>
<h3>What is Patungan?</h3>
<p>Bulky.id's co-purchase split-payment feature. Two or more registered buyers split the cost of one order.</p>$$ WHERE slug_id = 'faq';

UPDATE dokumen_kebijakan SET judul_en = 'Terms and Conditions', konten_en = $$<p>The full Buyer Terms and Conditions govern all transactions on Bulky.id. By purchasing on Bulky.id, you agree to be bound by the full Terms and Conditions.</p>
<h2>Key Principles</h2>
<ul>
<li><strong>Platform Purpose:</strong> B2B wholesale re-commerce for resellers, traders, and SMEs — not a consumer marketplace.</li>
<li><strong>Product Nature:</strong> All products are surplus, overstock, or returned goods sold on an AS-IS basis.</li>
<li><strong>No Warranties:</strong> Bulky.id disclaims all implied warranties.</li>
<li><strong>All Sales Final:</strong> Returns not permitted except for material quantity shortfalls reported within 24 hours.</li>
<li><strong>Governing Law:</strong> Indonesian law. Disputes via South Jakarta District Court.</li>
</ul>$$ WHERE slug_id = 'syarat-ketentuan';

UPDATE dokumen_kebijakan SET judul_en = 'Privacy Policy', konten_en = $$<p>Bulky.id is committed to protecting the privacy and personal data of all users. This Privacy Policy describes how we collect, use, store, share, and protect your personal information.</p>
<h2>Personal Data We Collect</h2>
<ul>
<li>Account registration: full name, email, phone, business information</li>
<li>Transaction data: order history, payment records, shipping addresses</li>
<li>Usage data: pages visited, device and access information</li>
</ul>
<h2>Your Rights</h2>
<p>You have the right to access, correct, delete, and port your personal data. Contact <a href="mailto:admin@bulky.id">admin@bulky.id</a> with subject "Data Privacy Request".</p>$$ WHERE slug_id = 'kebijakan-privasi';

-- slug_en for dokumen_kebijakan
UPDATE dokumen_kebijakan SET slug_en = 'about-us'             WHERE slug_id = 'tentang-kami'         AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'how-to-buy'           WHERE slug_id = 'cara-membeli'          AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'about-payment'        WHERE slug_id = 'tentang-pembayaran'    AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'contact-us'           WHERE slug_id = 'hubungi-kami'          AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'faq'                  WHERE slug_id = 'faq'                   AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'terms-and-conditions' WHERE slug_id = 'syarat-ketentuan'      AND slug_en IS NULL;
UPDATE dokumen_kebijakan SET slug_en = 'privacy-policy'       WHERE slug_id = 'kebijakan-privasi'     AND slug_en IS NULL;

-- ============================================================
-- 2. Warehouse — latitude & longitude untuk warehouse-cibinong
-- ============================================================
UPDATE warehouse
SET latitude = -6.46958024, longitude = 106.85984984
WHERE slug = 'warehouse-cibinong' AND (latitude IS NULL OR longitude IS NULL);

-- ============================================================
-- 3. Disclaimer — slug_en
-- ============================================================
UPDATE disclaimer
SET slug_en = 'purchasing-disclaimer'
WHERE slug_id = 'disclaimer-pembelian' AND slug_en IS NULL;

-- ============================================================
-- 4. Dokumen Kebijakan — konten (restore Indonesian content)
-- ============================================================
UPDATE dokumen_kebijakan SET konten = $$<h2>Siapa Kami</h2>
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
<p>Bulky.id dirancang khusus untuk industri re-commerce Indonesia — bukan marketplace generik.</p>$$ WHERE slug_id = 'tentang-kami';
UPDATE dokumen_kebijakan SET konten = $$<p>Membeli di Bulky.id sangat mudah. Ikuti langkah-langkah di bawah ini untuk menempatkan pesanan pertama Anda. Untuk pertanyaan, hubungi tim kami via WhatsApp atau email.</p>
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
</table>$$ WHERE slug_id = 'cara-membeli';
UPDATE dokumen_kebijakan SET konten = $$<h2>Virtual Account (Transfer Bank)</h2>
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
<p><strong>Catatan:</strong> Tidak ada pembayaran tunai atau COD. Semua pembayaran diproses melalui Xendit.</p>$$ WHERE slug_id = 'tentang-pembayaran';
UPDATE dokumen_kebijakan SET konten = $$<p>Kami siap membantu Anda. Hubungi kami melalui:</p>
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
</ul>$$ WHERE slug_id = 'hubungi-kami';
UPDATE dokumen_kebijakan SET konten = $$<p>Tidak menemukan yang Anda cari? Hubungi kami via WhatsApp di <a href="https://wa.me/62811833164">0811 833 164</a> atau email <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
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
<p>Refund yang disetujui diproses ke metode pembayaran asli dalam 7 hingga 14 hari kerja.</p>$$ WHERE slug_id = 'faq';
UPDATE dokumen_kebijakan SET konten = $$<p>Syarat dan Ketentuan Pembelian lengkap mengatur semua transaksi di Bulky.id. Ringkasan di bawah ini menyoroti poin-poin terpenting. Dengan membeli di Bulky.id, Anda setuju untuk terikat oleh Syarat dan Ketentuan lengkap — harap baca seluruhnya sebelum pembelian pertama Anda.</p>
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
<p><strong>SYARAT DAN KETENTUAN LENGKAP HARUS DIBACA SELURUHNYA SEBELUM MEMBELI. RINGKASAN INI BUKAN NASIHAT HUKUM.</strong></p>$$ WHERE slug_id = 'syarat-ketentuan';
UPDATE dokumen_kebijakan SET konten = $$<p>Bulky.id berkomitmen untuk melindungi privasi dan data pribadi semua pengguna platform kami. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, berbagi, dan melindungi informasi pribadi Anda saat Anda mengakses dan menggunakan situs web dan aplikasi mobile Bulky.id.</p>
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
</table>$$ WHERE slug_id = 'kebijakan-privasi';
