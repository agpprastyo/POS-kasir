package utils

import (
	"POS-kasir/config"
	"POS-kasir/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"time"
)

// Manager defines the interface for JWT operations.
type Manager interface {
	GenerateToken(username, email string, userID uuid.UUID, role repository.UserRole) (string, time.Time, error)
	VerifyToken(tokenStr string) (JWTClaims, error)
}

type JWTManager struct {
	cfg *config.AppConfig
}

type JWTClaims struct {
	Username string              `json:"username"`
	Email    string              `json:"email"`
	Role     repository.UserRole `json:"role"`
	UserID   uuid.UUID           `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWTManager.
func NewJWTManager(cfg *config.AppConfig) *JWTManager {
	return &JWTManager{
		cfg: cfg,
	}
}

// GenerateToken creates a JWT for the given username.
func (j *JWTManager) GenerateToken(username, email string, userID uuid.UUID, role repository.UserRole) (string, time.Time, error) {
	claims := JWTClaims{
		Username: username,
		Email:    email,
		Role:     role,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JWT.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.JWT.Issuer,
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, claims.ExpiresAt.Time, nil
}

// VerifyToken parses and validates a JWT.
func (j *JWTManager) VerifyToken(tokenStr string) (JWTClaims, error) {
	claims := JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.cfg.JWT.Secret), nil
	})
	if err != nil {
		return JWTClaims{}, err
	}
	if !token.Valid {
		return JWTClaims{}, errors.New("invalid token")
	}
	return claims, nil
}
