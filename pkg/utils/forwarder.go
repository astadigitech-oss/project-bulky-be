package utils

import "strings"

// NormalizeKota menormalisasi teks kota dari Google Maps untuk digunakan
// sebagai kota_pattern pada tabel forwarder_city_mapping.
// Contoh: "Kota Yogyakarta" -> "yogyakarta", "Kabupaten Bogor" -> "bogor"
func NormalizeKota(kota string) string {
	kota = strings.ToLower(strings.TrimSpace(kota))
	prefixes := []string{"kota ", "kabupaten ", "kab. ", "kab "}
	for _, p := range prefixes {
		kota = strings.TrimPrefix(kota, p)
	}
	return strings.TrimSpace(kota)
}

// NormalizeKecamatan menormalisasi teks kecamatan dari Google Maps untuk digunakan
// sebagai kecamatan_pattern pada tabel forwarder_subdistrict_mapping.
// Contoh: "Kecamatan Cibinong" -> "cibinong"
func NormalizeKecamatan(kec string) string {
	kec = strings.ToLower(strings.TrimSpace(kec))
	prefixes := []string{"kecamatan ", "kec. ", "kec "}
	for _, p := range prefixes {
		kec = strings.TrimPrefix(kec, p)
	}
	return strings.TrimSpace(kec)
}
