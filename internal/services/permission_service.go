package services

import (
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type PermissionService interface {
	GetAll() ([]models.Permission, error)
	GetByModul() (map[string][]models.PermissionSimpleResponse, error)
	GetByID(id uuid.UUID) (*models.Permission, error)
}

type permissionService struct {
	repo repositories.PermissionRepository
}

func NewPermissionService(repo repositories.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

func (s *permissionService) GetAll() ([]models.Permission, error) {
	permissions, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	// Ensure empty array instead of null
	if permissions == nil {
		permissions = []models.Permission{}
	}
	return permissions, nil
}

func (s *permissionService) GetByModul() (map[string][]models.PermissionSimpleResponse, error) {
	permissions, err := s.repo.FindByModul()
	if err != nil {
		return nil, err
	}

	// Map to simple response format
	result := make(map[string][]models.PermissionSimpleResponse)
	for modul, perms := range permissions {
		simplePerms := []models.PermissionSimpleResponse{}
		for _, p := range perms {
			simplePerms = append(simplePerms, models.PermissionSimpleResponse{
				ID:        p.ID.String(),
				Nama:      p.Nama,
				Deskripsi: p.Deskripsi,
			})
		}
		result[modul] = simplePerms
	}

	return result, nil
}

func (s *permissionService) GetByID(id uuid.UUID) (*models.Permission, error) {
	return s.repo.FindByID(id)
}
