package constants

// Order statuses yang dianggap sebagai "transaksi sah" (bukan dibatalkan).
// Dipakai oleh seluruh endpoint dasbor agar definisi transaksi konsisten.
const (
	// TransaksiExcludedOrderStatus adalah status pesanan yang TIDAK dihitung
	// sebagai transaksi (order dibatalkan).
	TransaksiExcludedOrderStatus = "CANCELLED"

	// TransaksiPaidPaymentStatus adalah status pembayaran yang dianggap lunas
	// untuk metrik uang (revenue, total belanja).
	TransaksiPaidPaymentStatus = "PAID"

	// TransaksiCompletedOrderStatus adalah status pesanan yang dianggap
	// "benar-benar diselesaikan" untuk metrik penjualan selesai
	// (penjualan-per-buyer, user-transaction).
	TransaksiCompletedOrderStatus = "COMPLETED"
)
