package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"
)

type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetProfile(ctx context.Context, adminID string) (*models.AdminResponse, error)
	UpdateProfile(ctx context.Context, adminID string, req *models.UpdateProfileRequest) (*models.AdminResponse, error)
	ChangePassword(ctx context.Context, adminID string, req *models.ChangePasswordRequest) error
}

type authService struct {
	adminRepo   repositories.AdminRepository
	sessionRepo repositories.AdminSessionRepository
	cfg         *config.Config
}

func NewAuthService(adminRepo repositories.AdminRepository, sessionRepo repositories.AdminSessionRepository) AuthService {
	return &authService{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
		cfg:         config.LoadConfig(),
	}
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	admin, err := s.adminRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		return nil, errors.New("email atau password salah")
	}

	if !admin.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan. Hubungi administrator")
	}

	// LEGACY AUTH SERVICE - NOT RECOMMENDED FOR NEW CODE
	// Use AuthV2Service instead which supports roles & permissions
	// This service no longer supports refresh tokens (single 24h token only)

	accessToken, err := utils.GenerateAccessToken(
		admin.ID,
		"ADMIN",
		admin.Email,
		"",  // roleID empty for legacy auth
		"",  // roleKode empty for legacy auth
		nil, // no permissions for legacy auth
	)
	if err != nil {
		return nil, err
	}

	// Update last login
	now := time.Now()
	admin.LastLoginAt = &now
	s.adminRepo.Update(ctx, admin)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: "", // No longer used - single 24h token
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours in seconds
		Admin:        *s.toAdminResponse(admin),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	// DEPRECATED: Refresh token mechanism removed
	// Return error to force re-login
	return nil, errors.New("refresh token tidak lagi didukung. Silakan login ulang")
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// DEPRECATED: Refresh token mechanism removed
	// Logout is now handled by client (remove token from storage)
	return nil
}

func (s *authService) GetProfile(ctx context.Context, adminID string) (*models.AdminResponse, error) {
	// Use FindByIDWithRole to optimize query with Preload
	admin, err := s.adminRepo.FindByIDWithRole(ctx, adminID)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}
	return s.toAdminResponse(admin), nil
}

func (s *authService) UpdateProfile(ctx context.Context, adminID string, req *models.UpdateProfileRequest) (*models.AdminResponse, error) {
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	if req.Email != nil && *req.Email != admin.Email {
		exists, _ := s.adminRepo.ExistsByEmail(ctx, *req.Email, &adminID)
		if exists {
			return nil, errors.New("email sudah digunakan oleh admin lain")
		}
		admin.Email = *req.Email
	}

	if req.Nama != nil {
		admin.Nama = *req.Nama
	}

	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return nil, err
	}

	return s.toAdminResponse(admin), nil
}

func (s *authService) ChangePassword(ctx context.Context, adminID string, req *models.ChangePasswordRequest) error {
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}

	if !utils.CheckPassword(req.OldPassword, admin.Password) {
		return errors.New("password lama tidak sesuai")
	}

	if utils.CheckPassword(req.NewPassword, admin.Password) {
		return errors.New("password baru tidak boleh sama dengan password lama")
	}

	// Use configured bcrypt cost
	hashedPassword, err := utils.HashPasswordWithCost(req.NewPassword, s.cfg.BcryptCost)
	if err != nil {
		return err
	}

	admin.Password = hashedPassword
	return s.adminRepo.Update(ctx, admin)
}

func (s *authService) toAdminResponse(a *models.Admin) *models.AdminResponse {
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
