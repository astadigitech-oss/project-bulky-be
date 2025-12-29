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
var refreshTokenDuration time.Duration

// JWTClaims adalah custom claims untuk JWT
type JWTClaims struct {
	UserID      string   `json:"sub"`
	UserType    string   `json:"type"`
	Email       string   `json:"email"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	TokenID     string   `json:"jti,omitempty"` // For refresh token

	// Legacy support (akan dihapus)
	AdminID string `json:"admin_id,omitempty"`
	jwt.RegisteredClaims
}

// SetJWTConfig sets JWT configuration
func SetJWTConfig(secret string, accessDuration, refreshDuration time.Duration) {
	jwtSecret = []byte(secret)
	accessTokenDuration = accessDuration
	refreshTokenDuration = refreshDuration
}

// SetJWTSecret sets JWT secret (legacy, deprecated)
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
	if accessTokenDuration == 0 {
		accessTokenDuration = time.Hour
	}
	if refreshTokenDuration == 0 {
		refreshTokenDuration = 7 * 24 * time.Hour
	}
}

// GenerateAccessToken generates a new access token with user type and permissions
func GenerateAccessToken(userID uuid.UUID, userType, email, role string, permissions []string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:      userID.String(),
		UserType:    userType,
		Email:       email,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken(userID uuid.UUID, userType string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(refreshTokenDuration)
	tokenID := uuid.New()

	claims := JWTClaims{
		UserID:   userID.String(),
		UserType: userType,
		TokenID:  tokenID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, expiresAt, err
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

// GetAccessTokenDuration returns the access token duration
func GetAccessTokenDuration() int {
	return int(accessTokenDuration.Seconds())
}

// GetRefreshTokenDuration returns the refresh token duration
func GetRefreshTokenDuration() time.Duration {
	return refreshTokenDuration
}
