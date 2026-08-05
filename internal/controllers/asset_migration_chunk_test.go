package controllers

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"project-bulky-be/internal/config"

	"github.com/gofiber/fiber/v2"
)

// buildChunkRequest membuat request multipart ke /import-v1/chunk
// dengan chunk_data dari data byte yang diberikan.
func buildChunkRequest(uploadID, index, total string, data []byte, sha string) *http.Request {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("upload_id", uploadID)
	w.WriteField("chunk_index", index)
	w.WriteField("total_chunks", total)
	if sha != "" {
		w.WriteField("chunk_sha1", sha)
	}
	fw, _ := w.CreateFormFile("chunk_data", "chunk.bin")
	fw.Write(data)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/import-v1/chunk", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func sha1Of(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

// setupChunkTestApp membuat Fiber app dengan UploadV1Chunk terdaftar
// (tanpa auth — untuk menguji logika resume/checksum saja).
func setupChunkTestApp(t *testing.T) (*fiber.App, *AssetMigrationController) {
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MiB — cukup untuk chunk 8 MiB uji
	})
	dir := t.TempDir()
	ctrl := &AssetMigrationController{
		cfg: &config.Config{UploadPath: filepath.Join(dir, "uploads")},
	}
	app.Post("/import-v1/chunk", ctrl.UploadV1Chunk)
	app.Delete("/import-v1/chunk", ctrl.AbortV1Chunk)
	return app, ctrl
}

func TestUploadV1ChunkResume(t *testing.T) {
	app, ctrl := setupChunkTestApp(t)

	chunkData := bytes.Repeat([]byte("A"), 8*1024*1024) // 8 MiB
	sum := sha1Of(chunkData)

	// 1) Upload pertama — chunk tersimpan, sudah_exist=false
	req := buildChunkRequest("resume-test", "0", "3", chunkData, sum)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request gagal: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload pertama: expected 200, got %d", resp.StatusCode)
	}
	body := readBody(t, resp)
	if !bytes.Contains(body, []byte(`"already_exist":false`)) {
		t.Fatalf("upload pertama: already_exist harus false, body: %s", body)
	}

	// 2) Kirim ulang dengan sha1 sama — server harus SKIP upload (already_exist=true)
	req = buildChunkRequest("resume-test", "0", "3", chunkData, sum)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("request resume gagal: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("resume: expected 200, got %d", resp.StatusCode)
	}
	body = readBody(t, resp)
	if !bytes.Contains(body, []byte(`"already_exist":true`)) {
		t.Fatalf("resume: already_exist harus true, body: %s", body)
	}

	// 3) Checksum beda (data korup/berbeda) → 422 + chunk dihapus
	corrupt := bytes.Repeat([]byte("B"), 8*1024*1024)
	req = buildChunkRequest("resume-test", "1", "3", corrupt, sha1Of(corrupt)+"deadbeef")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("request korup gagal: %v", err)
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("checksum beda: expected 422, got %d", resp.StatusCode)
	}
	// chunk_00001 tidak boleh tersimpan
	chunk1Path := filepath.Join(ctrl.chunksDirV1("resume-test"), "chunk_00001")
	if _, err := os.Stat(chunk1Path); !os.IsNotExist(err) {
		t.Fatalf("chunk korup seharusnya dihapus, tapi masih ada: %s", chunk1Path)
	}
}

// TestAbortV1ChunkSingle memastikan hapus per-chunk (reset) berjalan.
func TestAbortV1ChunkSingle(t *testing.T) {
	app, ctrl := setupChunkTestApp(t)

	data := bytes.Repeat([]byte("C"), 1024)
	sum := sha1Of(data)

	// Upload 2 chunk
	for i := 0; i < 2; i++ {
		req := buildChunkRequest("abort-test", fmt.Sprintf("%d", i), "2", data, sum)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("upload chunk %d gagal: %v status=%d", i, err, resp.StatusCode)
		}
	}

	// Hapus chunk 0 saja
	req := httptest.NewRequest(http.MethodDelete,
		"/import-v1/chunk?upload_id=abort-test&chunk_index=0", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request abort gagal: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("abort per-chunk: expected 200, got %d", resp.StatusCode)
	}

	// chunk_00000 hilang, chunk_00001 masih ada
	dir := ctrl.chunksDirV1("abort-test")
	if _, err := os.Stat(filepath.Join(dir, "chunk_00000")); !os.IsNotExist(err) {
		t.Fatal("chunk_00000 seharusnya terhapus")
	}
	if _, err := os.Stat(filepath.Join(dir, "chunk_00001")); err != nil {
		t.Fatal("chunk_00001 seharusnya masih ada")
	}

	// Hapus seluruh folder (abort semua) → idempotent sukses
	req = httptest.NewRequest(http.MethodDelete, "/import-v1/chunk?upload_id=abort-test", nil)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("abort semua: expected 200, got %d err=%v", resp.StatusCode, err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("folder chunk seharusnya terhapus total")
	}
}

func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("baca body gagal: %v", err)
	}
	return b
}
