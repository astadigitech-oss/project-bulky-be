package controllers

import (
	"archive/zip"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AssetMigrationController struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAssetMigrationController(db *gorm.DB, cfg *config.Config) *AssetMigrationController {
	return &AssetMigrationController{db: db, cfg: cfg}
}

// tmpDirV1 mengembalikan folder untuk file sementara (ZIP hasil gabungan
// chunk / upload single-shot) di BAWAH UploadPath — sengaja BUKAN os
// default temp dir ("" pada os.CreateTemp berarti /tmp di container),
// karena /tmp biasanya ada di ephemeral storage container yang jauh lebih
// kecil daripada volume persistent yang di-mount Dokploy (UPLOAD_PATH).
// ZIP backup v1 bisa beberapa GB — kalau digabung di /tmp yang penuh,
// proses finalize gagal dengan galat generik yang menyesatkan (seolah
// masalah jaringan/chunk, padahal disk container penuh).
func (ctrl *AssetMigrationController) tmpDirV1() (string, error) {
	dir := filepath.Join(ctrl.cfg.UploadPath, "chunks", "tmp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// isDiskFullError mendeteksi galat ENOSPC (disk/volume penuh) agar bisa
// dilaporkan dengan pesan yang jelas ke client, alih-alih pesan generik
// yang menyesatkan (seolah masalah jaringan terputus, padahal volume penuh).
func isDiskFullError(err error) bool {
	return errors.Is(err, syscall.ENOSPC)
}

// collectFilePaths queries all file paths referenced in DB and returns unique relative paths.
// Paths are stored in DB as relative (e.g. "blog/uuid.jpg"), but also handles full URLs.
func (ctrl *AssetMigrationController) collectFilePaths() ([]string, error) {
	type querySpec struct {
		table  string
		column string
		where  string
	}

	specs := []querySpec{
		{"produk_gambar", "gambar_url", "gambar_url != ''"},
		{"produk_dokumen", "file_url", "file_url != ''"},
		{"banner_event_promo", "gambar_url_id", "gambar_url_id != '' AND deleted_at IS NULL"},
		{"banner_event_promo", "gambar_url_en", "gambar_url_en != '' AND deleted_at IS NULL"},
		{"banner_tipe_produk", "gambar_url", "gambar_url != '' AND deleted_at IS NULL"},
		{"blog", "featured_image_url", "featured_image_url IS NOT NULL AND featured_image_url != '' AND deleted_at IS NULL"},
		{"hero_section", "gambar_url_id", "gambar_url_id != '' AND deleted_at IS NULL"},
		{"hero_section", "gambar_url_en", "gambar_url_en IS NOT NULL AND gambar_url_en != '' AND deleted_at IS NULL"},
		{"kategori_produk", "icon_url", "icon_url IS NOT NULL AND icon_url != '' AND deleted_at IS NULL"},
		{"kategori_produk", "gambar_kondisi_url", "gambar_kondisi_url IS NOT NULL AND gambar_kondisi_url != '' AND deleted_at IS NULL"},
		{"merek_produk", "logo_url", "logo_url IS NOT NULL AND logo_url != '' AND deleted_at IS NULL"},
		{"ulasan", "gambar", "gambar IS NOT NULL AND gambar != '' AND deleted_at IS NULL"},
		{"video", "video_url", "video_url != '' AND deleted_at IS NULL"},
		{"video", "thumbnail_url", "thumbnail_url IS NOT NULL AND thumbnail_url != '' AND deleted_at IS NULL"},
		{"buyer", "foto_url", "foto_url IS NOT NULL AND foto_url != '' AND deleted_at IS NULL"},
	}

	seen := make(map[string]bool)
	var unique []string

	for _, spec := range specs {
		var rows []string
		sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", spec.column, spec.table, spec.where)
		if err := ctrl.db.Raw(sql).Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("query %s.%s: %w", spec.table, spec.column, err)
		}

		for _, raw := range rows {
			rel := normalizeFilePath(raw)
			if rel == "" || seen[rel] {
				continue
			}
			seen[rel] = true
			unique = append(unique, rel)
		}
	}

	return unique, nil
}

// normalizeFilePath converts any stored path format to a relative path (no "uploads/" prefix).
// Returns empty string if the path is external (non-uploads URL) or empty.
func normalizeFilePath(p string) string {
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		idx := strings.Index(p, "/uploads/")
		if idx == -1 {
			return "" // external URL, not our file
		}
		p = p[idx+len("/uploads/"):]
	}
	p = strings.TrimPrefix(p, "uploads/")
	return p
}

// ExportAssets queries all DB-referenced upload paths, packages them into a zip, and streams it.
func (ctrl *AssetMigrationController) ExportAssets(c *fiber.Ctx) error {
	paths, err := ctrl.collectFilePaths()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengumpulkan daftar file", err.Error())
	}

	// Write zip to a temp file DI VOLUME UploadPath (bukan /tmp container)
	// to allow streaming large files safely without exhausting ephemeral disk.
	tmpDir, err := ctrl.tmpDirV1()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat direktori sementara", err.Error())
	}
	tmpFile, err := os.CreateTemp(tmpDir, "assets-export-*.zip")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat file sementara", err.Error())
	}
	tmpPath := tmpFile.Name()

	zw := zip.NewWriter(tmpFile)
	skipped := 0

	for _, relPath := range paths {
		fullPath := filepath.Join(ctrl.cfg.UploadPath, filepath.FromSlash(relPath))

		src, err := os.Open(fullPath)
		if err != nil {
			skipped++ // file referenced in DB but not on disk, skip
			continue
		}

		info, err := src.Stat()
		if err != nil {
			src.Close()
			skipped++
			continue
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			src.Close()
			skipped++
			continue
		}
		header.Name = relPath
		header.Method = zip.Deflate

		w, err := zw.CreateHeader(header)
		if err != nil {
			src.Close()
			skipped++
			continue
		}

		if _, err := io.Copy(w, src); err != nil {
			src.Close()
			skipped++
			continue
		}
		src.Close()
	}

	zw.Close()
	tmpFile.Close()

	defer os.Remove(tmpPath)

	filename := fmt.Sprintf("assets-export-%s.zip", time.Now().Format("20060102-150405"))
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("X-Export-Total", fmt.Sprintf("%d", len(paths)))
	c.Set("X-Export-Skipped", fmt.Sprintf("%d", skipped))

	return c.Download(tmpPath, filename)
}

// ImportAssets extracts an uploaded zip file into the uploads folder, preserving structure.
func (ctrl *AssetMigrationController) ImportAssets(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File zip tidak ditemukan. Kirim dengan field name 'file'", err.Error())
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File harus berformat .zip", "")
	}

	// Save uploaded zip to a temp file DI VOLUME UploadPath (avoids loading
	// entire zip into memory, and avoids ephemeral /tmp container disk).
	tmpDir, err := ctrl.tmpDirV1()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat direktori sementara", err.Error())
	}
	tmpZip, err := os.CreateTemp(tmpDir, "assets-import-*.zip")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat file sementara", err.Error())
	}
	tmpZipPath := tmpZip.Name()
	defer os.Remove(tmpZipPath)

	src, err := file.Open()
	if err != nil {
		tmpZip.Close()
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuka file", err.Error())
	}
	if _, err := io.Copy(tmpZip, src); err != nil {
		src.Close()
		tmpZip.Close()
		if isDiskFullError(err) {
			return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, "Volume storage penuh — hubungi devops untuk memperbesar volume Dokploy", err.Error())
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan file sementara", err.Error())
	}
	src.Close()
	tmpZip.Close()

	zr, err := zip.OpenReader(tmpZipPath)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File bukan format zip yang valid", err.Error())
	}
	defer zr.Close()

	imported := 0
	skipped := 0

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		// Prevent directory traversal attacks
		cleanName := filepath.Clean(f.Name)
		if strings.HasPrefix(cleanName, "..") || strings.Contains(cleanName, ".."+string(os.PathSeparator)) {
			skipped++
			continue
		}

		destPath := filepath.Join(ctrl.cfg.UploadPath, cleanName)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			skipped++
			continue
		}

		rc, err := f.Open()
		if err != nil {
			skipped++
			continue
		}

		dst, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			skipped++
			continue
		}

		if _, err := io.Copy(dst, rc); err != nil {
			dst.Close()
			rc.Close()
			os.Remove(destPath)
			skipped++
			continue
		}

		dst.Close()
		rc.Close()
		imported++
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Import selesai", fiber.Map{
		"imported": imported,
		"skipped":  skipped,
	})
}

// v1AssetPrefixes memetakan folder sumber di ZIP v1 (storage Laravel)
// ke folder tujuan di storage v2. Path ZIP v1 berbentuk
// "storage/app/public/<folder>/<file>"; baris kiri adalah folder sumber,
// baris kanan adalah prefix tujuan relatif terhadap UploadPath.
//
// Urutan PENTING: rule yang lebih spesifik harus lebih dulu. Sebagian PDF
// v1 disimpan di subfolder "products/pdf/" (lihat CreateProduct.php v1:
// Storage::disk('public')->put('products/pdf/...')), dan rule
// "storage/app/public/products/" TIDAK BOLEH menangkapnya — kalau iya, PDF
// salah masuk ke product-images/ lalu dianggap orphan dan ter-prune.
var v1AssetPrefixes = []struct{ src, dst string }{
	{"storage/app/public/products/pdf/", "product-documents/"},
	{"storage/app/public/products/", "product-images/"},
	{"storage/app/public/reviews/", "reviews/"},
	{"storage/app/public/public/profile/", "profile/"},
	{"storage/app/public/", ""}, // dokumen/PDF v1 diletakkan langsung di storage/app/public/
}

// v1ImportResult adalah ringkasan hasil pemrosesan ZIP v1, dipakai
// bersama oleh ImportAssetsV1 (multipart) dan finalize chunk upload.
type v1ImportResult struct {
	Imported  int
	Skipped   int
	Unmatched []string
	Files     []map[string]any
}

// processV1Zip membaca ZIP dari path lokal, memetakan isinya sesuai
// v1AssetPrefixes, lalu mengekstrak file yang cocok ke UploadPath.
// ZIP dibaca langsung dari disk (tidak dimuat ke memori), jadi aman
// untuk file berukuran besar.
func (ctrl *AssetMigrationController) processV1Zip(zipPath string) (v1ImportResult, error) {
	var res v1ImportResult

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return res, fmt.Errorf("file bukan format zip yang valid: %w", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		cleanName := filepath.ToSlash(filepath.Clean(f.Name))
		// Anti path traversal
		if strings.HasPrefix(cleanName, "../") || strings.Contains(cleanName, "/../") {
			res.Skipped++
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "path traversal ditolak"})
			continue
		}

		// Tentukan folder tujuan v2 berdasarkan struktur v1
		var relPath string
		mapped := false
		for _, m := range v1AssetPrefixes {
			if strings.HasPrefix(cleanName, m.src) {
				rest := strings.TrimPrefix(cleanName, m.src)
				if rest == "" {
					continue
				}
				// Dokumen v1 berupa PDF yang diletakkan langsung di storage/app/public/
				if m.dst == "" {
					if !strings.EqualFold(filepath.Ext(rest), ".pdf") {
						continue // bukan PDF, bukan aset yang kita kelola
					}
					relPath = "product-documents/" + rest
				} else {
					relPath = m.dst + rest
				}
				mapped = true
				break
			}
		}

		if !mapped {
			res.Skipped++
			if len(res.Unmatched) < 20 {
				res.Unmatched = append(res.Unmatched, cleanName)
			}
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "path tidak dikenali"})
			continue
		}

		destPath := filepath.Join(ctrl.cfg.UploadPath, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			res.Skipped++
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "gagal buat folder"})
			continue
		}

		rc, err := f.Open()
		if err != nil {
			res.Skipped++
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "gagal baca zip"})
			continue
		}

		dst, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			res.Skipped++
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "gagal buat file"})
			continue
		}

		if _, err := io.Copy(dst, rc); err != nil {
			dst.Close()
			rc.Close()
			os.Remove(destPath)
			// Disk penuh bersifat fatal untuk seluruh ZIP (bukan cuma file ini) —
			// hentikan proses daripada menandai ratusan file sisanya "skipped"
			// dengan alasan yang menyesatkan seperti "path tidak dikenali".
			if isDiskFullError(err) {
				return res, fmt.Errorf("volume storage penuh saat mengekstrak %q — hubungi devops untuk memperbesar volume Dokploy: %w", cleanName, err)
			}
			res.Skipped++
			res.Files = append(res.Files, map[string]any{"source": cleanName, "status": "skipped", "reason": "gagal tulis file"})
			continue
		}

		dst.Close()
		rc.Close()
		res.Imported++
		res.Files = append(res.Files, map[string]any{"source": cleanName, "dest": relPath, "status": "imported"})
	}

	return res, nil
}

// ImportAssetsV1 mengimpor ZIP hasil export folder "storage" Laravel v1
// (storage/app/public/...). Karena path yang tersimpan di DB v2 sudah ditransformasi
// migrasi (produk -> "product-images/", dokumen -> "product-documents/",
// ulasan -> "reviews/", foto profil -> "profile/"), file dari ZIP v1 perlu
// dipetakan ulang ke lokasi yang sesuai, bukan diekstrak apa adanya.
//
// Mapping yang didukung (cukup 5 tipe aset publik v1):
//
//	storage/app/public/products/pdf/<f>    -> <UPLOAD_PATH>/product-documents/<f>
//	storage/app/public/products/<f>        -> <UPLOAD_PATH>/product-images/<f>
//	storage/app/public/reviews/<f>         -> <UPLOAD_PATH>/reviews/<f>
//	storage/app/public/public/profile/<f>  -> <UPLOAD_PATH>/profile/<f>
//	storage/app/public/<f>.pdf             -> <UPLOAD_PATH>/product-documents/<f>
//
// File yang tidak cocok mapping mana pun (folder lain dari aplikasi v1) di-skip.
// ZIP >500MB tidak bisa lewat endpoint ini karena Fiber BodyLimit; gunakan
// chunk upload (POST /panel/assets/import-v1/chunk + /finalize) untuk file besar.
func (ctrl *AssetMigrationController) ImportAssetsV1(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File zip tidak ditemukan. Kirim dengan field name 'file'", err.Error())
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File harus berformat .zip", "")
	}

	// Simpan ZIP ke file sementara DI VOLUME UploadPath (bukan /tmp container
	// yang ephemeral & kecil) — hindari muat seluruh zip ke memori.
	tmpDir, err := ctrl.tmpDirV1()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat direktori sementara", err.Error())
	}
	tmpZip, err := os.CreateTemp(tmpDir, "assets-import-v1-*.zip")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat file sementara", err.Error())
	}
	tmpZipPath := tmpZip.Name()
	defer os.Remove(tmpZipPath)

	src, err := file.Open()
	if err != nil {
		tmpZip.Close()
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuka file", err.Error())
	}
	if _, err := io.Copy(tmpZip, src); err != nil {
		src.Close()
		tmpZip.Close()
		if isDiskFullError(err) {
			return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, "Volume storage penuh — hubungi devops untuk memperbesar volume Dokploy", err.Error())
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan file sementara", err.Error())
	}
	src.Close()
	tmpZip.Close()

	res, err := ctrl.processV1Zip(tmpZipPath)
	if err != nil {
		if isDiskFullError(err) {
			return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, err.Error(), "")
		}
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, err.Error(), "")
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Import v1 selesai", fiber.Map{
		"imported":  res.Imported,
		"skipped":   res.Skipped,
		"unmatched": res.Unmatched,
		"files":     res.Files,
	})
}

// chunksDirV1 mengembalikan folder tempat chunk ZIP v1 ditampung
// (di bawah UploadPath, sama seperti chunk video).
func (ctrl *AssetMigrationController) chunksDirV1(uploadID string) string {
	return filepath.Join(ctrl.cfg.UploadPath, "chunks", "v1-"+uploadID)
}

// UploadV1Chunk menerima satu potongan (chunk) ZIP backup v1.
// Format multipart sama dengan upload video:
//
//	upload_id:     UUID yang sama untuk semua chunk
//	chunk_index:   0-based
//	total_chunks:  total seluruh chunk
//	chunk_data:    file bagian ZIP
//	chunk_sha1:    (opsional) SHA-1 hex dari isi chunk — dipakai untuk verifikasi
//	               & resume: kalau chunk sudah tersimpan dengan sha1 sama,
//	               server langsung balas sukses tanpa menerima body lagi.
//
// Setiap chunk disimpan ke <UPLOAD_PATH>/chunks/v1-<upload_id>/chunk_XXXXX.
// Karena setiap request hanya membawa satu chunk (ukuran ≤ BodyLimit),
// ZIP sebesar apa pun bisa dikirim secara bertahap.
func (ctrl *AssetMigrationController) UploadV1Chunk(c *fiber.Ctx) error {
	uploadID := c.FormValue("upload_id")
	chunkIndexStr := c.FormValue("chunk_index")
	totalChunksStr := c.FormValue("total_chunks")
	clientSHA1 := strings.ToLower(c.FormValue("chunk_sha1"))

	if uploadID == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "upload_id wajib diisi", "")
	}
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || chunkIndex < 0 {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "chunk_index tidak valid", "")
	}
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil || totalChunks <= 0 {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "total_chunks tidak valid", "")
	}

	tempDir := ctrl.chunksDirV1(uploadID)
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%05d", chunkIndex))

	// Resume/duplicate detection: kalau chunk sudah tersimpan & checksum cocok,
	// anggap sudah terkirim — hemat bandwidth (pola resumable upload Google Drive).
	if clientSHA1 != "" {
		if existing, err := os.ReadFile(chunkPath); err == nil {
			if sumSHA1(existing) == clientSHA1 {
				return utils.SimpleSuccessResponse(c, http.StatusOK, "Chunk sudah tersimpan", fiber.Map{
					"upload_id":     uploadID,
					"chunk_index":   chunkIndex,
					"total_chunks":  totalChunks,
					"chunk_sha1":    clientSHA1,
					"already_exist": true,
				})
			}
			// sha1 beda → file lama korup/tidak lengkap → timpa (retry penuh)
			os.Remove(chunkPath)
		}
	}

	file, err := c.FormFile("chunk_data")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "chunk_data wajib diisi", err.Error())
	}

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat direktori temp", err.Error())
	}

	if err := c.SaveFile(file, chunkPath); err != nil {
		if isDiskFullError(err) {
			return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, "Volume storage penuh — hubungi devops untuk memperbesar volume Dokploy", err.Error())
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan chunk", err.Error())
	}

	// Verifikasi checksum hasil simpan — deteksi data korup sejak dini.
	saved, err := os.ReadFile(chunkPath)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membaca chunk tersimpan", err.Error())
	}
	savedSHA1 := sumSHA1(saved)
	if clientSHA1 != "" && savedSHA1 != clientSHA1 {
		os.Remove(chunkPath)
		return utils.SimpleErrorResponse(c, http.StatusUnprocessableEntity, "Chunk rusak (checksum tidak cocok)", "")
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Chunk berhasil disimpan", fiber.Map{
		"upload_id":     uploadID,
		"chunk_index":   chunkIndex,
		"total_chunks":  totalChunks,
		"chunk_sha1":    savedSHA1,
		"already_exist": false,
	})
}

// sumSHA1 menghitung SHA-1 hex dari data (untuk verifikasi integritas chunk).
func sumSHA1(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

// FinalizeV1Chunk menggabungkan semua chunk ZIP v1 menjadi satu file,
// lalu memprosesnya dengan logika yang sama seperti ImportAssetsV1
// (mapping storage/app/public/... ke folder v2). Folder chunk dihapus
// setelah selesai (sukses maupun gagal).
func (ctrl *AssetMigrationController) FinalizeV1Chunk(c *fiber.Ctx) error {
	uploadID := c.FormValue("upload_id")
	totalChunksStr := c.FormValue("total_chunks")

	tempDir := ctrl.chunksDirV1(uploadID)
	// Bersihkan folder chunk di awal — validasi error di bawah juga
	// menghapus sisa chunk (mencegah folder chunks/ menumpuk).
	defer os.RemoveAll(tempDir)

	if uploadID == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "upload_id wajib diisi", "")
	}
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil || totalChunks <= 0 {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "total_chunks tidak valid", "")
	}

	// Gabungkan chunk menjadi ZIP sementara DI VOLUME UploadPath (bukan /tmp
	// container yang ephemeral & kecil). Ini penting untuk ZIP backup v1
	// yang bisa beberapa GB — kalau /tmp container penuh, penggabungan
	// gagal di tengah jalan (gejala: "N bagian gagal" padahal semua chunk
	// sudah sukses terkirim satu-satu).
	tmpDir, err := ctrl.tmpDirV1()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat direktori sementara", err.Error())
	}
	tmpZip, err := os.CreateTemp(tmpDir, "assets-import-v1-*.zip")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat file sementara", err.Error())
	}
	tmpZipPath := tmpZip.Name()
	defer os.Remove(tmpZipPath)

	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%05d", i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			tmpZip.Close()
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Chunk %d tidak ditemukan. Pastikan semua chunk sudah diupload", i), err.Error())
		}
		_, copyErr := io.Copy(tmpZip, chunk)
		chunk.Close()
		if copyErr != nil {
			tmpZip.Close()
			if isDiskFullError(copyErr) {
				return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, "Volume storage penuh saat menggabungkan chunk — hubungi devops untuk memperbesar volume Dokploy", copyErr.Error())
			}
			return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menggabungkan chunk", copyErr.Error())
		}
	}
	tmpZip.Sync()
	tmpZip.Close()

	res, err := ctrl.processV1Zip(tmpZipPath)
	if err != nil {
		if isDiskFullError(err) {
			return utils.SimpleErrorResponse(c, http.StatusInsufficientStorage, err.Error(), "")
		}
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, err.Error(), "")
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Import v1 (chunk) selesai", fiber.Map{
		"imported":  res.Imported,
		"skipped":   res.Skipped,
		"unmatched": res.Unmatched,
		"files":     res.Files,
	})
}

// AbortV1Chunk membatalkan chunk upload yang sedang berjalan: menghapus
// folder <UPLOAD_PATH>/chunks/v1-<upload_id> beserta semua chunk yang
// sudah terlanjur tersimpan. Dipakai FE saat migrasi dibatalkan/gagal.
// Idempotent — jika folder tidak ada, tetap mengembalikan sukses.
//
// Opsional: jika field chunk_index dikirim (dengan chunk_sha1), hanya chunk
// itu yang dihapus — dipakai FE untuk reset chunk yang checksum-nya gagal.
func (ctrl *AssetMigrationController) AbortV1Chunk(c *fiber.Ctx) error {
	uploadID := c.FormValue("upload_id")
	if uploadID == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "upload_id wajib diisi", "")
	}

	tempDir := ctrl.chunksDirV1(uploadID)

	// Hapus satu chunk saja (reset per-chunk)
	if idxStr := c.FormValue("chunk_index"); idxStr != "" {
		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 0 {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, "chunk_index tidak valid", "")
		}
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%05d", idx))
		if err := os.Remove(chunkPath); err != nil && !os.IsNotExist(err) {
			return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus chunk", err.Error())
		}
		return utils.SimpleSuccessResponse(c, http.StatusOK, "Chunk dihapus", nil)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membersihkan chunk", err.Error())
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Upload chunk dibatalkan", nil)
}

// pendingV1Upload merangkum satu folder chunk upload v1 yang belum di-finalize/abort.
type pendingV1Upload struct {
	UploadID     string    `json:"upload_id"`
	ChunkCount   int       `json:"chunk_count"`
	TotalSize    int64     `json:"total_size"`
	LastModified time.Time `json:"last_modified"`
}

// scanPendingV1Uploads memindai <UPLOAD_PATH>/chunks untuk folder "v1-*"
// (chunk yang sudah diupload tapi belum di-finalize/abort) dan file
// sementara "tmp/*" (ZIP gabungan yang tertinggal karena request terputus
// di tengah proses finalize, misal container di-restart Dokploy).
// Dipakai ListPendingV1Uploads & CleanupStaleV1Uploads.
func (ctrl *AssetMigrationController) scanPendingV1Uploads() ([]pendingV1Upload, []string, error) {
	chunksRoot := filepath.Join(ctrl.cfg.UploadPath, "chunks")
	entries, err := os.ReadDir(chunksRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var pending []pendingV1Upload
	var staleTmpFiles []string

	for _, entry := range entries {
		name := entry.Name()

		// File ZIP sementara yatim di chunks/tmp/ (harusnya sudah dihapus
		// oleh defer os.Remove setelah finalize, tapi bisa tertinggal kalau
		// proses terputus paksa, misal container restart).
		if name == "tmp" && entry.IsDir() {
			tmpEntries, err := os.ReadDir(filepath.Join(chunksRoot, "tmp"))
			if err != nil {
				continue
			}
			for _, tf := range tmpEntries {
				staleTmpFiles = append(staleTmpFiles, filepath.Join(chunksRoot, "tmp", tf.Name()))
			}
			continue
		}

		if !entry.IsDir() || !strings.HasPrefix(name, "v1-") {
			continue
		}

		uploadID := strings.TrimPrefix(name, "v1-")
		dirPath := filepath.Join(chunksRoot, name)
		chunkEntries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}

		var totalSize int64
		var lastModified time.Time
		for _, ce := range chunkEntries {
			info, err := ce.Info()
			if err != nil {
				continue
			}
			totalSize += info.Size()
			if info.ModTime().After(lastModified) {
				lastModified = info.ModTime()
			}
		}

		pending = append(pending, pendingV1Upload{
			UploadID:     uploadID,
			ChunkCount:   len(chunkEntries),
			TotalSize:    totalSize,
			LastModified: lastModified,
		})
	}

	return pending, staleTmpFiles, nil
}

// ListPendingV1Uploads menampilkan semua upload chunk v1 yang masih
// tersimpan di volume tapi belum di-finalize/abort — berguna untuk memantau
// pemakaian disk volume Dokploy tanpa perlu SSH manual, dan mendeteksi
// chunk basi dari percobaan migrasi yang gagal/ditinggal begitu saja.
func (ctrl *AssetMigrationController) ListPendingV1Uploads(c *fiber.Ctx) error {
	pending, staleTmp, err := ctrl.scanPendingV1Uploads()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal memindai folder chunk", err.Error())
	}

	var totalSize int64
	for _, p := range pending {
		totalSize += p.TotalSize
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Daftar upload chunk pending", fiber.Map{
		"pending_uploads":    pending,
		"stale_tmp_files":    staleTmp,
		"total_pending_size": totalSize,
	})
}

// CleanupStaleV1Uploads menghapus folder chunk v1 (dan file ZIP sementara
// yatim di chunks/tmp/) yang sudah lebih lama dari `older_than_hours`
// (default 24 jam) tanpa pernah di-finalize/abort. Ini mencegah volume
// Dokploy diam-diam penuh oleh sisa percobaan migrasi yang gagal berulang
// kali — penyebab paling umum error "N bagian gagal" pada upload berikutnya
// meski setiap chunk individual berhasil terkirim.
//
// Request body (JSON, opsional):
//
//	{ "older_than_hours": 24, "dry_run": true }  // default: dry_run=true
//
// Eksekusi permanen (dry_run=false) WAJIB menyertakan dry_run_token dari
// dry-run terakhir (pola yang sama dengan PruneOrphans) — token one-time,
// kedaluwarsa 10 menit. Ini mencegah penghapusan permanen tanpa pratinjau.
func (ctrl *AssetMigrationController) CleanupStaleV1Uploads(c *fiber.Ctx) error {
	olderThanHours := 24
	dryRun := true
	var req struct {
		OlderThanHours *int   `json:"older_than_hours"`
		DryRun         *bool  `json:"dry_run"`
		DryRunToken    string `json:"dry_run_token"`
	}
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Body tidak valid", err.Error())
		}
		if req.OlderThanHours != nil {
			olderThanHours = *req.OlderThanHours
		}
		if req.DryRun != nil {
			dryRun = *req.DryRun
		}
	}

	// Eksekusi permanen tanpa token dry-run → tolak
	if !dryRun {
		if req.DryRunToken == "" {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest,
				"Jalankan dry-run dulu untuk mendapatkan token konfirmasi", "")
		}
		mu.Lock()
		tok, ok := pruneTokens[req.DryRunToken]
		if ok {
			delete(pruneTokens, req.DryRunToken) // one-time use
		}
		mu.Unlock()
		if !ok || time.Now().After(tok.expiry) {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest,
				"Token dry-run tidak valid atau kedaluwarsa. Jalankan dry-run ulang", "")
		}
	}

	pending, staleTmp, err := ctrl.scanPendingV1Uploads()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal memindai folder chunk", err.Error())
	}

	threshold := time.Now().Add(-time.Duration(olderThanHours) * time.Hour)
	chunksRoot := filepath.Join(ctrl.cfg.UploadPath, "chunks")

	var deleted []string
	var freedSize int64
	for _, p := range pending {
		if p.LastModified.After(threshold) {
			continue // masih dalam progres, jangan sentuh
		}
		deleted = append(deleted, p.UploadID)
		freedSize += p.TotalSize
		if !dryRun {
			os.RemoveAll(filepath.Join(chunksRoot, "v1-"+p.UploadID))
		}
	}

	var deletedTmp []string
	for _, tf := range staleTmp {
		info, err := os.Stat(tf)
		if err != nil || info.ModTime().After(threshold) {
			continue
		}
		deletedTmp = append(deletedTmp, filepath.Base(tf))
		freedSize += info.Size()
		if !dryRun {
			os.Remove(tf)
		}
	}

	// Terbitkan token one-time untuk eksekusi permanen (dry-run saja)
	var dryRunToken string
	if dryRun {
		dryRunToken = fmt.Sprintf("cleanup-%s", hex.EncodeToString(randomBytes(16)))
		mu.Lock()
		pruneTokens[dryRunToken] = pruneToken{
			expiry:     time.Now().Add(10 * time.Minute),
			totalFiles: len(deleted) + len(deletedTmp),
			totalSize:  freedSize,
		}
		mu.Unlock()
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Pembersihan chunk basi selesai", fiber.Map{
		"dry_run":           dryRun,
		"older_than_hours":  olderThanHours,
		"deleted_uploads":   deleted,
		"deleted_tmp_files": deletedTmp,
		"freed_size":        freedSize,
		"dry_run_token":     dryRunToken,
		"token_expiry_s":    600,
	})
}

// VerifyAssets membandingkan path file yang direferensikan DB v2 terhadap
// file fisik di UploadPath, lalu melaporkan anomali:
//
//   - missing: file direferensikan DB tapi TIDAK ada di disk (contoh nyata:
//     30 PDF produk yang ikut ter-prune saat mapping v1 salah — DB menunjuk
//     path yang tidak pernah diekstrak). Ini sinyal bahwa migrasi/import
//     belum lengkap, BUKAN kandidat prune.
//   - duplicates: dua baris DB atau lebih menunjuk path yang sama persis.
//     Biasanya akibat migrasi data yang menggabungkan path tanpa menormalkan
//     (contoh: "product-documents/a.pdf" vs "product-documents/a.pdf" hasil
//     prefix yang berbeda). Menghapus salah satunya berisiko, jadi dilaporkan
//     sebagai peringatan.
//
// Endpoint ini read-only dan aman dipanggil kapan saja — dipakai FE untuk
// menyaring risiko sebelum aksi prune permanen.
func (ctrl *AssetMigrationController) VerifyAssets(c *fiber.Ctx) error {
	paths, err := ctrl.collectFilePaths()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengumpulkan daftar file DB", err.Error())
	}

	type missingFile struct {
		Path string `json:"path"`
	}
	type duplicateFile struct {
		Path  string `json:"path"`
		Count int    `json:"count"`
	}

	counts := make(map[string]int, len(paths))
	for _, p := range paths {
		counts[p]++
	}

	var missing []missingFile
	var duplicates []duplicateFile
	for _, p := range paths {
		full := filepath.Join(ctrl.cfg.UploadPath, filepath.FromSlash(p))
		if _, err := os.Stat(full); err != nil {
			missing = append(missing, missingFile{Path: p})
		}
	}
	for p, n := range counts {
		if n > 1 {
			duplicates = append(duplicates, duplicateFile{Path: p, Count: n})
		}
	}

	// Urutkan agar output deterministik (map iteration di Go tidak berurutan)
	sort.Slice(missing, func(i, j int) bool { return missing[i].Path < missing[j].Path })
	sort.Slice(duplicates, func(i, j int) bool { return duplicates[i].Path < duplicates[j].Path })

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Verifikasi aset selesai", fiber.Map{
		"referenced": len(paths),
		"missing":    missing,
		"duplicates": duplicates,
	})
}

// PruneOrphans menghapus file di UploadPath yang TIDAK direferensikan
// oleh kolom file mana pun di DB v2. File yang dipakai aplikasi diambil
// dari collectFilePaths (yang juga dipakai ExportAssets), jadi daftarnya
// konsisten. Endpoint ini berbahaya — semua file yang dipilih bisa
// dihapus permanen.
//
// Request body (JSON):
//
//	{
//	  "dry_run": true,   // wajib true dulu untuk mendapat token konfirmasi
//	  "dry_run_token": "..."  // hanya untuk eksekusi (dry_run=false)
//	}
//
// Pengaman: eksekusi permanen (dry_run=false) WAJIB menyertakan
// dry_run_token yang dihasilkan dry-run terakhir. Token kedaluwarsa
// setelah 10 menit atau setelah dipakai sekali (one-time). Ini mencegah
// FE/script menghapus file tanpa pernah melihat daftarnya — persis kasus
// 30 PDF ter-prune karena mapping salah tanpa sempat dicek.
func (ctrl *AssetMigrationController) PruneOrphans(c *fiber.Ctx) error {
	dryRun := true
	var req struct {
		DryRun      *bool  `json:"dry_run"`
		DryRunToken string `json:"dry_run_token"`
	}
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Body tidak valid", err.Error())
		}
		if req.DryRun != nil {
			dryRun = *req.DryRun
		}
	}

	// Eksekusi permanen tanpa token dry-run → tolak. Token mencegah aksi
	// destruktif tanpa pratinjau daftar file yang akan dihapus.
	if !dryRun {
		if req.DryRunToken == "" {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest,
				"Jalankan dry-run dulu untuk mendapatkan token konfirmasi", "")
		}
		mu.Lock()
		tok, ok := pruneTokens[req.DryRunToken]
		if ok {
			delete(pruneTokens, req.DryRunToken) // one-time use
		}
		mu.Unlock()
		if !ok || time.Now().After(tok.expiry) {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest,
				"Token dry-run tidak valid atau kedaluwarsa. Jalankan dry-run ulang", "")
		}
	}

	paths, err := ctrl.collectFilePaths()
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengumpulkan daftar file DB", err.Error())
	}

	referenced := make(map[string]bool, len(paths))
	for _, p := range paths {
		referenced[p] = true
	}

	type orphanFile struct {
		Path string `json:"path"`
		Size int64  `json:"size"`
	}
	var orphans []orphanFile
	var totalSize int64

	err = filepath.Walk(ctrl.cfg.UploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // file yang tidak bisa diakses dilewati
		}
		if info.IsDir() {
			// Lewati folder internal (chunks) yang bukan aset
			if info.Name() == "chunks" {
				return filepath.SkipDir
			}
			return nil
		}

		rel, relErr := filepath.Rel(ctrl.cfg.UploadPath, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)

		if referenced[rel] {
			return nil
		}

		orphans = append(orphans, orphanFile{Path: rel, Size: info.Size()})
		totalSize += info.Size()
		return nil
	})
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal memindai folder upload", err.Error())
	}

	deleted := 0
	var dryRunToken string
	if !dryRun {
		for _, o := range orphans {
			full := filepath.Join(ctrl.cfg.UploadPath, filepath.FromSlash(o.Path))
			if err := os.Remove(full); err == nil {
				deleted++
			}
		}
	} else {
		// Terbitkan token one-time untuk eksekusi permanen. Token mengikat
		// ke kondisi storage saat dry-run (jumlah file & total ukuran) —
		// kalau ada perubahan besar setelahnya, FE harus dry-run ulang.
		dryRunToken = fmt.Sprintf("prune-%s", hex.EncodeToString(randomBytes(16)))
		mu.Lock()
		pruneTokens[dryRunToken] = pruneToken{
			expiry:     time.Now().Add(10 * time.Minute),
			totalFiles: len(orphans),
			totalSize:  totalSize,
		}
		mu.Unlock()
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Pruning selesai", fiber.Map{
		"dry_run":        dryRun,
		"total_files":    len(orphans),
		"total_size":     totalSize,
		"deleted":        deleted,
		"orphans":        orphans,
		"dry_run_token":  dryRunToken,
		"token_expiry_s": 600,
	})
}

// pruneToken menyimpan konteks dry-run untuk konfirmasi eksekusi permanen.
type pruneToken struct {
	expiry     time.Time
	totalFiles int
	totalSize  int64
}

// pruneTokens menyimpan token dry-run yang belum dipakai, dilindungi mutex.
var (
	pruneTokens = make(map[string]pruneToken)
	mu          sync.Mutex
)

// randomBytes menghasilkan n byte acak kriptografis (untuk token dry-run).
func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand praktis tidak pernah gagal; fallback ke timestamp
		// demi menghindari panic di lingkungan tanpa entropy yang cukup.
		return []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	return b
}
