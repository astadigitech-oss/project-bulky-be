package services

import (
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type RoleService interface {
	GetAll() ([]models.Role, error)
	GetByID(id uuid.UUID) (*models.Role, error)
	GetByIDWithPermissions(id uuid.UUID) (*models.Role, error)
	Create(role *models.Role, permissionIDs []uuid.UUID) error
	Update(role *models.Role, permissionIDs []uuid.UUID) error
	Delete(id uuid.UUID) error
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAll() ([]models.Role, error) {
	return s.repo.FindAll()
}

func (s *roleService) GetByID(id uuid.UUID) (*models.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleService) GetByIDWithPermissions(id uuid.UUID) (*models.Role, error) {
	return s.repo.FindByIDWithPermissions(id)
}

func (s *roleService) Create(role *models.Role, permissionIDs []uuid.UUID) error {
	// Check if kode already exists
	existing, err := s.repo.FindByKode(role.Kode)
	if err == nil && existing != nil {
		return errors.New("kode role sudah digunakan")
	}

	// Create role
	if err := s.repo.Create(role); err != nil {
		return err
	}

	// Assign permissions
	if len(permissionIDs) > 0 {
		return s.repo.AssignPermissions(role.ID, permissionIDs)
	}

	return nil
}

func (s *roleService) Update(role *models.Role, permissionIDs []uuid.UUID) error {
	// Check if exists
	existing, err := s.repo.FindByID(role.ID)
	if err != nil {
		return errors.New("role tidak ditemukan")
	}

	// Check if kode changed and already used
	if existing.Kode != role.Kode {
		existingByKode, err := s.repo.FindByKode(role.Kode)
		if err == nil && existingByKode != nil && existingByKode.ID != role.ID {
			return errors.New("kode role sudah digunakan")
		}
	}

	// Update role
	if err := s.repo.Update(role); err != nil {
		return err
	}

	// Update permissions
	if len(permissionIDs) > 0 {
		return s.repo.AssignPermissions(role.ID, permissionIDs)
	}

	return nil
}

func (s *roleService) Delete(id uuid.UUID) error {
	// Check if role exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role tidak ditemukan")
	}

	// TODO: Check if role is being used by any admin
	// For now, just delete
	return s.repo.Delete(id)
}
