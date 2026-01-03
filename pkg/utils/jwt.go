package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("token tidak valid")
	ErrExpiredToken     = errors.New("token sudah kadaluarsa")
	ErrInvalidSignature = errors.New("signature token tidak valid")
)

var jwtSecret []byte
var accessTokenDuration time.Duration

// JWTClaims adalah custom claims untuk JWT
type JWTClaims struct {
	UserID      string   `json:"sub"`
	UserType    string   `json:"type"`
	Email       string   `json:"email"`
	RoleID      string   `json:"role_id,omitempty"`     // Hanya untuk Admin
	RoleKode    string   `json:"role_kode,omitempty"`   // Hanya untuk Admin
	Permissions []string `json:"permissions,omitempty"` // Hanya untuk Admin

	// Legacy support (akan dihapus)
	AdminID string `json:"admin_id,omitempty"`
	Role    string `json:"role,omitempty"` // Deprecated, use RoleKode
	jwt.RegisteredClaims
}

// SetJWTConfig sets JWT configuration (24 hours for single token)
func SetJWTConfig(secret string, accessDuration time.Duration) {
	jwtSecret = []byte(secret)
	accessTokenDuration = accessDuration
}

// SetJWTSecret sets JWT secret (legacy, deprecated)
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
	if accessTokenDuration == 0 {
		accessTokenDuration = 24 * time.Hour // Default 24 hours
	}
}

// GenerateAccessToken generates a new access token with user type and permissions
func GenerateAccessToken(userID uuid.UUID, userType, email, roleID, roleKode string, permissions []string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:      userID.String(),
		UserType:    userType,
		Email:       email,
		RoleID:      roleID,
		RoleKode:    roleKode,
		Permissions: permissions,
		Role:        roleKode, // Legacy
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessTokenDuration returns the access token duration in seconds
func GetAccessTokenDuration() int {
	return int(accessTokenDuration.Seconds())
}
