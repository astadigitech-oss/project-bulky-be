package services

import (
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type ForceUpdateService interface {
	CreateForceUpdate(req *models.CreateForceUpdateRequest) (*models.ForceUpdateApp, error)
	UpdateForceUpdate(id string, req *models.UpdateForceUpdateRequest) (*models.ForceUpdateApp, error)
	DeleteForceUpdate(id string) error
	GetForceUpdateByID(id string) (*models.ForceUpdateApp, error)
	GetAllForceUpdates(page, limit int) ([]models.ForceUpdateApp, int64, error)
	SetActiveForceUpdate(id string) error
}

type forceUpdateService struct {
	repo repositories.ForceUpdateRepository
}

func NewForceUpdateService(repo repositories.ForceUpdateRepository) ForceUpdateService {
	return &forceUpdateService{
		repo: repo,
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
