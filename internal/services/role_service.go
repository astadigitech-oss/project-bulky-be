package services

import (
	"errors"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/google/uuid"
)

type RoleService interface {
	GetAll() ([]models.RoleResponseFormat, error)
	GetByID(id uuid.UUID) (*models.Role, error)
	GetByIDWithPermissions(id uuid.UUID) (*models.RoleDetailResponse, error)
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

func (s *roleService) GetAll() ([]models.RoleResponseFormat, error) {
	roles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	// Ensure empty array instead of null
	if roles == nil {
		roles = []models.RoleResponseFormat{}
	}
	return roles, nil
}

func (s *roleService) GetByID(id uuid.UUID) (*models.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleService) GetByIDWithPermissions(id uuid.UUID) (*models.RoleDetailResponse, error) {
	role, err := s.repo.FindByIDWithPermissions(id)
	if err != nil {
		return nil, err
	}

	// Map permissions to simple response
	var permissions []models.PermissionSimpleResponse
	for _, p := range role.Permissions {
		permissions = append(permissions, models.PermissionSimpleResponse{
			ID:        p.ID.String(),
			Nama:      p.Nama,
			Deskripsi: p.Deskripsi,
		})
	}

	return &models.RoleDetailResponse{
		ID:          role.ID.String(),
		Nama:        role.Nama,
		Kode:        role.Kode,
		Deskripsi:   role.Deskripsi,
		Permissions: permissions,
	}, nil
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

	// Assign permissions (replace all, bisa kosong)
	return s.repo.AssignPermissions(role.ID, permissionIDs)
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

	// Assign permissions (replace all, bisa kosong)
	return s.repo.AssignPermissions(role.ID, permissionIDs)
}

func (s *roleService) Delete(id uuid.UUID) error {
	// Check if role exists
	role, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role tidak ditemukan")
	}

	// Guard: role bawaan tidak boleh dihapus
	switch role.Kode {
	case models.RoleSuperAdmin, models.RoleAdmin, models.RoleStaff, models.RoleFinance, models.RoleMarketing:
		return errors.New("role bawaan tidak dapat dihapus")
	}

	// Guard: role yang masih dipakai admin tidak boleh dihapus
	count, err := s.repo.CountAdminByRoleID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("role masih digunakan oleh admin, tidak dapat dihapus")
	}

	return s.repo.Delete(id)
}
