package main

import (
	"fmt"
	"log"
	"os"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type DokumenKebijakanSeed struct {
	Judul  string
	Slug   string
	Konten string
	Urutan int
}

var dokumenSeeds = []DokumenKebijakanSeed{
	{
		Judul:  "Tentang Kami",
		Slug:   "tentang-kami",
		Urutan: 1,
		Konten: `<h1>Tentang Kami</h1>
<p>Bulky adalah platform e-commerce yang menyediakan berbagai produk elektronik berkualitas dengan harga terjangkau.</p>
<h2>Visi Kami</h2>
<p>Menjadi platform e-commerce terpercaya untuk produk elektronik di Indonesia.</p>
<h2>Misi Kami</h2>
<ul>
<li>Menyediakan produk elektronik berkualitas dengan harga kompetitif</li>
<li>Memberikan pelayanan terbaik kepada pelanggan</li>
<li>Membangun kepercayaan melalui transparansi dan integritas</li>
</ul>`,
	},
	{
		Judul:  "Cara Membeli",
		Slug:   "cara-membeli",
		Urutan: 2,
		Konten: `<h1>Cara Membeli</h1>
<p>Berikut adalah langkah-langkah untuk melakukan pembelian di Bulky:</p>
<ol>
<li><strong>Pilih Produk</strong> - Browse katalog produk dan pilih produk yang Anda inginkan</li>
<li><strong>Tambah ke Keranjang</strong> - Klik tombol "Tambah ke Keranjang" pada produk pilihan</li>
<li><strong>Review Keranjang</strong> - Periksa kembali produk di keranjang Anda</li>
<li><strong>Checkout</strong> - Isi data pengiriman dan pilih metode pembayaran</li>
<li><strong>Pembayaran</strong> - Lakukan pembayaran sesuai metode yang dipilih</li>
<li><strong>Konfirmasi</strong> - Tunggu konfirmasi pembayaran dari sistem</li>
<li><strong>Pengiriman</strong> - Produk akan dikirim setelah pembayaran dikonfirmasi</li>
</ol>`,
	},
	{
		Judul:  "Tentang Pembayaran",
		Slug:   "tentang-pembayaran",
		Urutan: 3,
		Konten: `<h1>Tentang Pembayaran</h1>
<p>Kami menyediakan berbagai metode pembayaran untuk kemudahan Anda:</p>
<h2>Bank Transfer / Virtual Account</h2>
<p>Transfer melalui bank BCA, Mandiri, BNI, BRI, dan bank lainnya.</p>
<h2>E-Wallet</h2>
<p>Pembayaran melalui GoPay, OVO, Dana, LinkAja, dan ShopeePay.</p>
<h2>QRIS</h2>
<p>Scan QRIS untuk pembayaran instan.</p>
<h2>Kartu Kredit</h2>
<p>Pembayaran menggunakan kartu kredit Visa, Mastercard, dan JCB.</p>
<p><strong>Catatan:</strong> Semua pembayaran akan diverifikasi dalam waktu maksimal 1x24 jam.</p>`,
	},
	{
		Judul:  "Hubungi Kami",
		Slug:   "hubungi-kami",
		Urutan: 4,
		Konten: `<h1>Hubungi Kami</h1>
<p>Kami siap membantu Anda. Hubungi kami melalui:</p>
<h2>Alamat</h2>
<p>Jl. Cilodong Raya No.89, Cilodong, Kec. Cilodong, Kota Depok, Jawa Barat 16414</p>
<h2>Telepon</h2>
<p>+62 811-833-164</p>
<h2>Email</h2>
<p>info@bulky.id</p>
<h2>Jam Operasional</h2>
<p>Senin - Sabtu: 09.00 - 18.00 WIB<br>Minggu: Tutup</p>
<h2>Media Sosial</h2>
<p>Instagram: @bulky.id<br>Facebook: Bulky Indonesia</p>`,
	},
	{
		Judul:  "Sering Ditanyakan",
		Slug:   "faq",
		Urutan: 5,
		Konten: `<h1>Pertanyaan yang Sering Ditanyakan (FAQ)</h1>
<h2>Bagaimana cara melacak pesanan saya?</h2>
<p>Anda dapat melacak pesanan melalui halaman "Pesanan Saya" di akun Anda. Nomor resi akan dikirimkan via email setelah pesanan dikirim.</p>
<h2>Berapa lama proses pengiriman?</h2>
<p>Estimasi pengiriman adalah 3-7 hari kerja tergantung lokasi tujuan.</p>
<h2>Apakah produk bergaransi?</h2>
<p>Ya, semua produk elektronik dilengkapi dengan garansi resmi dari distributor.</p>
<h2>Bagaimana cara mengembalikan produk?</h2>
<p>Produk dapat dikembalikan dalam 7 hari jika terdapat cacat atau kerusakan. Hubungi customer service kami untuk proses retur.</p>
<h2>Apakah bisa COD?</h2>
<p>Saat ini kami belum menyediakan layanan COD. Pembayaran dilakukan melalui transfer bank atau e-wallet.</p>`,
	},
	{
		Judul:  "Syarat dan Ketentuan",
		Slug:   "syarat-ketentuan",
		Urutan: 6,
		Konten: `<h1>Syarat dan Ketentuan</h1>
<p>Dengan menggunakan layanan Bulky, Anda menyetujui syarat dan ketentuan berikut:</p>
<h2>1. Penggunaan Layanan</h2>
<p>Anda harus berusia minimal 17 tahun atau memiliki izin dari orang tua/wali untuk menggunakan layanan kami.</p>
<h2>2. Akun Pengguna</h2>
<p>Anda bertanggung jawab untuk menjaga kerahasiaan akun dan password Anda.</p>
<h2>3. Pemesanan dan Pembayaran</h2>
<p>Semua pesanan harus dibayar sesuai dengan metode pembayaran yang dipilih. Pesanan akan diproses setelah pembayaran dikonfirmasi.</p>
<h2>4. Pengiriman</h2>
<p>Kami tidak bertanggung jawab atas keterlambatan pengiriman yang disebabkan oleh pihak ekspedisi atau force majeure.</p>
<h2>5. Pengembalian dan Penukaran</h2>
<p>Produk dapat dikembalikan dalam kondisi tertentu sesuai kebijakan retur kami.</p>
<h2>6. Harga</h2>
<p>Harga produk dapat berubah sewaktu-waktu tanpa pemberitahuan terlebih dahulu.</p>
<h2>7. Perubahan Syarat dan Ketentuan</h2>
<p>Kami berhak mengubah syarat dan ketentuan ini kapan saja. Perubahan akan diinformasikan melalui website.</p>`,
	},
	{
		Judul:  "Kebijakan Privasi",
		Slug:   "kebijakan-privasi",
		Urutan: 7,
		Konten: `<h1>Kebijakan Privasi</h1>
<p>Kami menghargai privasi Anda dan berkomitmen untuk melindungi data pribadi Anda.</p>
<h2>Informasi yang Kami Kumpulkan</h2>
<ul>
<li>Nama lengkap</li>
<li>Alamat email</li>
<li>Nomor telepon</li>
<li>Alamat pengiriman</li>
<li>Informasi pembayaran</li>
</ul>
<h2>Penggunaan Informasi</h2>
<p>Informasi yang kami kumpulkan digunakan untuk:</p>
<ul>
<li>Memproses pesanan Anda</li>
<li>Mengirimkan konfirmasi dan update pesanan</li>
<li>Meningkatkan layanan kami</li>
<li>Mengirimkan promosi dan penawaran (dengan persetujuan Anda)</li>
</ul>
<h2>Keamanan Data</h2>
<p>Kami menggunakan teknologi enkripsi dan sistem keamanan untuk melindungi data pribadi Anda.</p>
<h2>Berbagi Informasi</h2>
<p>Kami tidak akan menjual atau menyewakan informasi pribadi Anda kepada pihak ketiga tanpa persetujuan Anda, kecuali diwajibkan oleh hukum.</p>
<h2>Hak Anda</h2>
<p>Anda memiliki hak untuk mengakses, mengubah, atau menghapus data pribadi Anda. Hubungi kami untuk melakukan perubahan.</p>
<h2>Cookies</h2>
<p>Website kami menggunakan cookies untuk meningkatkan pengalaman pengguna. Anda dapat menonaktifkan cookies melalui pengaturan browser Anda.</p>`,
	},
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg)
	db := database.GetDB()

	log.Println("🌱 Seeding Dokumen Kebijakan...")

	for _, seed := range dokumenSeeds {
		// Check if exists
		var existing models.DokumenKebijakan
		if err := db.Where("slug = ?", seed.Slug).First(&existing).Error; err == nil {
			// Already exists, update content
			existing.Judul = seed.Judul
			existing.Konten = seed.Konten
			existing.Urutan = seed.Urutan
			existing.IsActive = true

			if err := db.Save(&existing).Error; err != nil {
				log.Printf("❌ Failed to update %s: %v", seed.Judul, err)
				os.Exit(1)
			}
			log.Printf("✓ Updated: %s", seed.Judul)
			continue
		}

		// Create new record
		dokumen := models.DokumenKebijakan{
			ID:       uuid.New(),
			Judul:    seed.Judul,
			Slug:     seed.Slug,
			Konten:   seed.Konten,
			Urutan:   seed.Urutan,
			IsActive: true,
		}

		if err := db.Create(&dokumen).Error; err != nil {
			log.Printf("❌ Failed to create %s: %v", seed.Judul, err)
			os.Exit(1)
		}

		log.Printf("✓ Created: %s", seed.Judul)
	}

	fmt.Println()
	log.Println("✅ Dokumen Kebijakan seeding completed!")
}
