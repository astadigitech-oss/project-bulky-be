package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

type EmailService interface {
	SendFormulirNotification(to, subject string, data map[string]interface{}) error
}

type emailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromName     string
	fromEmail    string
}

func NewEmailService() EmailService {
	return &emailService{
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnv("SMTP_PORT", "587"),
		smtpUsername: getEnv("SMTP_USERNAME", ""),
		smtpPassword: getEnv("SMTP_PASSWORD", ""),
		fromName:     getEnv("SMTP_FROM_NAME", "Bulky Indonesia"),
		fromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@bulky.id"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

const formulirEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f5a623; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { padding: 30px; background: #f9f9f9; border: 1px solid #ddd; border-top: none; border-radius: 0 0 5px 5px; }
        .field { margin-bottom: 20px; }
        .label { font-weight: bold; color: #555; margin-bottom: 5px; }
        .value { color: #333; padding: 10px; background: white; border-radius: 3px; }
        .footer { text-align: center; margin-top: 20px; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2 style="margin: 0;">📦 Formulir Pemesanan Partai Besar</h2>
        </div>
        <div class="content">
            <p>Ada pengajuan baru untuk pemesanan partai besar:</p>
            
            <div class="field">
                <div class="label">👤 Nama:</div>
                <div class="value">{{.Nama}}</div>
            </div>
            
            <div class="field">
                <div class="label">📞 Telepon:</div>
                <div class="value">{{.Telepon}}</div>
            </div>
            
            <div class="field">
                <div class="label">📍 Alamat:</div>
                <div class="value">{{.Alamat}}</div>
            </div>
            
            <div class="field">
                <div class="label">💰 Anggaran:</div>
                <div class="value">{{.Anggaran}}</div>
            </div>
            
            <div class="field">
                <div class="label">📦 Kategori Produk:</div>
                <div class="value">{{.KategoriStr}}</div>
            </div>
            
            <div class="field">
                <div class="label">🕐 Waktu Submit:</div>
                <div class="value">{{.CreatedAt}}</div>
            </div>
        </div>
        <div class="footer">
            <p>Email ini dikirim otomatis oleh sistem Bulky Indonesia</p>
        </div>
    </div>
</body>
</html>
`

func (s *emailService) SendFormulirNotification(to, subject string, data map[string]interface{}) error {
	// Skip if SMTP not configured
	if s.smtpUsername == "" || s.smtpPassword == "" {
		fmt.Println("SMTP not configured, skipping email send")
		return nil
	}

	// Convert kategori array to string
	if kategori, ok := data["Kategori"].([]string); ok {
		data["KategoriStr"] = strings.Join(kategori, ", ")
	}

	// Parse template
	tmpl, err := template.New("email").Parse(formulirEmailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Prepare email
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", from, to, subject, body.String()))

	// Send email
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	if err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
