package services

import (
	"context"
	"errors"
	"fmt"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormulirPartaiBesarService interface {
	// Config
	GetConfig(ctx context.Context) (*models.FormulirConfigResponse, error)
	UpdateConfig(ctx context.Context, req *models.UpdateFormulirConfigRequest) (*models.FormulirConfigResponse, error)

	// Anggaran
	CreateAnggaran(ctx context.Context, req *models.CreateAnggaranRequest) (*models.AnggaranResponse, error)
	FindAllAnggaran(ctx context.Context) ([]models.AnggaranResponse, error)
	FindAnggaranByID(ctx context.Context, id string) (*models.AnggaranResponse, error)
	UpdateAnggaran(ctx context.Context, id string, req *models.UpdateAnggaranRequest) (*models.AnggaranResponse, error)
	DeleteAnggaran(ctx context.Context, id string) error
	ReorderAnggaran(ctx context.Context, req *models.ReorderRequest) error

	// Submission
	SubmitFormulir(ctx context.Context, buyerID *uuid.UUID, req *models.CreateFormulirSubmissionRequest) (string, error)
	FindAllSubmission(ctx context.Context, params *models.FormulirSubmissionFilterRequest) ([]models.FormulirSubmissionListResponse, *models.PaginationMeta, error)
	FindSubmissionByID(ctx context.Context, id string) (*models.FormulirSubmissionDetailResponse, error)
	ResendEmail(ctx context.Context, id string) error

	// Options for buyer
	GetOptions(ctx context.Context) (*models.FormulirOptionsResponse, error)
}

type formulirPartaiBesarService struct {
	repo         repositories.FormulirPartaiBesarRepository
	kategoriRepo repositories.KategoriProdukRepository
	emailService EmailService
}

func NewFormulirPartaiBesarService(
	repo repositories.FormulirPartaiBesarRepository,
	kategoriRepo repositories.KategoriProdukRepository,
	emailService EmailService,
) FormulirPartaiBesarService {
	return &formulirPartaiBesarService{
		repo:         repo,
		kategoriRepo: kategoriRepo,
		emailService: emailService,
	}
}

// ========================================
// Config
// ========================================

func (s *formulirPartaiBesarService) GetConfig(ctx context.Context) (*models.FormulirConfigResponse, error) {
	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil konfigurasi")
	}

	return &models.FormulirConfigResponse{
		ID:          config.ID.String(),
		DaftarEmail: config.GetEmails(),
		UpdatedAt:   config.UpdatedAt,
	}, nil
}

func (s *formulirPartaiBesarService) UpdateConfig(ctx context.Context, req *models.UpdateFormulirConfigRequest) (*models.FormulirConfigResponse, error) {
	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil konfigurasi")
	}

	config.SetEmails(req.DaftarEmail)

	if err := s.repo.UpdateConfig(ctx, config); err != nil {
		return nil, errors.New("gagal mengupdate konfigurasi")
	}

	return &models.FormulirConfigResponse{
		ID:          config.ID.String(),
		DaftarEmail: config.GetEmails(),
		UpdatedAt:   config.UpdatedAt,
	}, nil
}

// ========================================
// Anggaran
// ========================================

func (s *formulirPartaiBesarService) CreateAnggaran(ctx context.Context, req *models.CreateAnggaranRequest) (*models.AnggaranResponse, error) {
	anggaran := &models.FormulirPartaiBesarAnggaran{
		Label:  req.Label,
		Urutan: req.Urutan,
	}

	if err := s.repo.CreateAnggaran(ctx, anggaran); err != nil {
		return nil, errors.New("gagal membuat anggaran")
	}

	return &models.AnggaranResponse{
		ID:     anggaran.ID.String(),
		Label:  anggaran.Label,
		Urutan: anggaran.Urutan,
	}, nil
}

func (s *formulirPartaiBesarService) FindAllAnggaran(ctx context.Context) ([]models.AnggaranResponse, error) {
	items, err := s.repo.FindAllAnggaran(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data anggaran")
	}

	var responses []models.AnggaranResponse
	for _, item := range items {
		responses = append(responses, models.AnggaranResponse{
			ID:     item.ID.String(),
			Label:  item.Label,
			Urutan: item.Urutan,
		})
	}

	return responses, nil
}

func (s *formulirPartaiBesarService) FindAnggaranByID(ctx context.Context, id string) (*models.AnggaranResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	anggaran, err := s.repo.FindAnggaranByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("anggaran tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data anggaran")
	}

	return &models.AnggaranResponse{
		ID:     anggaran.ID.String(),
		Label:  anggaran.Label,
		Urutan: anggaran.Urutan,
	}, nil
}

func (s *formulirPartaiBesarService) UpdateAnggaran(ctx context.Context, id string, req *models.UpdateAnggaranRequest) (*models.AnggaranResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	anggaran, err := s.repo.FindAnggaranByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("anggaran tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data anggaran")
	}

	if req.Label != nil {
		anggaran.Label = *req.Label
	}
	if req.Urutan != nil {
		anggaran.Urutan = *req.Urutan
	}

	if err := s.repo.UpdateAnggaran(ctx, anggaran); err != nil {
		return nil, errors.New("gagal mengupdate anggaran")
	}

	return &models.AnggaranResponse{
		ID:     anggaran.ID.String(),
		Label:  anggaran.Label,
		Urutan: anggaran.Urutan,
	}, nil
}

func (s *formulirPartaiBesarService) DeleteAnggaran(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	// Check if exists
	_, err = s.repo.FindAnggaranByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("anggaran tidak ditemukan")
		}
		return errors.New("gagal mengambil data anggaran")
	}

	if err := s.repo.DeleteAnggaran(ctx, uid); err != nil {
		return errors.New("gagal menghapus anggaran")
	}

	return nil
}

func (s *formulirPartaiBesarService) ReorderAnggaran(ctx context.Context, req *models.ReorderRequest) error {
	if err := s.repo.ReorderAnggaran(ctx, req.Items); err != nil {
		return errors.New("gagal mengubah urutan anggaran")
	}
	return nil
}

// ========================================
// Submission
// ========================================

func (s *formulirPartaiBesarService) SubmitFormulir(ctx context.Context, buyerID *uuid.UUID, req *models.CreateFormulirSubmissionRequest) (string, error) {
	// Validate anggaran exists
	anggaranID, _ := uuid.Parse(req.AnggaranID)
	_, err := s.repo.FindAnggaranByID(ctx, anggaranID)
	if err != nil {
		return "", errors.New("anggaran tidak ditemukan")
	}

	// Create submission
	submission := &models.FormulirPartaiBesarSubmission{
		BuyerID:    buyerID,
		Nama:       req.Nama,
		Telepon:    req.Telepon,
		Alamat:     req.Alamat,
		AnggaranID: &anggaranID,
	}
	submission.SetKategoriIDs(req.KategoriIDs)

	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		return "", errors.New("gagal menyimpan formulir")
	}

	// Get config for email recipients
	config, err := s.repo.GetConfig(ctx)
	if err == nil && len(config.GetEmails()) > 0 {
		// Send email notification
		go func() {
			if err := s.sendEmailNotification(context.Background(), submission); err != nil {
				fmt.Printf("Failed to send email: %v\n", err)
			}
		}()
	}

	return submission.ID.String(), nil
}

func (s *formulirPartaiBesarService) sendEmailNotification(ctx context.Context, submission *models.FormulirPartaiBesarSubmission) error {
	// Get config
	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Get anggaran label
	anggaranLabel := ""
	if submission.AnggaranID != nil {
		anggaran, err := s.repo.FindAnggaranByID(ctx, *submission.AnggaranID)
		if err == nil {
			anggaranLabel = anggaran.Label
		}
	}

	// Get kategori names
	kategoriIDs := submission.GetKategoriIDs()
	var kategoriNames []string
	for _, kid := range kategoriIDs {
		kategori, err := s.kategoriRepo.FindByID(ctx, kid.String())
		if err == nil {
			kategoriNames = append(kategoriNames, kategori.NamaID)
		}
	}

	// Send email to all recipients
	recipients := config.GetEmails()
	subject := "Formulir Pemesanan Partai Besar Baru"

	for _, recipient := range recipients {
		if err := s.emailService.SendFormulirNotification(recipient, subject, map[string]interface{}{
			"Nama":      submission.Nama,
			"Telepon":   submission.Telepon,
			"Alamat":    submission.Alamat,
			"Anggaran":  anggaranLabel,
			"Kategori":  kategoriNames,
			"CreatedAt": submission.CreatedAt.Format("02 January 2006 15:04"),
		}); err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", recipient, err)
		}
	}

	// Update email sent status
	now := time.Now()
	submission.EmailSent = true
	submission.EmailSentAt = &now
	s.repo.UpdateSubmission(ctx, submission)

	return nil
}

func (s *formulirPartaiBesarService) FindAllSubmission(ctx context.Context, params *models.FormulirSubmissionFilterRequest) ([]models.FormulirSubmissionListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	items, total, err := s.repo.FindAllSubmission(ctx, params)
	if err != nil {
		return nil, nil, errors.New("gagal mengambil data submission")
	}

	var responses []models.FormulirSubmissionListResponse
	for _, item := range items {
		anggaranLabel := ""
		if item.Anggaran != nil {
			anggaranLabel = item.Anggaran.Label
		}

		// Get kategori names
		kategoriIDs := item.GetKategoriIDs()
		var kategoriNames []string
		for _, kid := range kategoriIDs {
			kategori, err := s.kategoriRepo.FindByID(ctx, kid.String())
			if err == nil {
				kategoriNames = append(kategoriNames, kategori.NamaID)
			}
		}

		responses = append(responses, models.FormulirSubmissionListResponse{
			ID:        item.ID.String(),
			Nama:      item.Nama,
			Telepon:   item.Telepon,
			Alamat:    item.Alamat,
			Anggaran:  anggaranLabel,
			Kategori:  kategoriNames,
			EmailSent: item.EmailSent,
			CreatedAt: item.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return responses, &meta, nil
}

func (s *formulirPartaiBesarService) FindSubmissionByID(ctx context.Context, id string) (*models.FormulirSubmissionDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	submission, err := s.repo.FindSubmissionByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("submission tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data submission")
	}

	response := &models.FormulirSubmissionDetailResponse{
		ID:          submission.ID.String(),
		Nama:        submission.Nama,
		Telepon:     submission.Telepon,
		Alamat:      submission.Alamat,
		EmailSent:   submission.EmailSent,
		EmailSentAt: submission.EmailSentAt,
		CreatedAt:   submission.CreatedAt,
	}

	// Buyer info
	if submission.Buyer != nil {
		response.Buyer = &models.BuyerListResponse{
			ID:       submission.Buyer.ID.String(),
			Nama:     submission.Buyer.Nama,
			Username: submission.Buyer.Username,
			Email:    submission.Buyer.Email,
			Telepon:  submission.Buyer.Telepon,
		}
	}

	// Anggaran info
	if submission.Anggaran != nil {
		response.Anggaran = &models.AnggaranResponse{
			ID:     submission.Anggaran.ID.String(),
			Label:  submission.Anggaran.Label,
			Urutan: submission.Anggaran.Urutan,
		}
	}

	// Kategori info
	kategoriIDs := submission.GetKategoriIDs()
	for _, kid := range kategoriIDs {
		kategori, err := s.kategoriRepo.FindByID(ctx, kid.String())
		if err == nil {
			response.Kategori = append(response.Kategori, models.DropdownItem{
				ID:   kategori.ID.String(),
				Nama: kategori.NamaID,
				Slug: kategori.Slug,
			})
		}
	}

	return response, nil
}

func (s *formulirPartaiBesarService) ResendEmail(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	submission, err := s.repo.FindSubmissionByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("submission tidak ditemukan")
		}
		return errors.New("gagal mengambil data submission")
	}

	// Send email
	if err := s.sendEmailNotification(ctx, submission); err != nil {
		return errors.New("gagal mengirim email")
	}

	return nil
}

func (s *formulirPartaiBesarService) GetOptions(ctx context.Context) (*models.FormulirOptionsResponse, error) {
	// Get anggaran
	anggaranItems, err := s.repo.FindAllAnggaran(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data anggaran")
	}

	var anggaranList []models.AnggaranResponse
	for _, item := range anggaranItems {
		anggaranList = append(anggaranList, models.AnggaranResponse{
			ID:     item.ID.String(),
			Label:  item.Label,
			Urutan: item.Urutan,
		})
	}

	// Get kategori (active only)
	kategoriItems, err := s.kategoriRepo.GetAllForDropdown(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil data kategori")
	}

	var kategoriList []models.DropdownItem
	for _, item := range kategoriItems {
		kategoriList = append(kategoriList, models.DropdownItem{
			ID:   item.ID.String(),
			Nama: item.NamaID,
			Slug: item.Slug,
		})
	}

	return &models.FormulirOptionsResponse{
		Anggaran: anggaranList,
		Kategori: kategoriList,
	}, nil
}
