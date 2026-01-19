package services

import (
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type ForceUpdateService interface {
	CreateForceUpdate(req *models.CreateForceUpdateRequest) (*models.ForceUpdateApp, error)
	UpdateForceUpdate(id string, req *models.UpdateForceUpdateRequest) (*models.ForceUpdateApp, error)
	DeleteForceUpdate(id string) error
	GetForceUpdateByID(id string) (*models.ForceUpdateApp, error)
	GetAllForceUpdates(page, limit int) ([]models.ForceUpdateApp, int64, error)
	SetActiveForceUpdate(id string) error
	CheckVersion(currentVersion string) (*models.CheckVersionResponse, error)
}

type forceUpdateService struct {
	repo         repositories.ForceUpdateRepository
	playStoreURL string
	appStoreURL  string
}

func NewForceUpdateService(repo repositories.ForceUpdateRepository, playStoreURL, appStoreURL string) ForceUpdateService {
	return &forceUpdateService{
		repo:         repo,
		playStoreURL: playStoreURL,
		appStoreURL:  appStoreURL,
	}
}

func (s *forceUpdateService) CreateForceUpdate(req *models.CreateForceUpdateRequest) (*models.ForceUpdateApp, error) {
	forceUpdate := &models.ForceUpdateApp{
		KodeVersi:       req.KodeVersi,
		UpdateType:      models.UpdateType(req.UpdateType),
		InformasiUpdate: req.InformasiUpdate,
		IsActive:        req.IsActive,
	}

	err := s.repo.Create(forceUpdate)
	if err != nil {
		return nil, err
	}

	return forceUpdate, nil
}

func (s *forceUpdateService) UpdateForceUpdate(id string, req *models.UpdateForceUpdateRequest) (*models.ForceUpdateApp, error) {
	forceUpdate, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("force update not found")
	}

	if req.KodeVersi != nil {
		forceUpdate.KodeVersi = *req.KodeVersi
	}
	if req.UpdateType != nil {
		forceUpdate.UpdateType = models.UpdateType(*req.UpdateType)
	}
	if req.InformasiUpdate != nil {
		forceUpdate.InformasiUpdate = *req.InformasiUpdate
	}
	if req.IsActive != nil {
		forceUpdate.IsActive = *req.IsActive
	}

	err = s.repo.Update(forceUpdate)
	if err != nil {
		return nil, err
	}

	return forceUpdate, nil
}

func (s *forceUpdateService) DeleteForceUpdate(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("force update not found")
	}

	return s.repo.Delete(id)
}

func (s *forceUpdateService) GetForceUpdateByID(id string) (*models.ForceUpdateApp, error) {
	forceUpdate, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("force update not found")
	}
	return forceUpdate, nil
}

func (s *forceUpdateService) GetAllForceUpdates(page, limit int) ([]models.ForceUpdateApp, int64, error) {
	return s.repo.FindAll(page, limit)
}

func (s *forceUpdateService) SetActiveForceUpdate(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("force update not found")
	}

	return s.repo.SetActive(id)
}

func (s *forceUpdateService) CheckVersion(currentVersion string) (*models.CheckVersionResponse, error) {
	activeUpdate, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}

	// No active force update configured
	if activeUpdate == nil {
		return &models.CheckVersionResponse{
			ShouldUpdate:    false,
			UpdateType:      nil,
			LatestVersion:   nil,
			CurrentVersion:  currentVersion,
			InformasiUpdate: nil,
			StoreURL:        nil,
		}, nil
	}

	needsUpdate := utils.NeedsUpdate(currentVersion, activeUpdate.KodeVersi)
	updateTypeStr := string(activeUpdate.UpdateType)

	response := &models.CheckVersionResponse{
		ShouldUpdate:    needsUpdate,
		CurrentVersion:  currentVersion,
		LatestVersion:   &activeUpdate.KodeVersi,
		UpdateType:      &updateTypeStr,
		InformasiUpdate: &activeUpdate.InformasiUpdate,
		StoreURL:        &s.playStoreURL,
	}

	return response, nil
}
