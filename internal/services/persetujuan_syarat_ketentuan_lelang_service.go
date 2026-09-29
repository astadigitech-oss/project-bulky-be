package services

import (
	"context"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type PersetujuanSyaratKetentuanLelangService interface {
	AmbilSemua(ctx context.Context, params *models.PaginationRequest) ([]models.PersetujuanSyaratKetentuanLelangResponse, *models.PaginationMeta, error)
	AmbilBerdasarkanID(ctx context.Context, id string) (*models.PersetujuanSyaratKetentuanLelangResponse, error)
}

type persetujuanSyaratKetentuanLelangService struct {
	repo repositories.PersetujuanSyaratKetentuanLelangRepository
}

func NewPersetujuanSyaratKetentuanLelangService(
	repo repositories.PersetujuanSyaratKetentuanLelangRepository,
) PersetujuanSyaratKetentuanLelangService {
	return &persetujuanSyaratKetentuanLelangService{repo: repo}
}

func (s *persetujuanSyaratKetentuanLelangService) AmbilSemua(
	ctx context.Context,
	params *models.PaginationRequest,
) ([]models.PersetujuanSyaratKetentuanLelangResponse, *models.PaginationMeta, error) {
	params.SetDefaults()
	consents, total, err := s.repo.AmbilSemua(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	items := make([]models.PersetujuanSyaratKetentuanLelangResponse, 0, len(consents))
	for i := range consents {
		items = append(items, ubahKePersetujuanSyaratKetentuanLelangResponse(&consents[i]))
	}
	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return items, &meta, nil
}

func (s *persetujuanSyaratKetentuanLelangService) AmbilBerdasarkanID(
	ctx context.Context,
	id string,
) (*models.PersetujuanSyaratKetentuanLelangResponse, error) {
	consent, err := s.repo.AmbilBerdasarkanID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := ubahKePersetujuanSyaratKetentuanLelangResponse(consent)
	return &response, nil
}

func ubahKePersetujuanSyaratKetentuanLelangResponse(
	consent *models.PersetujuanSyaratKetentuanLelang,
) models.PersetujuanSyaratKetentuanLelangResponse {
	response := models.PersetujuanSyaratKetentuanLelangResponse{
		ID:              consent.ID.String(),
		BuyerID:         consent.BuyerID.String(),
		BidID:           consent.BidID.String(),
		TermsDocumentID: consent.TermsDocumentID.String(),
		Locale:          consent.Locale,
		ContentHash:     consent.ContentHash,
		DisetujuiAt:     consent.DisetujuiAt,
		IPAddress:       consent.IPAddress,
		UserAgent:       consent.UserAgent,
		CreatedAt:       consent.CreatedAt,
	}
	if consent.Buyer != nil {
		response.BuyerNama = consent.Buyer.Nama
		if consent.Buyer.Email != nil {
			response.BuyerEmail = *consent.Buyer.Email
		}
	}
	if consent.Bid != nil {
		response.BidSequence = consent.Bid.Sequence
		response.BidAmount = consent.Bid.Amount.StringFixed(0)
		if consent.Bid.Batch != nil {
			response.BatchID = consent.Bid.Batch.ID.String()
			response.BatchKode = consent.Bid.Batch.Code
			response.BatchNama = consent.Bid.Batch.NamaID
			if consent.Locale == "en" && consent.Bid.Batch.NamaEN != nil && *consent.Bid.Batch.NamaEN != "" {
				response.BatchNama = *consent.Bid.Batch.NamaEN
			}
		}
	}
	response.TermsJudul = "Syarat dan Ketentuan Lelang"
	if consent.Locale == "en" {
		response.TermsJudul = "Auction Terms and Conditions"
	}
	return response
}
