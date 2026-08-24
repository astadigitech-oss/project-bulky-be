package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BackupController struct {
	service services.BackupService
}

func NewBackupController(service services.BackupService) *BackupController {
	return &BackupController{service: service}
}

// List menampilkan daftar file backup database beserta statistik storage
func (ctrl *BackupController) List(c *fiber.Ctx) error {
	result, err := ctrl.service.ListBackups(c.Context())
	if err != nil {
		log.Printf("[backup-controller] Gagal mengambil daftar backup: %v", err)
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar backup database", err.Error())
	}
	return utils.SimpleSuccessResponse(c, http.StatusOK, "Daftar backup database berhasil diambil", result)
}

// Create memicu proses backup database secara manual
func (ctrl *BackupController) Create(c *fiber.Ctx) error {
	adminID := getAdminUUIDFromContext(c)
	adminEmail, _ := c.Locals("user_email").(string)

	item, err := ctrl.service.CreateBackup(c.Context(), models.BackupTriggerManual, adminID, adminEmail)
	if err != nil {
		log.Printf("[backup-controller] Gagal membuat backup database: %v", err)
		if errors.Is(err, services.ErrBackupInProgress) {
			return utils.SimpleErrorResponse(c, http.StatusConflict, err.Error(), "")
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal membuat backup database", err.Error())
	}

	return utils.SimpleSuccessResponse(c, http.StatusCreated, "Backup database berhasil dibuat", item)
}

// Download mengunduh file backup database .sql.gz
func (ctrl *BackupController) Download(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Parameter nama file diperlukan", "")
	}

	filePath, err := ctrl.service.GetBackupFilePath(filename)
	if err != nil {
		if errors.Is(err, services.ErrInvalidFilename) {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, err.Error(), "")
		}
		if errors.Is(err, services.ErrBackupNotFound) {
			return utils.SimpleErrorResponse(c, http.StatusNotFound, err.Error(), "")
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal memproses file backup", err.Error())
	}

	c.Set("Content-Type", "application/gzip")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return c.Download(filePath, filename)
}

// Delete menghapus file backup database tertentu
func (ctrl *BackupController) Delete(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return utils.SimpleErrorResponse(c, http.StatusBadRequest, "Parameter nama file diperlukan", "")
	}

	adminID := getAdminUUIDFromContext(c)
	adminEmail, _ := c.Locals("user_email").(string)

	if err := ctrl.service.DeleteBackup(c.Context(), filename, adminID, adminEmail); err != nil {
		if errors.Is(err, services.ErrInvalidFilename) {
			return utils.SimpleErrorResponse(c, http.StatusBadRequest, err.Error(), "")
		}
		if errors.Is(err, services.ErrBackupNotFound) {
			return utils.SimpleErrorResponse(c, http.StatusNotFound, err.Error(), "")
		}
		return utils.SimpleErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus file backup", err.Error())
	}

	return utils.SimpleSuccessResponse(c, http.StatusOK, "File backup database berhasil dihapus", nil)
}

func getAdminUUIDFromContext(c *fiber.Ctx) *uuid.UUID {
	if val := c.Locals("user_id"); val != nil {
		switch v := val.(type) {
		case uuid.UUID:
			return &v
		case *uuid.UUID:
			return v
		case string:
			if parsed, err := uuid.Parse(v); err == nil {
				return &parsed
			}
		}
	}
	if val := c.Locals("admin_id"); val != nil {
		switch v := val.(type) {
		case uuid.UUID:
			return &v
		case *uuid.UUID:
			return v
		case string:
			if parsed, err := uuid.Parse(v); err == nil {
				return &parsed
			}
		}
	}
	return nil
}
