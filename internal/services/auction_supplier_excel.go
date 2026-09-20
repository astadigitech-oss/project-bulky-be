package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
)

const (
	maxSupplierWorkbookSize = 10 << 20
	maxSupplierWorkbookRows = 2000
)

type supplierWorkbook struct {
	sheet     string
	rows      [][]string
	headerRow int
}

func (s *auctionService) PreviewSupplierExcel(ctx context.Context, file *multipart.FileHeader) (*dto.AuctionSupplierExcelPreview, error) {
	workbook, err := readSupplierWorkbook(file)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	header := workbook.rows[workbook.headerRow]
	columns := make([]dto.AuctionSupplierExcelColumn, len(header))
	for index, value := range header {
		label := strings.TrimSpace(value)
		if label == "" {
			label = "Kolom " + excelColumnLetter(index)
		}
		column := dto.AuctionSupplierExcelColumn{Index: index, Letter: excelColumnLetter(index), Header: label, Samples: []string{}}
		for rowIndex := workbook.headerRow + 1; rowIndex < len(workbook.rows) && len(column.Samples) < 3; rowIndex++ {
			if index < len(workbook.rows[rowIndex]) && strings.TrimSpace(workbook.rows[rowIndex][index]) != "" {
				column.Samples = append(column.Samples, strings.TrimSpace(workbook.rows[rowIndex][index]))
			}
		}
		columns[index] = column
	}

	return &dto.AuctionSupplierExcelPreview{
		SheetName: workbook.sheet,
		HeaderRow: workbook.headerRow,
		Columns:   columns,
	}, nil
}

func (s *auctionService) ImportSupplierExcel(ctx context.Context, file *multipart.FileHeader, mapping dto.AuctionSupplierExcelMapping, title string, adminID uuid.UUID) (*dto.AuctionSupplierExcelImport, error) {
	workbook, err := readSupplierWorkbook(file)
	if err != nil {
		return nil, err
	}
	if mapping.NameColumn < 0 || mapping.PriceColumn < 0 || mapping.QuantityColumn < 0 ||
		mapping.NameColumn == mapping.PriceColumn || mapping.NameColumn == mapping.QuantityColumn || mapping.PriceColumn == mapping.QuantityColumn {
		return nil, auctionErr(400, "Pilih kolom nama, harga, dan qty yang berbeda")
	}
	if mapping.HeaderRow != workbook.headerRow {
		return nil, auctionErr(400, "Header Excel berubah. Unggah ulang file untuk membaca kolom terbaru")
	}
	maxColumn := max(mapping.NameColumn, mapping.PriceColumn, mapping.QuantityColumn)
	if maxColumn >= len(workbook.rows[workbook.headerRow]) {
		return nil, auctionErr(400, "Kolom yang dipilih tidak ditemukan di file Excel")
	}

	items := make([]dto.AuctionDraftItem, 0)
	for rowIndex := workbook.headerRow + 1; rowIndex < len(workbook.rows); rowIndex++ {
		row := workbook.rows[rowIndex]
		name := strings.TrimSpace(excelCell(row, mapping.NameColumn))
		priceText := strings.TrimSpace(excelCell(row, mapping.PriceColumn))
		quantityText := strings.TrimSpace(excelCell(row, mapping.QuantityColumn))
		if name == "" && priceText == "" && quantityText == "" {
			continue
		}
		excelRow := rowIndex + 1
		if name == "" || len([]rune(name)) > 255 {
			return nil, auctionErr(400, fmt.Sprintf("Baris %d: nama item wajib diisi dan maksimal 255 karakter", excelRow))
		}
		price, ok := parseExcelWholeNumber(priceText)
		if !ok || !price.IsPositive() {
			return nil, auctionErr(400, fmt.Sprintf("Baris %d: harga harus berupa Rupiah bulat lebih dari nol", excelRow))
		}
		quantityDecimal, ok := parseExcelWholeNumber(quantityText)
		if !ok || !quantityDecimal.IsPositive() || !quantityDecimal.IsInteger() {
			return nil, auctionErr(400, fmt.Sprintf("Baris %d: qty harus berupa bilangan bulat lebih dari nol", excelRow))
		}
		quantity := quantityDecimal.IntPart()
		if int64(int(quantity)) != quantity {
			return nil, auctionErr(400, fmt.Sprintf("Baris %d: qty terlalu besar", excelRow))
		}
		items = append(items, dto.AuctionDraftItem{
			SourceType: "MANUAL",
			Nama:       name,
			UnitPrice:  price.StringFixed(0),
			Quantity:   int(quantity),
		})
		if len(items) > maxSupplierWorkbookRows {
			return nil, auctionErr(400, fmt.Sprintf("Excel maksimal berisi %d baris item", maxSupplierWorkbookRows))
		}
	}
	if len(items) == 0 {
		return nil, auctionErr(400, "Tidak ada item yang dapat diimpor dari Excel")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(title) == "" {
		title = "Supplier Auction Item List"
	}
	pdfBytes, err := buildSupplierItemsPDF(title, items)
	if err != nil {
		return nil, auctionErr(500, "PDF daftar item gagal dibuat")
	}
	asset, err := s.storeGeneratedAuctionPDF(ctx, pdfBytes, title, adminID)
	if err != nil {
		return nil, err
	}
	return &dto.AuctionSupplierExcelImport{Items: items, PDF: mapAsset(asset, s.cfg)}, nil
}

func readSupplierWorkbook(file *multipart.FileHeader) (*supplierWorkbook, error) {
	if file == nil {
		return nil, auctionErr(400, "File Excel wajib dipilih")
	}
	if file.Size <= 0 || file.Size > maxSupplierWorkbookSize {
		return nil, auctionErr(413, "Ukuran Excel maksimal 10MB")
	}
	if strings.ToLower(filepath.Ext(file.Filename)) != ".xlsx" {
		return nil, auctionErr(415, "Format Excel tidak didukung. Gunakan file .xlsx")
	}
	reader, err := file.Open()
	if err != nil {
		return nil, auctionErr(400, "File Excel tidak dapat dibaca")
	}
	defer reader.Close()
	content, err := io.ReadAll(io.LimitReader(reader, maxSupplierWorkbookSize+1))
	if err != nil || len(content) == 0 || len(content) > maxSupplierWorkbookSize {
		return nil, auctionErr(400, "File Excel tidak valid atau melebihi 10MB")
	}
	book, err := excelize.OpenReader(bytes.NewReader(content), excelize.Options{
		UnzipSizeLimit:    64 << 20,
		UnzipXMLSizeLimit: 16 << 20,
	})
	if err != nil {
		return nil, auctionErr(400, "File tidak dapat dibuka sebagai workbook Excel .xlsx")
	}
	defer book.Close()
	sheets := book.GetSheetList()
	if len(sheets) == 0 {
		return nil, auctionErr(400, "Workbook Excel tidak memiliki sheet")
	}
	rows, err := book.GetRows(sheets[0])
	if err != nil || len(rows) == 0 {
		return nil, auctionErr(400, "Sheet pertama Excel kosong atau tidak dapat dibaca")
	}
	headerRow := -1
	for index, row := range rows {
		for _, value := range row {
			if strings.TrimSpace(value) != "" {
				headerRow = index
				break
			}
		}
		if headerRow >= 0 {
			break
		}
	}
	if headerRow < 0 {
		return nil, auctionErr(400, "Sheet pertama Excel tidak memiliki baris header")
	}
	return &supplierWorkbook{sheet: sheets[0], rows: rows, headerRow: headerRow}, nil
}

func excelCell(row []string, column int) string {
	if column < 0 || column >= len(row) {
		return ""
	}
	return row[column]
}

func excelColumnLetter(index int) string {
	letter := ""
	for index >= 0 {
		letter = string(rune('A'+index%26)) + letter
		index = index/26 - 1
	}
	return letter
}

func parseExcelWholeNumber(value string) (decimal.Decimal, bool) {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(strings.ToLower(value), "rp", "")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "\u00a0", "")
	if value == "" {
		return decimal.Zero, false
	}
	lastDot, lastComma := strings.LastIndex(value, "."), strings.LastIndex(value, ",")
	decimalSeparator := byte(0)
	if lastDot >= 0 && lastComma >= 0 {
		if lastDot > lastComma {
			decimalSeparator = '.'
		} else {
			decimalSeparator = ','
		}
	} else if lastDot >= 0 || lastComma >= 0 {
		separator := byte('.')
		last := lastDot
		if lastComma >= 0 {
			separator, last = ',', lastComma
		}
		tailDigits := len(value) - last - 1
		separatorCount := strings.Count(value, string(separator))
		if separatorCount == 1 && tailDigits != 3 {
			decimalSeparator = separator
		}
	}
	var normalized strings.Builder
	for index := 0; index < len(value); index++ {
		character := value[index]
		if character >= '0' && character <= '9' || character == '-' {
			normalized.WriteByte(character)
			continue
		}
		if (character == '.' || character == ',') && decimalSeparator == character {
			normalized.WriteByte('.')
		}
	}
	number, err := decimal.NewFromString(normalized.String())
	if err != nil || !number.IsInteger() {
		return decimal.Zero, false
	}
	return number, true
}

func buildSupplierItemsPDF(title string, items []dto.AuctionDraftItem) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.SetAutoPageBreak(true, 14)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.MultiCell(0, 8, pdf.UnicodeTranslatorFromDescriptor("cp1252")(title), "", "C", false)
	pdf.Ln(3)

	widths := []float64{145, 55, 22, 43}
	drawHeader := func() {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(235, 238, 243)
		for index, heading := range []string{"Product/Item Name", "Unit Price (IDR)", "Quantity", "Subtotal (IDR)"} {
			align := "L"
			if index > 0 {
				align = "R"
			}
			pdf.CellFormat(widths[index], 9, heading, "1", 0, align, true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 9)
	}
	drawHeader()

	for _, item := range items {
		name := pdf.UnicodeTranslatorFromDescriptor("cp1252")(item.Nama)
		lines := pdf.SplitText(name, widths[0]-4)
		rowHeight := float64(len(lines)) * 5
		if rowHeight < 9 {
			rowHeight = 9
		}
		if pdf.GetY()+rowHeight > 196 {
			pdf.AddPage()
			drawHeader()
		}
		startX, startY := pdf.GetX(), pdf.GetY()
		pdf.MultiCell(widths[0], 5, name, "", "L", false)
		pdf.Rect(startX, startY, widths[0], rowHeight, "D")
		pdf.SetXY(startX+widths[0], startY)
		pdf.CellFormat(widths[1], rowHeight, formatExcelRupiah(item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[2], rowHeight, strconv.Itoa(item.Quantity), "1", 0, "R", false, 0, "")
		price, _ := decimal.NewFromString(item.UnitPrice)
		subtotal := price.Mul(decimal.NewFromInt(int64(item.Quantity))).StringFixed(0)
		pdf.CellFormat(widths[3], rowHeight, formatExcelRupiah(subtotal), "1", 0, "R", false, 0, "")
		pdf.SetXY(startX, startY+rowHeight)
	}

	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func formatExcelRupiah(value string) string {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return value
	}
	text := strconv.FormatInt(number, 10)
	if number < 0 {
		text = text[1:]
	}
	for position := len(text) - 3; position > 0; position -= 3 {
		text = text[:position] + "." + text[position:]
	}
	if number < 0 {
		text = "-" + text
	}
	return text
}

func (s *auctionService) storeGeneratedAuctionPDF(ctx context.Context, content []byte, title string, adminID uuid.UUID) (*models.AuctionAsset, error) {
	assetID := uuid.New()
	filename := "auction-items-" + assetID.String() + ".pdf"
	storageKey := filepath.ToSlash(filepath.Join("auction", filename))
	path := filepath.Join(s.cfg.UploadPath, filepath.FromSlash(storageKey))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, auctionErr(500, "Folder penyimpanan PDF tidak dapat dibuat")
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return nil, auctionErr(500, "PDF daftar item tidak dapat disimpan")
	}
	asset := &models.AuctionAsset{
		ID:           assetID,
		UploadedBy:   adminID,
		Kind:         "PDF",
		StorageKey:   storageKey,
		OriginalName: "daftar-item-supplier.pdf",
		MimeType:     "application/pdf",
		SizeBytes:    int64(len(content)),
	}
	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	return asset, nil
}
