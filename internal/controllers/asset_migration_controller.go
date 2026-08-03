package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	// Write zip to a temp file to allow streaming large files safely
	tmpFile, err := os.CreateTemp("", "assets-export-*.zip")
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

	// Save uploaded zip to a temp file (avoids loading entire zip into memory)
	tmpZip, err := os.CreateTemp("", "assets-import-*.zip")
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
var v1AssetPrefixes = []struct{ src, dst string }{
	{"storage/app/public/products/", "product-images/"},
	{"storage/app/public/reviews/", "reviews/"},
	{"storage/app/public/public/profile/", "profile/"},
	{"storage/app/public/", ""}, // dokumen/PDF v1 diletakkan langsung di storage/app/public/
}

// ImportAssetsV1 mengimpor ZIP hasil export folder "storage" Laravel v1
// (storage/app/public/...). Karena path yang tersimpan di DB v2 sudah ditransformasi
// migrasi (produk -> "product-images/", dokumen -> "product-documents/",
// ulasan -> "reviews/", foto profil -> "profile/"), file dari ZIP v1 perlu
// dipetakan ulang ke lokasi yang sesuai, bukan diekstrak apa adanya.
//
// Mapping yang didukung (cukup 4 tipe aset publik v1):
//
//	storage/app/public/products/<f>        -> <UPLOAD_PATH>/product-images/<f>
//	storage/app/public/reviews/<f>         -> <UPLOAD_PATH>/reviews/<f>
//	storage/app/public/public/profile/<f>  -> <UPLOAD_PATH>/profile/<f>
//	storage/app/public/<f>.pdf             -> <UPLOAD_PATH>/product-documents/<f>
//
// File yang tidak cocok mapping mana pun (folder lain dari aplikasi v1) di-skip.
func (ctrl *AssetMigrationController) ImportAssetsV1(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File zip tidak ditemukan. Kirim dengan field name 'file'", err.Error())
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File harus berformat .zip", "")
	}

	// Simpan ZIP ke file sementara (hindari muat seluruh zip ke memori)
	tmpZip, err := os.CreateTemp("", "assets-import-v1-*.zip")
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
	var unmatched []string
	type fileResult struct {
		Source string `json:"source"`
		Dest   string `json:"dest,omitempty"`
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	}
	var files []fileResult

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		cleanName := filepath.ToSlash(filepath.Clean(f.Name))
		// Anti path traversal
		if strings.HasPrefix(cleanName, "../") || strings.Contains(cleanName, "/../") {
			skipped++
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "path traversal ditolak"})
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
			skipped++
			if len(unmatched) < 20 {
				unmatched = append(unmatched, cleanName)
			}
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "path tidak dikenali"})
			continue
		}

		destPath := filepath.Join(ctrl.cfg.UploadPath, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			skipped++
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "gagal buat folder"})
			continue
		}

		rc, err := f.Open()
		if err != nil {
			skipped++
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "gagal baca zip"})
			continue
		}

		dst, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			skipped++
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "gagal buat file"})
			continue
		}

		if _, err := io.Copy(dst, rc); err != nil {
			dst.Close()
			rc.Close()
			os.Remove(destPath)
			skipped++
			files = append(files, fileResult{Source: cleanName, Status: "skipped", Reason: "gagal tulis file"})
			continue
		}

		dst.Close()
		rc.Close()
		imported++
		files = append(files, fileResult{Source: cleanName, Dest: relPath, Status: "imported"})
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "Import v1 selesai", fiber.Map{
		"imported":  imported,
		"skipped":   skipped,
		"unmatched": unmatched,
		"files":     files,
	})
}
