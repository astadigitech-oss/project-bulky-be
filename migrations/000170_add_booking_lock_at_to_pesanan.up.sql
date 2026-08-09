-- Tambah kolom lock anti double-booking yang TERPISAH dari booking_status.
--
-- Latar belakang: kolom booking_status (migration 000169) sebelumnya dipakai
-- ganda — sebagai lock claim ('IN_PROGRESS') oleh ClaimBooking DAN sebagai
-- penyimpan status provider oleh webhook Deliveree (locating_driver, dsb.).
-- Jika webhook masuk saat proses booking berjalan, lock tertimpa sehingga
-- proses lain bisa ikut memanggil API provider → berpotensi double booking.
--
-- Solusi: lock dipindah ke kolom khusus booking_lock_at (timestamptz) yang
-- TIDAK pernah disentuh webhook. booking_status kembali murni milik provider.

ALTER TABLE pesanan
    ADD COLUMN IF NOT EXISTS booking_lock_at TIMESTAMPTZ;

-- Bersihkan sisa lock IN_PROGRESS yang mungkin tertinggal dari versi lama
-- (mis. proses crash saat booking). Setelah migration ini, lock hanya berupa
-- timestamp pada booking_lock_at, bukan nilai di booking_status.
UPDATE pesanan
   SET booking_status = NULL
 WHERE booking_status = 'IN_PROGRESS';

COMMENT ON COLUMN pesanan.booking_lock_at IS 'Timestamp claim proses booking berlangsung (anti double-booking). NULL = bebas di-claim. Tidak disentuh webhook.';
