package services

import (
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type PermissionService interface {
	GetAll() ([]models.Permission, error)
	GetByModul() (map[string][]models.Permission, error)
	GetByID(id uuid.UUID) (*models.Permission, error)
}

type permissionService struct {
	repo repositories.PermissionRepository
}

func NewPermissionService(repo repositories.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) GetAll() ([]models.Permission, error) {
	return s.repo.FindAll()
}

func (s *permissionService) GetByModul() (map[string][]models.Permission, error) {
	return s.repo.FindByModul()
}

func (s *permissionService) GetByID(id uuid.UUID) (*models.Permission, error) {
	return s.repo.FindByID(id)
}
