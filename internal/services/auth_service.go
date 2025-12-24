package services

import (
	"context"
	"errors"
	"time"

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
}

func NewAuthService(adminRepo repositories.AdminRepository, sessionRepo repositories.AdminSessionRepository) AuthService {
	return &authService{adminRepo: adminRepo, sessionRepo: sessionRepo}
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

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(admin.ID.String(), admin.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := utils.GenerateRefreshToken(admin.ID.String(), admin.Email)
	if err != nil {
		return nil, err
	}

	// Save session
	session := &models.AdminSession{
		AdminID:   admin.ID,
		Token:     refreshToken,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		ExpiresAt: expiresAt,
	}
	s.sessionRepo.Create(ctx, session)

	// Update last login
	now := time.Now()
	admin.LastLoginAt = &now
	s.adminRepo.Update(ctx, admin)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Admin:        *s.toAdminResponse(admin),
	}, nil
}


func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateJWT(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau sudah expired. Silakan login ulang")
	}

	// Check if session exists
	session, err := s.sessionRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau sudah expired. Silakan login ulang")
	}

	// Check if expired
	if session.ExpiresAt.Before(time.Now()) {
		s.sessionRepo.DeleteByToken(ctx, refreshToken)
		return nil, errors.New("refresh token tidak valid atau sudah expired. Silakan login ulang")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(claims.AdminID, claims.Email)
	if err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.sessionRepo.DeleteByToken(ctx, refreshToken)
}

func (s *authService) GetProfile(ctx context.Context, adminID string) (*models.AdminResponse, error) {
	admin, err := s.adminRepo.FindByID(ctx, adminID)
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

	hashedPassword, err := utils.HashPassword(req.NewPassword)
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
