package services

import (
	"context"
	"errors"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/google/uuid"
)

type AdminService interface {
	Create(ctx context.Context, req *models.CreateAdminRequest) (*models.AdminResponse, error)
	FindByID(ctx context.Context, id string) (*models.AdminResponse, error)
	FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.AdminListResponse, *models.PaginationMeta, error)
	Update(ctx context.Context, id string, req *models.UpdateAdminRequest) (*models.AdminResponse, error)
	Delete(ctx context.Context, id, currentAdminID string) error
	ToggleStatus(ctx context.Context, id, currentAdminID string) (*models.ToggleStatusResponse, error)
	ResetPassword(ctx context.Context, id string, req *models.ResetPasswordRequest) error
	UpdateProfile(ctx context.Context, id string, nama, email string) (*models.Admin, error)
	IsEmailExistExcludeID(ctx context.Context, email, excludeID string) (bool, error)
}

type adminService struct {
	repo        repositories.AdminRepository
	sessionRepo repositories.AdminSessionRepository
	cfg         *config.Config
}

func NewAdminService(repo repositories.AdminRepository, sessionRepo repositories.AdminSessionRepository) AdminService {
	return &adminService{
		repo:        repo,
		sessionRepo: sessionRepo,
		cfg:         config.LoadConfig(),
	}
}

func (s *adminService) Create(ctx context.Context, req *models.CreateAdminRequest) (*models.AdminResponse, error) {
	exists, _ := s.repo.ExistsByEmail(ctx, req.Email, nil)
	if exists {
		return nil, errors.New("email sudah terdaftar")
	}

	// Parse RoleID
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return nil, errors.New("role_id tidak valid")
	}

	// Use configured bcrypt cost
	hashedPassword, err := utils.HashPasswordWithCost(req.Password, s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	admin := &models.Admin{
		ID:       uuid.New(),
		Nama:     req.Nama,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   roleID,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return s.toResponse(admin), nil
}

func (s *adminService) FindByID(ctx context.Context, id string) (*models.AdminResponse, error) {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}
	return s.toResponse(admin), nil
}

func (s *adminService) FindAll(ctx context.Context, params *models.PaginationRequest) ([]models.AdminListResponse, *models.PaginationMeta, error) {
	params.SetDefaults()

	admins, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var items []models.AdminListResponse
	for _, a := range admins {
		roleName := ""
		if a.Role != nil {
			roleName = a.Role.Nama
		}

		items = append(items, models.AdminListResponse{
			ID:          a.ID.String(),
			Nama:        a.Nama,
			Email:       a.Email,
			Role:        roleName,
			IsActive:    a.IsActive,
			LastLoginAt: a.LastLoginAt,
			CreatedAt:   a.CreatedAt,
		})
	}

	meta := models.NewPaginationMeta(params.Page, params.PerPage, total)

	return items, &meta, nil
}

func (s *adminService) Update(ctx context.Context, id string, req *models.UpdateAdminRequest) (*models.AdminResponse, error) {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	if req.Email != nil && *req.Email != admin.Email {
		exists, err := s.repo.ExistsByEmail(ctx, *req.Email, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email sudah digunakan oleh admin lain")
		}
		admin.Email = *req.Email
	}

	if req.Nama != nil {
		admin.Nama = *req.Nama
	}
	if req.IsActive != nil {
		admin.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	return s.toResponse(admin), nil
}

func (s *adminService) Delete(ctx context.Context, id, currentAdminID string) error {
	if id == currentAdminID {
		return errors.New("tidak dapat menghapus akun sendiri")
	}

	count, _ := s.repo.Count(ctx)
	if count <= 1 {
		return errors.New("tidak dapat menghapus admin terakhir")
	}

	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}

	// Delete all sessions
	s.sessionRepo.DeleteByAdminID(ctx, id)

	return s.repo.Delete(ctx, id)
}

func (s *adminService) ToggleStatus(ctx context.Context, id, currentAdminID string) (*models.ToggleStatusResponse, error) {
	if id == currentAdminID {
		return nil, errors.New("tidak dapat menonaktifkan akun sendiri")
	}

	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	admin.IsActive = !admin.IsActive
	if err := s.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	// If deactivated, delete all sessions
	if !admin.IsActive {
		s.sessionRepo.DeleteByAdminID(ctx, id)
	}

	return &models.ToggleStatusResponse{
		ID:       admin.ID.String(),
		IsActive: admin.IsActive,
	}, nil
}

func (s *adminService) ResetPassword(ctx context.Context, id string, req *models.ResetPasswordRequest) error {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}

	// Use configured bcrypt cost
	hashedPassword, err := utils.HashPasswordWithCost(req.NewPassword, s.cfg.BcryptCost)
	if err != nil {
		return err
	}

	admin.Password = hashedPassword
	if err := s.repo.Update(ctx, admin); err != nil {
		return err
	}

	// Delete all sessions to force re-login
	s.sessionRepo.DeleteByAdminID(ctx, id)

	return nil
}

func (s *adminService) toResponse(a *models.Admin) *models.AdminResponse {
	return &models.AdminResponse{
		ID:          a.ID.String(),
		Nama:        a.Nama,
		Email:       a.Email,
		IsActive:    a.IsActive,
		LastLoginAt: a.LastLoginAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// UpdateProfile updates admin profile (nama and email)
func (s *adminService) UpdateProfile(ctx context.Context, id string, nama, email string) (*models.Admin, error) {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	admin.Nama = nama
	admin.Email = email

	if err := s.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

// IsEmailExistExcludeID checks if email exists excluding specific ID
func (s *adminService) IsEmailExistExcludeID(ctx context.Context, email, excludeID string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email, &excludeID)
}
