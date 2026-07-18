package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type BuyerDisclaimerConsentService interface {
	// Admin: semua consent (audit log)
	GetAllConsents(ctx context.Context, params *models.PaginationRequest) ([]models.DisclaimerConsentAdminResponse, *models.PaginationMeta, error)

	// Admin: detail consent berdasarkan pesanan_id
	GetConsentByPesanan(ctx context.Context, pesananID string) (*models.DisclaimerConsentAdminResponse, error)
}

type buyerDisclaimerConsentService struct {
	repo repositories.BuyerDisclaimerConsentRepository
}

func NewBuyerDisclaimerConsentService(
	repo repositories.BuyerDisclaimerConsentRepository,
) BuyerDisclaimerConsentService {
	return &buyerDisclaimerConsentService{repo: repo}
}

func (s *buyerDisclaimerConsentService) GetAllConsents(
	ctx context.Context,
	params *models.PaginationRequest,
) ([]models.DisclaimerConsentAdminResponse, *models.PaginationMeta, error) {
	consents, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := make([]models.DisclaimerConsentAdminResponse, 0, len(consents))
	for i := range consents {
		items = append(items, *toConsentAdminResponse(&consents[i]))
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return items, &meta, nil
}

func (s *buyerDisclaimerConsentService) GetConsentByPesanan(
	ctx context.Context,
	pesananID string,
) (*models.DisclaimerConsentAdminResponse, error) {
	consent, err := s.repo.FindByPesananID(ctx, pesananID)
	if err != nil {
		return nil, errors.New("data persetujuan disclaimer tidak ditemukan")
	}
	return toConsentAdminResponse(consent), nil
}

// --- helpers ---

func toConsentAdminResponse(c *models.BuyerDisclaimerConsent) *models.DisclaimerConsentAdminResponse {
	resp := &models.DisclaimerConsentAdminResponse{
		ID:           c.ID.String(),
		BuyerID:      c.BuyerID.String(),
		PesananID:    c.PesananID.String(),
		DisclaimerID: c.DisclaimerID.String(),
		DisetujuiAt:  c.DisetujuiAt,
		IPAddress:    c.IPAddress,
		UserAgent:    c.UserAgent,
		CreatedAt:    c.CreatedAt,
	}

	if c.Buyer != nil {
		resp.BuyerNama = c.Buyer.Nama
		if c.Buyer.Email != nil {
			resp.BuyerEmail = *c.Buyer.Email
		}
	}
	if c.Pesanan != nil {
		resp.PesananKode = c.Pesanan.Kode
	}

	return resp
}
