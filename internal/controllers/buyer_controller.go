package controllers

import (
	"fmt"
	"net/http"
	"time"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/services"
	"project-bulky-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type BuyerController struct {
	service     services.BuyerService
	activityLog services.ActivityLogService
}

func NewBuyerController(service services.BuyerService, activityLog services.ActivityLogService) *BuyerController {
	return &BuyerController{service: service, activityLog: activityLog}
}

func (c *BuyerController) FindAll(ctx *fiber.Ctx) error {
	var params models.BuyerFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, meta, err := c.service.FindAll(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.PaginatedSuccessResponse(ctx, "Data buyer berhasil diambil", items, *meta)
}

// Export returns the complete buyer dataset using the same search and sort
// filters as the list endpoint.
func (c *BuyerController) Export(ctx *fiber.Ctx) error {
	var params models.BuyerFilterRequest
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	items, err := c.service.FindAllForExport(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengekspor data buyer", nil)
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Buyer"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"No", "Nama", "Username", "Email", "Telepon", "Aktif", "Terverifikasi", "Login Terakhir", "Terdaftar"}
	for column, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(column+1, 1)
		_ = f.SetCellValue(sheet, cell, header)
	}
	for index, buyer := range items {
		row := index + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), index+1)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), buyer.Nama)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), nullableString(buyer.Username))
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), nullableString(buyer.Email))
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), buyer.Telepon)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), buyer.IsActive)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), buyer.IsVerified)
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), formatExportTime(buyer.LastLoginAt))
		_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", row), buyer.CreatedAt.Format(time.RFC3339))
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"buyer-bulky-%s.xlsx\"", time.Now().Format("20060102")))
	if err := f.Write(ctx.Response().BodyWriter()); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat file Excel", nil)
	}
	return nil
}

func nullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func formatExportTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

func (c *BuyerController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result, err := c.service.FindByID(ctx.UserContext(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Detail buyer berhasil diambil", result)
}

func (c *BuyerController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.service.Delete(ctx.UserContext(), id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "buyer tidak ditemukan" {
			status = http.StatusNotFound
		}
		return utils.ErrorResponse(ctx, status, err.Error(), nil)
	}
	c.activityLog.Log(ctx, models.ActionDelete, "buyer", "Buyer berhasil dihapus")
	return utils.SuccessResponse(ctx, "Buyer berhasil dihapus", nil)
}

func (c *BuyerController) ResetPassword(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.ResetBuyerPasswordRequest
	if err := BindJSON(ctx, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validasi gagal", parseValidationErrors(err))
	}

	if err := c.service.ResetPassword(ctx.UserContext(), id, &req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	c.activityLog.Log(ctx, models.ActionUpdate, "buyer", "Password buyer berhasil direset")
	return utils.SuccessResponse(ctx, "Password buyer berhasil direset", nil)
}

func (c *BuyerController) GetStatistik(ctx *fiber.Ctx) error {
	result, err := c.service.GetStatistik(ctx.UserContext())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Statistik buyer berhasil diambil", result)
}

func (c *BuyerController) GetChart(ctx *fiber.Ctx) error {
	var params models.ChartParams
	if err := ctx.QueryParser(&params); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Parameter tidak valid", nil)
	}

	result, err := c.service.GetChart(ctx.UserContext(), &params)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SuccessResponse(ctx, "Data chart berhasil diambil", result)
}
