package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeasonalCampaignService interface {
	Create(context.Context, *models.CreateSeasonalCampaignRequest, *uuid.UUID) (*models.SeasonalCampaignResponse, error)
	FindByID(context.Context, string) (*models.SeasonalCampaignResponse, error)
	FindAll(context.Context, *models.SeasonalCampaignFilterRequest) ([]models.SeasonalCampaignResponse, *models.PaginationMeta, error)
	Update(context.Context, string, *models.UpdateSeasonalCampaignRequest, *uuid.UUID) (*models.SeasonalCampaignResponse, error)
	Publish(context.Context, string, *uuid.UUID) (*models.SeasonalCampaignResponse, error)
	Cancel(context.Context, string, *uuid.UUID) (*models.SeasonalCampaignResponse, error)
	Delete(context.Context, string) (*models.SeasonalCampaignResponse, error)
	Preview(context.Context, string) (*models.SeasonalCampaignResponse, error)
}

type seasonalCampaignService struct {
	repo repositories.SeasonalCampaignRepository
	cfg  *config.Config
	now  func() time.Time
}

func NewSeasonalCampaignService(repo repositories.SeasonalCampaignRepository, cfg *config.Config) SeasonalCampaignService {
	return &seasonalCampaignService{repo: repo, cfg: cfg, now: func() time.Time { return time.Now().UTC() }}
}

func (s *seasonalCampaignService) Create(ctx context.Context, req *models.CreateSeasonalCampaignRequest, adminID *uuid.UUID) (*models.SeasonalCampaignResponse, error) {
	if req.Nama == "" {
		return nil, errors.New("nama wajib diisi")
	}
	start, end, err := parseCampaignPeriod(req.TanggalMulai, req.TanggalSelesai)
	if err != nil {
		return nil, err
	}
	campaign := &models.SeasonalCampaign{ID: uuid.New(), Nama: req.Nama, TanggalMulai: start, TanggalSelesai: end, WebLogoURL: req.WebLogoURL, MobileLoadingLogoURL: req.MobileLoadingLogoURL, WebNavbarDecorationURL: req.WebNavbarDecorationURL, MobileTopAppBarOrnamentURL: req.MobileTopAppBarOrnamentURL, CreatedBy: adminID, UpdatedBy: adminID}
	if err := s.repo.Create(ctx, campaign); err != nil {
		return nil, err
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) FindByID(ctx context.Context, id string) (*models.SeasonalCampaignResponse, error) {
	campaign, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("campaign seasonal tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) FindAll(ctx context.Context, params *models.SeasonalCampaignFilterRequest) ([]models.SeasonalCampaignResponse, *models.PaginationMeta, error) {
	params.SetDefaults()
	campaigns, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	result := make([]models.SeasonalCampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		result = append(result, *s.toResponse(&campaigns[i]))
	}
	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return result, &meta, nil
}

func (s *seasonalCampaignService) Update(ctx context.Context, id string, req *models.UpdateSeasonalCampaignRequest, adminID *uuid.UUID) (*models.SeasonalCampaignResponse, error) {
	campaign, err := s.find(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Nama != nil {
		if *req.Nama == "" {
			return nil, errors.New("nama wajib diisi")
		}
		campaign.Nama = *req.Nama
	}
	if req.TanggalMulai != nil || req.TanggalSelesai != nil {
		if req.TanggalMulai != nil && req.TanggalSelesai != nil && *req.TanggalMulai == "" && *req.TanggalSelesai == "" {
			campaign.TanggalMulai, campaign.TanggalSelesai = nil, nil
		} else {
			start, end, parseErr := parseCampaignPeriod(valueOrCurrent(req.TanggalMulai, campaign.TanggalMulai), valueOrCurrent(req.TanggalSelesai, campaign.TanggalSelesai))
			if parseErr != nil {
				return nil, parseErr
			}
			campaign.TanggalMulai, campaign.TanggalSelesai = start, end
		}
	}
	if req.WebLogoURL != nil {
		campaign.WebLogoURL = req.WebLogoURL
	}
	if req.MobileLoadingLogoURL != nil {
		campaign.MobileLoadingLogoURL = req.MobileLoadingLogoURL
	}
	if req.WebNavbarDecorationURL != nil {
		campaign.WebNavbarDecorationURL = req.WebNavbarDecorationURL
	}
	if req.MobileTopAppBarOrnamentURL != nil {
		campaign.MobileTopAppBarOrnamentURL = req.MobileTopAppBarOrnamentURL
	}
	if req.RemoveWebLogo != nil && *req.RemoveWebLogo {
		campaign.WebLogoURL = nil
		campaign.MobileLoadingLogoURL = nil
	}
	if req.RemoveWebNavbarDecoration != nil && *req.RemoveWebNavbarDecoration {
		campaign.WebNavbarDecorationURL = nil
	}
	if req.RemoveMobileTopAppBarOrnament != nil && *req.RemoveMobileTopAppBarOrnament {
		campaign.MobileTopAppBarOrnamentURL = nil
	}
	if campaign.IsPublished {
		if err := s.validatePublish(ctx, campaign); err != nil {
			return nil, err
		}
	}
	campaign.UpdatedBy = adminID
	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, mapOverlapError(err)
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) Publish(ctx context.Context, id string, adminID *uuid.UUID) (*models.SeasonalCampaignResponse, error) {
	campaign, err := s.find(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.validatePublish(ctx, campaign); err != nil {
		return nil, err
	}
	campaign.IsPublished, campaign.CancelledAt, campaign.UpdatedBy = true, nil, adminID
	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, mapOverlapError(err)
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) Cancel(ctx context.Context, id string, adminID *uuid.UUID) (*models.SeasonalCampaignResponse, error) {
	campaign, err := s.find(ctx, id)
	if err != nil {
		return nil, err
	}
	if !campaign.IsPublished || campaign.CancelledAt != nil {
		return nil, errors.New("hanya campaign published atau scheduled yang dapat dibatalkan")
	}
	now := s.now()
	campaign.IsPublished, campaign.CancelledAt, campaign.UpdatedBy = false, &now, adminID
	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, err
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) Delete(ctx context.Context, id string) (*models.SeasonalCampaignResponse, error) {
	campaign, err := s.find(ctx, id)
	if err != nil {
		return nil, err
	}
	if campaign.IsPublished && campaign.CancelledAt == nil {
		return nil, errors.New("campaign published atau scheduled harus dibatalkan sebelum dihapus")
	}
	if err := s.repo.Delete(ctx, campaign); err != nil {
		return nil, err
	}
	return s.toResponse(campaign), nil
}

func (s *seasonalCampaignService) Preview(ctx context.Context, id string) (*models.SeasonalCampaignResponse, error) {
	return s.FindByID(ctx, id)
}

func (s *seasonalCampaignService) find(ctx context.Context, id string) (*models.SeasonalCampaign, error) {
	campaign, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("campaign seasonal tidak ditemukan")
	}
	return campaign, err
}

func (s *seasonalCampaignService) validatePublish(ctx context.Context, campaign *models.SeasonalCampaign) error {
	if campaign.Nama == "" {
		return errors.New("nama wajib diisi sebelum publish")
	}
	if campaign.TanggalMulai == nil || campaign.TanggalSelesai == nil {
		return errors.New("tanggal_mulai dan tanggal_selesai wajib diisi sebelum publish")
	}
	if !campaign.TanggalSelesai.After(*campaign.TanggalMulai) {
		return errors.New("tanggal_selesai harus setelah tanggal_mulai")
	}
	overlap, err := s.repo.HasPublishedOverlap(ctx, *campaign.TanggalMulai, *campaign.TanggalSelesai, campaign.ID.String())
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("periode campaign overlap dengan campaign published lain")
	}
	return nil
}

func (s *seasonalCampaignService) toResponse(c *models.SeasonalCampaign) *models.SeasonalCampaignResponse {
	return &models.SeasonalCampaignResponse{ID: c.ID.String(), Nama: c.Nama, Status: c.StatusAt(s.now()), TanggalMulai: c.TanggalMulai, TanggalSelesai: c.TanggalSelesai, Assets: models.SeasonalCampaignAssetsResponse{WebLogoURL: prefixedAssetURL(c.WebLogoURL, s.cfg.BaseURL), MobileLoadingLogoURL: prefixedAssetURL(c.MobileLoadingLogoURL, s.cfg.BaseURL), WebNavbarDecorationURL: prefixedAssetURL(c.WebNavbarDecorationURL, s.cfg.BaseURL), MobileTopAppBarOrnamentURL: prefixedAssetURL(c.MobileTopAppBarOrnamentURL, s.cfg.BaseURL)}, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

func parseCampaignPeriod(startValue, endValue *string) (*time.Time, *time.Time, error) {
	if startValue == nil && endValue == nil {
		return nil, nil, nil
	}
	if startValue == nil || endValue == nil {
		return nil, nil, errors.New("tanggal_mulai dan tanggal_selesai harus diisi bersamaan")
	}
	start, err := time.Parse(time.RFC3339, *startValue)
	if err != nil {
		return nil, nil, errors.New("tanggal_mulai harus berformat RFC3339 dengan offset")
	}
	end, err := time.Parse(time.RFC3339, *endValue)
	if err != nil {
		return nil, nil, errors.New("tanggal_selesai harus berformat RFC3339 dengan offset")
	}
	start, end = start.UTC(), end.UTC()
	if !end.After(start) {
		return nil, nil, errors.New("tanggal_selesai harus setelah tanggal_mulai")
	}
	return &start, &end, nil
}

func valueOrCurrent(value *string, current *time.Time) *string {
	if value != nil {
		return value
	}
	if current == nil {
		return nil
	}
	formatted := current.Format(time.RFC3339)
	return &formatted
}

func prefixedAssetURL(path *string, baseURL string) *string {
	if path == nil {
		return nil
	}
	if baseURL == "" {
		value := *path
		return &value
	}
	value := fmt.Sprintf("%s/uploads/%s", trimTrailingSlash(baseURL), *path)
	return &value
}

func trimTrailingSlash(value string) string {
	for len(value) > 0 && value[len(value)-1] == '/' {
		value = value[:len(value)-1]
	}
	return value
}

func mapOverlapError(err error) error {
	// PostgreSQL reports a stable constraint name; this keeps concurrent publishes user-friendly.
	if err != nil && strings.Contains(err.Error(), "seasonal_campaign_no_published_overlap") {
		return errors.New("periode campaign overlap dengan campaign published lain")
	}
	return err
}
