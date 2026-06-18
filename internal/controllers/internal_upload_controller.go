package controllers

import (
	"net/http"
	"project-bulky-be/internal/config"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type InternalUploadController struct {
	cfg *config.Config
}

func NewInternalUploadController(cfg *config.Config) *InternalUploadController {
	return &InternalUploadController{cfg: cfg}
}

// UploadUlasanGambar menerima file gambar dari storefront BE dan menyimpannya ke folder uploads.
// Endpoint ini hanya bisa diakses oleh internal service via X-Internal-Key header.
func (ctrl *InternalUploadController) UploadUlasanGambar(c *fiber.Ctx) error {
	file, err := c.FormFile("gambar")
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "File gambar tidak ditemukan", err.Error())
	}

	if !utils.IsValidImageType(file) {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Tipe file tidak didukung. Hanya jpg, png, webp yang diperbolehkan", "")
	}

	// Max 5MB
	if file.Size > 5*1024*1024 {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Ukuran file maksimal 5MB", "")
	}

	relativePath, err := utils.SaveUploadedFile(file, "ulasan", ctrl.cfg)
	if err != nil {
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan file", err.Error())
	}

	fileURL := utils.GetFileURL(relativePath, ctrl.cfg)

	return utils.SimpleSuccessResponse(c, http.StatusOK, "File berhasil diupload", fiber.Map{
		"path": relativePath,
		"url":  fileURL,
	})
}
