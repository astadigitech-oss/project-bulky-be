package services

import (
	"context"
	"errors"
	"log"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoService interface {
	Create(ctx context.Context, req *dto.CreateVideoRequest) (*dto.VideoResponse, error)
	CreateDraft(ctx context.Context, req *dto.CreateVideoRequest) (*models.Video, error)
	MarkReady(ctx context.Context, id uuid.UUID, videoURL string, durasiDetik int) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	RecoverStuckJobs(ctx context.Context)
	GetTranscodeStatus(ctx context.Context, id uuid.UUID) (*dto.VideoTranscodeStatusResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateVideoRequest) (*dto.VideoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.VideoResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.VideoResponse, error)
	GetAll(ctx context.Context, params *dto.VideoFilterRequest) ([]dto.VideoListResponse, *models.PaginationMeta, error)
	Search(ctx context.Context, keyword string, isActive *bool, page, limit int) ([]dto.VideoListResponse, int64, error)
	GetPopular(ctx context.Context, limit int) ([]dto.VideoListResponse, error)
	IncrementView(ctx context.Context, id uuid.UUID) error
	GetRelated(ctx context.Context, videoID uuid.UUID, limit int) ([]dto.VideoListResponse, error)
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) error
}

type videoService struct {
	videoRepo    repositories.VideoRepository
	kategoriRepo repositories.KategoriVideoRepository
	cfg          *config.Config
}

func NewVideoService(
	videoRepo repositories.VideoRepository,
	kategoriRepo repositories.KategoriVideoRepository,
	cfg *config.Config,
) VideoService {
	return &videoService{
		videoRepo:    videoRepo,
		kategoriRepo: kategoriRepo,
		cfg:          cfg,
	}
}

func (s *videoService) Create(ctx context.Context, req *dto.CreateVideoRequest) (*dto.VideoResponse, error) {
	// Validate kategori exists
	_, err := s.kategoriRepo.FindByID(ctx, req.KategoriID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kategori not found")
		}
		return nil, err
	}

	// Convert is_active string to boolean
	isActive := req.IsActive == "true"

	// Generate/use slug_id
	var slugIDVal string
	if req.SlugID != nil && *req.SlugID != "" {
		slugIDVal = *req.SlugID
	} else {
		slugIDVal = utils.GenerateSlug(req.JudulID)
	}
	slugIDPtr := &slugIDVal

	// Generate/use slug_en
	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		slugEN = req.SlugEN
	} else if req.JudulEN != "" {
		s := utils.GenerateSlug(req.JudulEN)
		slugEN = &s
	}

	// Auto-generate meta SEO fields when not provided
	metaTitleID := req.MetaTitleID
	if metaTitleID == nil {
		v := truncateMeta(req.JudulID, 200)
		metaTitleID = &v
	}

	metaTitleEN := req.MetaTitleEN
	if metaTitleEN == nil && req.JudulEN != "" {
		v := truncateMeta(req.JudulEN, 200)
		metaTitleEN = &v
	}

	metaDescriptionID := req.MetaDescriptionID
	if metaDescriptionID == nil {
		v := truncateMeta(stripHTML(req.DeskripsiID), 160)
		metaDescriptionID = &v
	}

	metaDescriptionEN := req.MetaDescriptionEN
	if metaDescriptionEN == nil && req.DeskripsiEN != "" {
		v := truncateMeta(stripHTML(req.DeskripsiEN), 160)
		metaDescriptionEN = &v
	}

	metaKeywords := req.MetaKeywords
	if metaKeywords == nil {
		judulEN := &req.JudulEN
		kw := generateMetaKeywords(req.JudulID, judulEN, nil)
		metaKeywords = &kw
	}

	video := &models.Video{
		JudulID:           req.JudulID,
		JudulEN:           &req.JudulEN,
		Slug:              slugIDVal,
		SlugID:            slugIDPtr,
		SlugEN:            slugEN,
		DeskripsiID:       req.DeskripsiID,
		DeskripsiEN:       &req.DeskripsiEN,
		VideoURL:          req.VideoURL,
		ThumbnailURL:      req.ThumbnailURL,
		KategoriID:        req.KategoriID,
		MetaTitleID:       metaTitleID,
		MetaTitleEN:       metaTitleEN,
		MetaDescriptionID: metaDescriptionID,
		MetaDescriptionEN: metaDescriptionEN,
		MetaKeywords:      metaKeywords,
		IsActive:          isActive,
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, video.ID)
}

func (s *videoService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateVideoRequest) (*dto.VideoResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.JudulID != nil {
		video.JudulID = *req.JudulID
		if req.SlugID == nil {
			s := utils.GenerateSlug(*req.JudulID)
			video.SlugID = &s
			video.Slug = s
		}
	}
	if req.JudulEN != nil {
		video.JudulEN = req.JudulEN
		if req.SlugEN == nil && req.JudulEN != nil && *req.JudulEN != "" {
			s := utils.GenerateSlug(*req.JudulEN)
			video.SlugEN = &s
		}
	}
	if req.SlugID != nil && *req.SlugID != "" {
		video.SlugID = req.SlugID
		video.Slug = *req.SlugID
	}
	if req.SlugEN != nil && *req.SlugEN != "" {
		video.SlugEN = req.SlugEN
	}
	if req.DeskripsiID != nil {
		video.DeskripsiID = *req.DeskripsiID
	}
	if req.DeskripsiEN != nil {
		video.DeskripsiEN = req.DeskripsiEN
	}
	if req.VideoURL != nil {
		video.VideoURL = *req.VideoURL
	}
	if req.ThumbnailURL != nil {
		video.ThumbnailURL = req.ThumbnailURL
	}
	if req.KategoriID != nil {
		video.KategoriID = *req.KategoriID
	}
	// Meta SEO: gunakan nilai eksplisit jika ada, auto-generate jika source berubah
	if req.MetaTitleID != nil {
		video.MetaTitleID = req.MetaTitleID
	} else if req.JudulID != nil {
		v := truncateMeta(video.JudulID, 200)
		video.MetaTitleID = &v
	}

	if req.MetaTitleEN != nil {
		video.MetaTitleEN = req.MetaTitleEN
	} else if req.JudulEN != nil && *req.JudulEN != "" {
		v := truncateMeta(*req.JudulEN, 200)
		video.MetaTitleEN = &v
	}

	if req.MetaDescriptionID != nil {
		video.MetaDescriptionID = req.MetaDescriptionID
	} else if req.DeskripsiID != nil {
		v := truncateMeta(stripHTML(video.DeskripsiID), 160)
		video.MetaDescriptionID = &v
	}

	if req.MetaDescriptionEN != nil {
		video.MetaDescriptionEN = req.MetaDescriptionEN
	} else if req.DeskripsiEN != nil && video.DeskripsiEN != nil {
		v := truncateMeta(stripHTML(*video.DeskripsiEN), 160)
		video.MetaDescriptionEN = &v
	}

	if req.MetaKeywords != nil {
		video.MetaKeywords = req.MetaKeywords
	} else if req.JudulID != nil || req.JudulEN != nil {
		kw := generateMetaKeywords(video.JudulID, video.JudulEN, nil)
		video.MetaKeywords = &kw
	}
	if req.IsActive != nil {
		// Convert string to boolean
		isActive := *req.IsActive == "true"
		video.IsActive = isActive
	}

	if err := s.videoRepo.Update(ctx, video); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *videoService) Delete(ctx context.Context, id uuid.UUID) error {
	video, err := s.videoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.videoRepo.Delete(ctx, video)
}

func (s *videoService) GetByID(ctx context.Context, id uuid.UUID) (*dto.VideoResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toVideoResponse(video), nil
}

func (s *videoService) GetBySlug(ctx context.Context, slug string) (*dto.VideoResponse, error) {
	video, err := s.videoRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return s.toVideoResponse(video), nil
}

func (s *videoService) GetAll(ctx context.Context, params *dto.VideoFilterRequest) ([]dto.VideoListResponse, *models.PaginationMeta, error) {
	offset := params.GetOffset()
	videos, total, err := s.videoRepo.FindAll(ctx, params.Search, params.IsActive, params.KategoriID, params.SortBy, params.Order, params.PerPage, offset)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.VideoListResponse, len(videos))
	for i, video := range videos {
		responses[i] = s.toVideoListResponse(&video)
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)
	return responses, &meta, nil
}

func (s *videoService) Search(ctx context.Context, keyword string, isActive *bool, page, limit int) ([]dto.VideoListResponse, int64, error) {
	offset := (page - 1) * limit
	videos, total, err := s.videoRepo.Search(ctx, keyword, isActive, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.VideoListResponse, len(videos))
	for i, video := range videos {
		responses[i] = s.toVideoListResponse(&video)
	}

	return responses, total, nil
}

func (s *videoService) GetPopular(ctx context.Context, limit int) ([]dto.VideoListResponse, error) {
	videos, err := s.videoRepo.FindPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VideoListResponse, len(videos))
	for i, video := range videos {
		responses[i] = s.toVideoListResponse(&video)
	}

	return responses, nil
}

func (s *videoService) IncrementView(ctx context.Context, id uuid.UUID) error {
	return s.videoRepo.IncrementViewCount(ctx, id)
}

func (s *videoService) GetRelated(ctx context.Context, videoID uuid.UUID, limit int) ([]dto.VideoListResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	videos, err := s.videoRepo.FindRelated(ctx, videoID, video.KategoriID, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VideoListResponse, len(videos))
	for i, v := range videos {
		responses[i] = s.toVideoListResponse(&v)
	}

	return responses, nil
}

func (s *videoService) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats, err := s.videoRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// Get popular videos
	popularVideos, _ := s.videoRepo.FindPopular(ctx, 5)
	stats["video_populer"] = popularVideos

	return stats, nil
}

func (s *videoService) ToggleStatus(ctx context.Context, id uuid.UUID) error {
	return s.videoRepo.ToggleStatus(ctx, id)
}

func (s *videoService) CreateDraft(ctx context.Context, req *dto.CreateVideoRequest) (*models.Video, error) {
	_, err := s.kategoriRepo.FindByID(ctx, req.KategoriID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kategori not found")
		}
		return nil, err
	}

	isActive := req.IsActive == "true"

	var slugIDVal string
	if req.SlugID != nil && *req.SlugID != "" {
		slugIDVal = *req.SlugID
	} else {
		slugIDVal = utils.GenerateSlug(req.JudulID)
	}
	slugIDPtr := &slugIDVal

	var slugEN *string
	if req.SlugEN != nil && *req.SlugEN != "" {
		slugEN = req.SlugEN
	} else if req.JudulEN != "" {
		s := utils.GenerateSlug(req.JudulEN)
		slugEN = &s
	}

	metaTitleID := req.MetaTitleID
	if metaTitleID == nil {
		v := truncateMeta(req.JudulID, 200)
		metaTitleID = &v
	}
	metaTitleEN := req.MetaTitleEN
	if metaTitleEN == nil && req.JudulEN != "" {
		v := truncateMeta(req.JudulEN, 200)
		metaTitleEN = &v
	}
	metaDescriptionID := req.MetaDescriptionID
	if metaDescriptionID == nil {
		v := truncateMeta(stripHTML(req.DeskripsiID), 160)
		metaDescriptionID = &v
	}
	metaDescriptionEN := req.MetaDescriptionEN
	if metaDescriptionEN == nil && req.DeskripsiEN != "" {
		v := truncateMeta(stripHTML(req.DeskripsiEN), 160)
		metaDescriptionEN = &v
	}
	metaKeywords := req.MetaKeywords
	if metaKeywords == nil {
		judulEN := &req.JudulEN
		kw := generateMetaKeywords(req.JudulID, judulEN, nil)
		metaKeywords = &kw
	}

	video := &models.Video{
		JudulID:           req.JudulID,
		JudulEN:           &req.JudulEN,
		Slug:              slugIDVal,
		SlugID:            slugIDPtr,
		SlugEN:            slugEN,
		DeskripsiID:       req.DeskripsiID,
		DeskripsiEN:       &req.DeskripsiEN,
		VideoURL:          req.VideoURL,
		ThumbnailURL:      req.ThumbnailURL,
		KategoriID:        req.KategoriID,
		MetaTitleID:       metaTitleID,
		MetaTitleEN:       metaTitleEN,
		MetaDescriptionID: metaDescriptionID,
		MetaDescriptionEN: metaDescriptionEN,
		MetaKeywords:      metaKeywords,
		IsActive:          isActive,
		TranscodeStatus:   "processing",
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	return video, nil
}

func (s *videoService) MarkReady(ctx context.Context, id uuid.UUID, videoURL string, durasiDetik int) error {
	return s.videoRepo.UpdateFields(ctx, id, map[string]interface{}{
		"video_url":        videoURL,
		"durasi_detik":     durasiDetik,
		"transcode_status": "ready",
		"transcode_error":  gorm.Expr("NULL"),
	})
}

func (s *videoService) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return s.videoRepo.UpdateFields(ctx, id, map[string]interface{}{
		"transcode_status": "failed",
		"transcode_error":  errMsg,
	})
}

func (s *videoService) RecoverStuckJobs(ctx context.Context) {
	affected, err := s.videoRepo.MarkStuckAsFailed(ctx, 15*time.Minute)
	if err != nil {
		log.Printf("[video] RecoverStuckJobs error: %v", err)
		return
	}
	if affected > 0 {
		log.Printf("[video] RecoverStuckJobs: %d video dikembalikan ke status failed", affected)
	}
}

func (s *videoService) GetTranscodeStatus(ctx context.Context, id uuid.UUID) (*dto.VideoTranscodeStatusResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.VideoTranscodeStatusResponse{
		ID:              video.ID,
		TranscodeStatus: video.TranscodeStatus,
		TranscodeError:  video.TranscodeError,
	}, nil
}

func (s *videoService) toVideoResponse(video *models.Video) *dto.VideoResponse {
	resp := &dto.VideoResponse{
		ID:                video.ID,
		JudulID:           video.JudulID,
		JudulEN:           video.JudulEN,
		SlugID:            video.SlugID,
		SlugEN:            video.SlugEN,
		DeskripsiID:       video.DeskripsiID,
		DeskripsiEN:       video.DeskripsiEN,
		VideoURL:          utils.GetFileURL(video.VideoURL, s.cfg),
		ThumbnailURL:      utils.GetFileURLPtr(video.ThumbnailURL, s.cfg),
		KategoriID:        video.KategoriID,
		MetaTitleID:       video.MetaTitleID,
		MetaTitleEN:       video.MetaTitleEN,
		MetaDescriptionID: video.MetaDescriptionID,
		MetaDescriptionEN: video.MetaDescriptionEN,
		MetaKeywords:      video.MetaKeywords,
		TranscodeStatus:   video.TranscodeStatus,
		TranscodeError:    video.TranscodeError,
		IsActive:          video.IsActive,
		ViewCount:         video.ViewCount,
		PublishedAt:       video.PublishedAt,
		CreatedAt:         video.CreatedAt,
		UpdatedAt:         video.UpdatedAt,
	}

	if video.Kategori != nil {
		resp.Kategori = &dto.KategoriVideoBrief{
			ID:     video.Kategori.ID,
			NamaID: video.Kategori.NamaID,
			NamaEN: video.Kategori.NamaEN,
			SlugID: video.Kategori.SlugID,
			SlugEN: video.Kategori.SlugEN,
		}
	}

	return resp
}

func (s *videoService) toVideoListResponse(video *models.Video) dto.VideoListResponse {
	resp := dto.VideoListResponse{
		ID:              video.ID,
		JudulID:         video.JudulID,
		JudulEN:         video.JudulEN,
		SlugID:          video.SlugID,
		SlugEN:          video.SlugEN,
		DeskripsiID:     video.DeskripsiID,
		DeskripsiEN:     video.DeskripsiEN,
		ThumbnailURL:    utils.GetFileURLPtr(video.ThumbnailURL, s.cfg),
		Kategori:        nil,
		TranscodeStatus: video.TranscodeStatus,
		IsActive:        video.IsActive,
		ViewCount:       video.ViewCount,
		PublishedAt:     video.PublishedAt,
		CreatedAt:       video.CreatedAt,
	}

	if video.Kategori != nil {
		resp.Kategori = &dto.KategoriVideoBrief{
			ID:     video.Kategori.ID,
			NamaID: video.Kategori.NamaID,
			NamaEN: video.Kategori.NamaEN,
			SlugID: video.Kategori.SlugID,
			SlugEN: video.Kategori.SlugEN,
		}
	}

	return resp
}
