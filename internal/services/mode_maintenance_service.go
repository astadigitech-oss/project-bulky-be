package services

import (
	"errors"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
)

type ModeMaintenanceService interface {
	CreateMaintenance(req *models.CreateMaintenanceRequest) (*models.ModeMaintenance, error)
	UpdateMaintenance(id string, req *models.UpdateMaintenanceRequest) (*models.ModeMaintenance, error)
	DeleteMaintenance(id string) error
	GetMaintenanceByID(id string) (*models.ModeMaintenance, error)
	GetAllMaintenances(page, limit int) ([]models.ModeMaintenance, int64, error)
	ActivateMaintenance(id string) error
	DeactivateMaintenance(id string) error
}

type modeMaintenanceService struct {
	repo repositories.ModeMaintenanceRepository
}

func NewModeMaintenanceService(repo repositories.ModeMaintenanceRepository) ModeMaintenanceService {
	return &modeMaintenanceService{
		repo: repo,
	}
}

func (s *modeMaintenanceService) CreateMaintenance(req *models.CreateMaintenanceRequest) (*models.ModeMaintenance, error) {
	maintenance := &models.ModeMaintenance{
		Judul:           req.Judul,
		TipeMaintenance: models.MaintenanceType(req.TipeMaintenance),
		Deskripsi:       req.Deskripsi,
		IsActive:        req.IsActive,
	}

	err := s.repo.Create(maintenance)
	if err != nil {
		return nil, err
	}

	return maintenance, nil
}

func (s *modeMaintenanceService) UpdateMaintenance(id string, req *models.UpdateMaintenanceRequest) (*models.ModeMaintenance, error) {
	maintenance, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("maintenance not found")
	}

	if req.Judul != nil {
		maintenance.Judul = *req.Judul
	}
	if req.TipeMaintenance != nil {
		maintenance.TipeMaintenance = models.MaintenanceType(*req.TipeMaintenance)
	}
	if req.Deskripsi != nil {
		maintenance.Deskripsi = *req.Deskripsi
	}
	if req.IsActive != nil {
		maintenance.IsActive = *req.IsActive
	}

	err = s.repo.Update(maintenance)
	if err != nil {
		return nil, err
	}

	return maintenance, nil
}

func (s *modeMaintenanceService) DeleteMaintenance(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("maintenance not found")
	}

	return s.repo.Delete(id)
}

func (s *modeMaintenanceService) GetMaintenanceByID(id string) (*models.ModeMaintenance, error) {
	maintenance, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("maintenance not found")
	}
	return maintenance, nil
}

func (s *modeMaintenanceService) GetAllMaintenances(page, limit int) ([]models.ModeMaintenance, int64, error) {
	return s.repo.FindAll(page, limit)
}

func (s *modeMaintenanceService) ActivateMaintenance(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("maintenance not found")
	}

	return s.repo.Activate(id)
}

func (s *modeMaintenanceService) DeactivateMaintenance(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("maintenance not found")
	}

	return s.repo.Deactivate(id)
}
