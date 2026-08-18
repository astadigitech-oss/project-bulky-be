package models
package models

// WMSConnectionInfo hasil cek koneksi/identitas client dari WMS
// (GET /api/integration/me). Field Data dibiarkan bertipe interface{} karena
// bentuk persisnya belum didokumentasikan detail oleh tim WMS.
type WMSConnectionInfo struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
