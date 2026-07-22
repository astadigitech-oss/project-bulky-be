package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"project-bulky-be/pkg/utils"
)

var nonDigit = regexp.MustCompile(`[^0-9]`)

// normalizePhone menormalkan nomor v1 ke E.164 Indonesia (+62...).
// Mengembalikan ok=false bila nomor kosong/tidak masuk akal.
func normalizePhone(raw string) (string, bool) {
	d := nonDigit.ReplaceAllString(raw, "")
	switch {
	case d == "":
		return "", false
	case strings.HasPrefix(d, "62"):
		// sudah kode negara
	case strings.HasPrefix(d, "0"):
		d = "62" + d[1:]
	case strings.HasPrefix(d, "8"):
		d = "62" + d
	}
	// panjang wajar nomor Indonesia: +62 diikuti 8-13 digit
	if len(d) < 10 || len(d) > 16 {
		return "", false
	}
	p := "+" + d
	if len(p) > 20 { // batas kolom telepon varchar(20)
		return "", false
	}
	return p, true
}

// convertBcrypt mengubah prefix hash Laravel $2y$ menjadi $2a$ agar diterima bcrypt Go.
func convertBcrypt(hash string) string {
	if strings.HasPrefix(hash, "$2y$") {
		return "$2a$" + hash[len("$2y$"):]
	}
	return hash
}

// parseTrans membaca kolom JSON {"id": "...", "en": "..."} v1.
func parseTrans(raw sql.NullString) (idVal, enVal string) {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return "", ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw.String), &m); err != nil {
		return "", ""
	}
	if v, ok := m["id"].(string); ok {
		idVal = strings.TrimSpace(v)
	}
	if v, ok := m["en"].(string); ok {
		enVal = strings.TrimSpace(v)
	}
	return idVal, enVal
}

// parseImages membaca kolom products.images (JSON array of filename).
func parseImages(raw sql.NullString) ([]string, error) {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return nil, nil
	}
	var files []string
	if err := json.Unmarshal([]byte(raw.String), &files); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(files))
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		// simpan basename saja; folder tujuan ditentukan konvensi v2
		out = append(out, path.Base(strings.ReplaceAll(f, "\\", "/")))
	}
	return out, nil
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

func id8(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

// slugSpace melacak slug terpakai dalam satu ruang unique (existing target + alokasi baru).
type slugSpace struct{ used map[string]bool }

func newSlugSpace(existing map[string]bool) *slugSpace {
	u := map[string]bool{}
	for k := range existing {
		u[k] = true
	}
	return &slugSpace{used: u}
}

// alloc mengambil slug unik. Baris terhapus yang bentrok memakai pola
// "<slug>-deleted-<id8>" (mengikuti migrasi v2 000157); baris aktif memakai suffix -1, -2, ...
func (s *slugSpace) alloc(base string, deleted bool, v1ID string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "item-" + id8(v1ID)
	}
	if !s.used[base] {
		s.used[base] = true
		return base
	}
	if deleted {
		cand := base + "-deleted-" + id8(v1ID)
		for i := 2; s.used[cand]; i++ {
			cand = fmt.Sprintf("%s-deleted-%s-%d", base, id8(v1ID), i)
		}
		s.used[cand] = true
		return cand
	}
	for i := 1; ; i++ {
		cand := fmt.Sprintf("%s-%d", base, i)
		if !s.used[cand] {
			s.used[cand] = true
			return cand
		}
	}
}

func slugify(text string) string { return utils.GenerateSlug(text) }

// helper nullable
func strPtr(s string) *string { return &s }

func nullStrPtr(ns sql.NullString) *string {
	if !ns.Valid || strings.TrimSpace(ns.String) == "" {
		return nil
	}
	v := ns.String
	return &v
}

func nullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	v := nt.Time
	return &v
}

func floatOrZero(nf sql.NullFloat64) float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

func boolFromTiny(ni sql.NullInt64) bool {
	return ni.Valid && ni.Int64 != 0
}
