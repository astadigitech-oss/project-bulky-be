package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
)

type seasonalCampaignRepoStub struct {
	campaign *models.SeasonalCampaign
	overlap  bool
	saved    int
}

func (r *seasonalCampaignRepoStub) Create(_ context.Context, campaign *models.SeasonalCampaign) error {
	r.campaign = campaign
	return nil
}
func (r *seasonalCampaignRepoStub) FindByID(_ context.Context, _ string) (*models.SeasonalCampaign, error) {
	if r.campaign == nil {
		return nil, errors.New("record not found")
	}
	return r.campaign, nil
}
func (r *seasonalCampaignRepoStub) FindAll(_ context.Context, _ *models.SeasonalCampaignFilterRequest) ([]models.SeasonalCampaign, int64, error) {
	if r.campaign == nil {
		return nil, 0, nil
	}
	return []models.SeasonalCampaign{*r.campaign}, 1, nil
}
func (r *seasonalCampaignRepoStub) Save(_ context.Context, campaign *models.SeasonalCampaign) error {
	r.campaign, r.saved = campaign, r.saved+1
	return nil
}
func (r *seasonalCampaignRepoStub) Delete(_ context.Context, _ *models.SeasonalCampaign) error {
	return nil
}
func (r *seasonalCampaignRepoStub) HasPublishedOverlap(_ context.Context, _, _ time.Time, _ string) (bool, error) {
	return r.overlap, nil
}

func newSeasonalServiceForTest(repo *seasonalCampaignRepoStub, now time.Time) *seasonalCampaignService {
	service := NewSeasonalCampaignService(repo, &config.Config{BaseURL: "https://api.bulky.test/"}).(*seasonalCampaignService)
	service.now = func() time.Time { return now }
	return service
}

func campaignForTest(start, end time.Time) *models.SeasonalCampaign {
	return &models.SeasonalCampaign{ID: uuid.New(), Nama: "HUT RI", TanggalMulai: &start, TanggalSelesai: &end}
}

func TestSeasonalCampaignPublishRejectsOverlap(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	repo := &seasonalCampaignRepoStub{campaign: campaignForTest(start, end)}
	service := newSeasonalServiceForTest(repo, start)
	repo.overlap = true
	if _, err := service.Publish(context.Background(), repo.campaign.ID.String(), nil); err == nil {
		t.Fatal("publish with overlapping period must fail")
	}
}

func TestSeasonalCampaignCreateDraftWithoutAssetsAndRejectsInvalidPeriod(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	repo := &seasonalCampaignRepoStub{}
	service := newSeasonalServiceForTest(repo, now)
	result, err := service.Create(context.Background(), &models.CreateSeasonalCampaignRequest{Nama: "Promo September"}, nil)
	if err != nil {
		t.Fatalf("draft without assets should be valid: %v", err)
	}
	if result.Status != "draft" || result.Assets.WebLogoURL != nil || result.Assets.WebNavbarDecorationURL != nil || result.Assets.MobileTopAppBarOrnamentURL != nil {
		t.Fatalf("unexpected draft response: %+v", result)
	}
	start, end := "2026-09-02T00:00:00+07:00", "2026-09-01T00:00:00+07:00"
	if _, err := service.Create(context.Background(), &models.CreateSeasonalCampaignRequest{Nama: "Invalid", TanggalMulai: &start, TanggalSelesai: &end}, nil); err == nil {
		t.Fatal("campaign with an end before start must be rejected")
	}
}

func TestSeasonalCampaignPublishAndCancelLifecycle(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	start, end := now.Add(-time.Hour), now.Add(time.Hour)
	repo := &seasonalCampaignRepoStub{campaign: campaignForTest(start, end)}
	service := newSeasonalServiceForTest(repo, now)
	result, err := service.Publish(context.Background(), repo.campaign.ID.String(), nil)
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	if result.Status != "active" || !repo.campaign.IsPublished || repo.campaign.CancelledAt != nil {
		t.Fatalf("unexpected publish state: %+v", result)
	}
	result, err = service.Cancel(context.Background(), repo.campaign.ID.String(), nil)
	if err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	if result.Status != "cancelled" || repo.campaign.IsPublished || repo.campaign.CancelledAt == nil {
		t.Fatalf("unexpected cancel state: %+v", result)
	}
}

func TestSeasonalCampaignUpdateOneAssetKeepsOtherAssets(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	start, end := now.Add(-time.Hour), now.Add(time.Hour)
	logo, navbar, mobile := "seasonal-campaign/logo.webp", "seasonal-campaign/navbar.webp", "seasonal-campaign/mobile.webp"
	repo := &seasonalCampaignRepoStub{campaign: campaignForTest(start, end)}
	repo.campaign.WebLogoURL, repo.campaign.WebNavbarDecorationURL, repo.campaign.MobileTopAppBarOrnamentURL = &logo, &navbar, &mobile
	service := newSeasonalServiceForTest(repo, now)
	newLogo := "seasonal-campaign/new-logo.webp"
	result, err := service.Update(context.Background(), repo.campaign.ID.String(), &models.UpdateSeasonalCampaignRequest{WebLogoURL: &newLogo}, nil)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if *repo.campaign.WebLogoURL != newLogo || *repo.campaign.WebNavbarDecorationURL != navbar || *repo.campaign.MobileTopAppBarOrnamentURL != mobile {
		t.Fatal("single asset update must preserve assets not supplied")
	}
	if got := *result.Assets.WebLogoURL; got != "https://api.bulky.test/uploads/seasonal-campaign/new-logo.webp" {
		t.Fatalf("asset URL was not prefixed correctly: %s", got)
	}
	remove := true
	result, err = service.Update(context.Background(), repo.campaign.ID.String(), &models.UpdateSeasonalCampaignRequest{RemoveWebNavbarDecoration: &remove}, nil)
	if err != nil {
		t.Fatalf("remove asset failed: %v", err)
	}
	if result.Assets.WebNavbarDecorationURL != nil || *repo.campaign.MobileTopAppBarOrnamentURL != mobile {
		t.Fatal("remove flag must only clear the chosen asset")
	}
}

func TestSeasonalCampaignUpdateCanClearDraftPeriod(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	start, end := now.Add(time.Hour), now.Add(2*time.Hour)
	logo := "seasonal-campaign/logo.webp"
	repo := &seasonalCampaignRepoStub{campaign: campaignForTest(start, end)}
	repo.campaign.WebLogoURL = &logo
	service := newSeasonalServiceForTest(repo, now)
	empty := ""

	result, err := service.Update(context.Background(), repo.campaign.ID.String(), &models.UpdateSeasonalCampaignRequest{
		TanggalMulai:   &empty,
		TanggalSelesai: &empty,
	}, nil)
	if err != nil {
		t.Fatalf("clearing a draft period failed: %v", err)
	}
	if repo.campaign.TanggalMulai != nil || repo.campaign.TanggalSelesai != nil {
		t.Fatalf("period was not cleared: %+v", repo.campaign)
	}
	if result.TanggalMulai != nil || result.TanggalSelesai != nil || result.Assets.WebLogoURL == nil {
		t.Fatalf("unexpected response after clearing period: %+v", result)
	}
}

func TestSeasonalCampaignStatusAt(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	start, end := now.Add(time.Hour), now.Add(2*time.Hour)
	campaign := campaignForTest(start, end)
	if status := campaign.StatusAt(now); status != "draft" {
		t.Fatalf("got %s, want draft", status)
	}
	campaign.IsPublished = true
	if status := campaign.StatusAt(now); status != "scheduled" {
		t.Fatalf("got %s, want scheduled", status)
	}
	start, end = now.Add(-2*time.Hour), now.Add(-time.Hour)
	campaign.TanggalMulai, campaign.TanggalSelesai = &start, &end
	if status := campaign.StatusAt(now); status != "ended" {
		t.Fatalf("got %s, want ended", status)
	}
}
