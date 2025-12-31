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
	// Authentication
	Login(ctx context.Context, email, password, userType, deviceInfo, ipAddress string) (*LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error

	// Profile
	GetCurrentUser(ctx context.Context, userID uuid.UUID, userType string) (interface{}, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, userType, nama string) (interface{}, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, userType, currentPassword, newPassword string) error
}

type LoginResult struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
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

func (s *authV2Service) Login(ctx context.Context, email, password, userType, deviceInfo, ipAddress string) (*LoginResult, error) {
	if userType == "ADMIN" {
		return s.loginAdmin(ctx, email, password, deviceInfo, ipAddress)
	} else if userType == "BUYER" {
		return s.loginBuyer(ctx, email, password, deviceInfo, ipAddress)
	}
	return nil, errors.New("tipe user tidak valid")
}

func (s *authV2Service) loginAdmin(ctx context.Context, email, password, deviceInfo, ipAddress string) (*LoginResult, error) {
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

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(
		admin.ID,
		"ADMIN",
		admin.Email,
		admin.Role.Kode,
		permissions,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, expiresAt, err := utils.GenerateRefreshToken(admin.ID, "ADMIN")
	if err != nil {
		return nil, err
	}

	// Hash and store refresh token
	hashedToken := utils.HashToken(refreshTokenStr)
	refreshToken := &models.RefreshToken{
		UserType:   models.UserTypeAdmin,
		UserID:     admin.ID,
		Token:      hashedToken,
		DeviceInfo: &deviceInfo,
		IPAddress:  &ipAddress,
		ExpiredAt:  expiresAt,
	}
	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// Update last login
	s.authRepo.UpdateAdminLastLogin(admin.ID)

	// Log successful login
	s.logActivity(ctx, &admin.ID, "ADMIN", models.ActionLogin, "auth", "Login berhasil", ipAddress)

	// Prepare response
	result := &LoginResult{
		User: map[string]interface{}{
			"id":    admin.ID,
			"nama":  admin.Nama,
			"email": admin.Email,
			"role": map[string]interface{}{
				"id":   admin.Role.ID,
				"nama": admin.Role.Nama,
				"kode": admin.Role.Kode,
			},
			"permissions": permissions,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    utils.GetAccessTokenDuration(),
	}

	return result, nil
}

func (s *authV2Service) loginBuyer(ctx context.Context, email, password, deviceInfo, ipAddress string) (*LoginResult, error) {
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

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(
		buyer.ID,
		"BUYER",
		buyer.Email,
		"",
		nil,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, expiresAt, err := utils.GenerateRefreshToken(buyer.ID, "BUYER")
	if err != nil {
		return nil, err
	}

	// Hash and store refresh token
	hashedToken := utils.HashToken(refreshTokenStr)
	refreshToken := &models.RefreshToken{
		UserType:   models.UserTypeBuyer,
		UserID:     buyer.ID,
		Token:      hashedToken,
		DeviceInfo: &deviceInfo,
		IPAddress:  &ipAddress,
		ExpiredAt:  expiresAt,
	}
	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// Update last login
	s.authRepo.UpdateBuyerLastLogin(buyer.ID)

	// Log successful login
	s.logActivity(ctx, &buyer.ID, "BUYER", models.ActionLogin, "auth", "Login berhasil", ipAddress)

	// Prepare response
	result := &LoginResult{
		User: map[string]interface{}{
			"id":      buyer.ID,
			"nama":    buyer.Nama,
			"email":   buyer.Email,
			"telepon": buyer.Telepon,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    utils.GetAccessTokenDuration(),
	}

	return result, nil
}

func (s *authV2Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate token format
	claims, err := utils.ValidateJWT(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	// Check if token exists in database
	hashedToken := utils.HashToken(refreshToken)
	storedToken, err := s.authRepo.FindRefreshToken(hashedToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau sudah expired")
	}

	// Get user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("user ID tidak valid")
	}

	var newAccessToken string
	var permissions []string

	// Generate new access token based on user type
	if storedToken.UserType == models.UserTypeAdmin {
		admin, err := s.authRepo.FindAdminWithRole(userID)
		if err != nil {
			return nil, errors.New("user tidak ditemukan")
		}

		// Extract permissions
		if admin.Role != nil {
			for _, perm := range admin.Role.Permissions {
				permissions = append(permissions, perm.Kode)
			}
		}

		newAccessToken, err = utils.GenerateAccessToken(
			admin.ID,
			"ADMIN",
			admin.Email,
			admin.Role.Kode,
			permissions,
		)
		if err != nil {
			return nil, err
		}
	} else if storedToken.UserType == models.UserTypeBuyer {
		buyer, err := s.authRepo.FindBuyerByID(userID)
		if err != nil {
			return nil, errors.New("user tidak ditemukan")
		}

		newAccessToken, err = utils.GenerateAccessToken(
			buyer.ID,
			"BUYER",
			buyer.Email,
			"",
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	// Optionally rotate refresh token (for better security)
	// For now, we'll keep the same refresh token

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    utils.GetAccessTokenDuration(),
	}, nil
}

func (s *authV2Service) Logout(ctx context.Context, refreshToken string) error {
	hashedToken := utils.HashToken(refreshToken)
	return s.authRepo.RevokeRefreshToken(hashedToken)
}

func (s *authV2Service) GetCurrentUser(ctx context.Context, userID uuid.UUID, userType string) (interface{}, error) {
	if userType == "ADMIN" {
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
			"id":    admin.ID,
			"nama":  admin.Nama,
			"email": admin.Email,
			"role": map[string]interface{}{
				"id":   admin.Role.ID,
				"nama": admin.Role.Nama,
				"kode": admin.Role.Kode,
			},
			"permissions": permissions,
		}, nil
	} else if userType == "BUYER" {
		buyer, err := s.authRepo.FindBuyerByID(userID)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"id":      buyer.ID,
			"nama":    buyer.Nama,
			"email":   buyer.Email,
			"telepon": buyer.Telepon,
		}, nil
	}

	return nil, errors.New("tipe user tidak valid")
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

// UpdateProfile updates user profile (nama)
func (s *authV2Service) UpdateProfile(ctx context.Context, userID uuid.UUID, userType, nama string) (interface{}, error) {
	if userType == "ADMIN" {
		admin, err := s.authRepo.FindAdminByID(userID)
		if err != nil {
			return nil, errors.New("admin tidak ditemukan")
		}

		// Update nama
		oldNama := admin.Nama
		admin.Nama = nama

		// Save to database (need to add Update method to repository)
		if err := s.authRepo.UpdateAdmin(admin); err != nil {
			return nil, err
		}

		// Log activity
		ipAddress := ""
		if ginCtx, ok := ctx.(*gin.Context); ok {
			ipAddress = ginCtx.ClientIP()
		}
		s.logActivity(ctx, &userID, "ADMIN", models.ActionUpdate, "profile",
			"Mengupdate profile dari '"+oldNama+"' ke '"+nama+"'", ipAddress)

		// Get full data with role
		adminWithRole, _ := s.authRepo.FindAdminWithRole(userID)
		if adminWithRole != nil {
			admin = adminWithRole
		}

		permissions := []string{}
		if admin.Role != nil {
			for _, p := range admin.Role.Permissions {
				permissions = append(permissions, p.Kode)
			}
		}

		return map[string]interface{}{
			"id":    admin.ID,
			"nama":  admin.Nama,
			"email": admin.Email,
			"role": map[string]interface{}{
				"id":   admin.Role.ID,
				"nama": admin.Role.Nama,
				"kode": admin.Role.Kode,
			},
			"permissions": permissions,
		}, nil

	} else if userType == "BUYER" {
		buyer, err := s.authRepo.FindBuyerByID(userID)
		if err != nil {
			return nil, errors.New("buyer tidak ditemukan")
		}

		// Update nama
		oldNama := buyer.Nama
		buyer.Nama = nama

		// Save to database
		if err := s.authRepo.UpdateBuyer(buyer); err != nil {
			return nil, err
		}

		// Log activity
		ipAddress := ""
		if ginCtx, ok := ctx.(*gin.Context); ok {
			ipAddress = ginCtx.ClientIP()
		}
		s.logActivity(ctx, &userID, "BUYER", models.ActionUpdate, "profile",
			"Mengupdate profile dari '"+oldNama+"' ke '"+nama+"'", ipAddress)

		return map[string]interface{}{
			"id":      buyer.ID,
			"nama":    buyer.Nama,
			"email":   buyer.Email,
			"telepon": buyer.Telepon,
		}, nil
	}

	return nil, errors.New("tipe user tidak valid")
}

// ChangePassword changes user password
func (s *authV2Service) ChangePassword(ctx context.Context, userID uuid.UUID, userType, currentPassword, newPassword string) error {
	if userType == "ADMIN" {
		admin, err := s.authRepo.FindAdminByID(userID)
		if err != nil {
			return errors.New("admin tidak ditemukan")
		}

		// Verify current password
		if !utils.CheckPassword(currentPassword, admin.Password) {
			return errors.New("password lama tidak sesuai")
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
		ipAddress := ""
		if ginCtx, ok := ctx.(*gin.Context); ok {
			ipAddress = ginCtx.ClientIP()
		}
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
			return errors.New("password lama tidak sesuai")
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
		ipAddress := ""
		if ginCtx, ok := ctx.(*gin.Context); ok {
			ipAddress = ginCtx.ClientIP()
		}
		s.logActivity(ctx, &userID, "BUYER", models.ActionUpdate, "security",
			"Mengubah password", ipAddress)

		return nil
	}

	return errors.New("tipe user tidak valid")
}
