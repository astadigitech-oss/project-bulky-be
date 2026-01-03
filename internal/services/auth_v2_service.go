package services

import (
	"context"
	"errors"
	"time"

	"project-bulky-be/internal/config"
	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"
	"project-bulky-be/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthV2Service interface {
	// Authentication - Separated by user type
	AdminLogin(ctx context.Context, email, password string) (*LoginResultSimplified, error)
	BuyerLogin(ctx context.Context, email, password string) (*LoginResultSimplified, error)

	// Profile
	GetAdminWithPermissions(ctx context.Context, userID uuid.UUID) (interface{}, error)
	GetBuyer(ctx context.Context, userID uuid.UUID) (interface{}, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, userType, currentPassword, newPassword string) error
}

type LoginResultSimplified struct {
	User        interface{} `json:"user"`
	AccessToken string      `json:"access_token"`
}

type authV2Service struct {
	authRepo     repositories.AuthRepository
	activityRepo repositories.ActivityLogRepository
	cfg          *config.Config
}

func NewAuthV2Service(
	authRepo repositories.AuthRepository,
	activityRepo repositories.ActivityLogRepository,
) AuthV2Service {
	return &authV2Service{
		authRepo:     authRepo,
		activityRepo: activityRepo,
		cfg:          config.LoadConfig(),
	}
}

func (s *authV2Service) AdminLogin(ctx context.Context, email, password string) (*LoginResultSimplified, error) {
	ipAddress := ""
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ipAddress = ginCtx.ClientIP()
	}

	// Find admin by email
	admin, err := s.authRepo.FindAdminByEmail(email)
	if err != nil {
		// Log failed login
		s.logActivity(ctx, nil, "ADMIN", models.ActionLoginFailed, "auth", "Login gagal: email tidak ditemukan", ipAddress)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(password, admin.Password) {
		// Log failed login
		s.logActivity(ctx, &admin.ID, "ADMIN", models.ActionLoginFailed, "auth", "Login gagal: password salah", ipAddress)
		return nil, errors.New("email atau password salah")
	}

	// Check if active
	if !admin.IsActive {
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi admin")
	}

	// Get role with permissions
	admin, err = s.authRepo.FindAdminWithRole(admin.ID)
	if err != nil {
		return nil, errors.New("gagal memuat data role")
	}

	// Extract permissions
	var permissions []string
	if admin.Role != nil {
		for _, perm := range admin.Role.Permissions {
			permissions = append(permissions, perm.Kode)
		}
	}

	// Generate token (24 hours)
	accessToken, err := utils.GenerateAccessToken(
		admin.ID,
		"ADMIN",
		admin.Email,
		admin.Role.ID.String(),
		admin.Role.Kode,
		permissions,
	)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.authRepo.UpdateAdminLastLogin(admin.ID)

	// Log successful login
	s.logActivity(ctx, &admin.ID, "ADMIN", models.ActionLogin, "auth", "Login berhasil", ipAddress)

	// Simplified response
	result := &LoginResultSimplified{
		User: map[string]interface{}{
			"id":    admin.ID.String(),
			"nama":  admin.Nama,
			"email": admin.Email,
		},
		AccessToken: accessToken,
	}

	return result, nil
}

func (s *authV2Service) BuyerLogin(ctx context.Context, email, password string) (*LoginResultSimplified, error) {
	ipAddress := ""
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ipAddress = ginCtx.ClientIP()
	}

	// Find buyer by email
	buyer, err := s.authRepo.FindBuyerByEmail(email)
	if err != nil {
		// Log failed login
		s.logActivity(ctx, nil, "BUYER", models.ActionLoginFailed, "auth", "Login gagal: email tidak ditemukan", ipAddress)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(password, buyer.Password) {
		// Log failed login
		s.logActivity(ctx, &buyer.ID, "BUYER", models.ActionLoginFailed, "auth", "Login gagal: password salah", ipAddress)
		return nil, errors.New("email atau password salah")
	}

	// Check if active
	if !buyer.IsActive {
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi admin")
	}

	// Generate token (24 hours)
	accessToken, err := utils.GenerateAccessToken(
		buyer.ID,
		"BUYER",
		buyer.Email,
		"",
		"",
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.authRepo.UpdateBuyerLastLogin(buyer.ID)

	// Log successful login
	s.logActivity(ctx, &buyer.ID, "BUYER", models.ActionLogin, "auth", "Login berhasil", ipAddress)

	// Simplified response
	result := &LoginResultSimplified{
		User: map[string]interface{}{
			"id":    buyer.ID.String(),
			"nama":  buyer.Nama,
			"email": buyer.Email,
		},
		AccessToken: accessToken,
	}

	return result, nil
}

func (s *authV2Service) GetAdminWithPermissions(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	admin, err := s.authRepo.FindAdminWithRole(userID)
	if err != nil {
		return nil, err
	}

	var permissions []string
	if admin.Role != nil {
		for _, perm := range admin.Role.Permissions {
			permissions = append(permissions, perm.Kode)
		}
	}

	return map[string]interface{}{
		"id":    admin.ID.String(),
		"nama":  admin.Nama,
		"email": admin.Email,
		"role": map[string]interface{}{
			"nama": admin.Role.Nama,
		},
		"permissions": permissions,
	}, nil
}

func (s *authV2Service) GetBuyer(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	buyer, err := s.authRepo.FindBuyerByID(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":       buyer.ID.String(),
		"nama":     buyer.Nama,
		"username": buyer.Username,
		"email":    buyer.Email,
		"telepon":  buyer.Telepon,
	}, nil
}

// Helper to log activity
func (s *authV2Service) logActivity(ctx context.Context, userID *uuid.UUID, userType string, action models.ActivityAction, modul, deskripsi, ipAddress string) {
	log := &models.ActivityLog{
		UserType:  userType,
		UserID:    userID,
		Action:    action,
		Modul:     modul,
		Deskripsi: deskripsi,
		IPAddress: &ipAddress,
		CreatedAt: time.Now(),
	}

	// Extract user agent from context if it's gin.Context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		userAgent := ginCtx.GetHeader("User-Agent")
		log.UserAgent = &userAgent
	}

	// Log asynchronously
	go s.activityRepo.Create(log)
}

// ChangePassword changes user password
func (s *authV2Service) ChangePassword(ctx context.Context, userID uuid.UUID, userType, currentPassword, newPassword string) error {
	ipAddress := ""
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ipAddress = ginCtx.ClientIP()
	}

	if userType == "ADMIN" {
		admin, err := s.authRepo.FindAdminByID(userID)
		if err != nil {
			return errors.New("admin tidak ditemukan")
		}

		// Verify current password
		if !utils.CheckPassword(currentPassword, admin.Password) {
			return errors.New("password saat ini salah")
		}

		// Check if new password is same as old
		if utils.CheckPassword(newPassword, admin.Password) {
			return errors.New("password baru tidak boleh sama dengan password lama")
		}

		// Hash new password with configured cost
		hashedPassword, err := utils.HashPasswordWithCost(newPassword, s.cfg.BcryptCost)
		if err != nil {
			return errors.New("gagal meng-hash password")
		}

		admin.Password = hashedPassword

		// Save to database
		if err := s.authRepo.UpdateAdmin(admin); err != nil {
			return err
		}

		// Log activity
		s.logActivity(ctx, &userID, "ADMIN", models.ActionUpdate, "security",
			"Mengubah password", ipAddress)

		return nil

	} else if userType == "BUYER" {
		buyer, err := s.authRepo.FindBuyerByID(userID)
		if err != nil {
			return errors.New("buyer tidak ditemukan")
		}

		// Verify current password
		if !utils.CheckPassword(currentPassword, buyer.Password) {
			return errors.New("password saat ini salah")
		}

		// Check if new password is same as old
		if utils.CheckPassword(newPassword, buyer.Password) {
			return errors.New("password baru tidak boleh sama dengan password lama")
		}

		// Hash new password with configured cost
		hashedPassword, err := utils.HashPasswordWithCost(newPassword, s.cfg.BcryptCost)
		if err != nil {
			return errors.New("gagal meng-hash password")
		}

		buyer.Password = hashedPassword

		// Save to database
		if err := s.authRepo.UpdateBuyer(buyer); err != nil {
			return err
		}

		// Log activity
		s.logActivity(ctx, &userID, "BUYER", models.ActionUpdate, "security",
			"Mengubah password", ipAddress)

		return nil
	}

	return errors.New("tipe user tidak valid")
}
