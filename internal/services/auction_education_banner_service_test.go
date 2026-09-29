package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
)

type auctionEducationBannerRepoStub struct {
	banner *models.AuctionEducationBanner
}

func (r *auctionEducationBannerRepoStub) Create(_ context.Context, banner *models.AuctionEducationBanner) error {
	r.banner = banner
	return nil
}

func (r *auctionEducationBannerRepoStub) FindByID(_ context.Context, id uuid.UUID) (*models.AuctionEducationBanner, error) {
	if r.banner == nil || r.banner.ID != id {
		return nil, gorm.ErrRecordNotFound
	}
	return r.banner, nil
}

func (r *auctionEducationBannerRepoStub) List(_ context.Context, _ string, _ int, _ int) ([]models.AuctionEducationBanner, int64, error) {
	if r.banner == nil {
		return nil, 0, nil
	}
	return []models.AuctionEducationBanner{*r.banner}, 1, nil
}

func (r *auctionEducationBannerRepoStub) Save(_ context.Context, banner *models.AuctionEducationBanner) error {
	r.banner = banner
	return nil
}

func (r *auctionEducationBannerRepoStub) Reorder(_ context.Context, _ []uuid.UUID) ([]models.AuctionEducationBanner, error) {
	return nil, nil
}

func (r *auctionEducationBannerRepoStub) Delete(_ context.Context, _ *models.AuctionEducationBanner) error {
	return nil
}

func TestAuctionEducationBannerPublishDraftLifecycle(t *testing.T) {
	repo := &auctionEducationBannerRepoStub{}
	service := NewAuctionEducationBannerService(repo, &config.Config{BaseURL: "https://api.bulky.test"})

	created, err := service.Create(context.Background(), "Cara ikut lelang", "banner/id.webp", "banner/en.webp", 0)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Status != "draft" || repo.banner.IsPublished {
		t.Fatalf("new banner must be draft, got status=%q published=%t", created.Status, repo.banner.IsPublished)
	}

	published, err := service.Publish(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if published.Status != "published" || !repo.banner.IsPublished {
		t.Fatalf("published banner has status=%q published=%t", published.Status, repo.banner.IsPublished)
	}

	draft, err := service.Draft(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Draft() error = %v", err)
	}
	if draft.Status != "draft" || repo.banner.IsPublished {
		t.Fatalf("draft banner has status=%q published=%t", draft.Status, repo.banner.IsPublished)
	}
}
