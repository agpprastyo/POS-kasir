package utils

import (
	"POS-kasir/config"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"time"
)

// Manager defines the interface for JWT operations.
type Manager interface {
	GenerateToken(username, email string, userID uuid.UUID, role string) (string, time.Time, error)
	GenerateRefreshToken(username, email string, userID uuid.UUID, role string) (string, time.Time, error)
	VerifyToken(tokenStr string) (JWTClaims, error)
}

type JWTManager struct {
	cfg *config.AppConfig
}

func NewJWTManager(cfg *config.AppConfig) *JWTManager {
	return &JWTManager{
		cfg: cfg,
	}
}

type JWTClaims struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	UserID   uuid.UUID `json:"user_id"`
	Type     string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// ...
// GenerateToken creates a short-lived Access Token.
func (j *JWTManager) GenerateToken(username, email string, userID uuid.UUID, role string) (string, time.Time, error) {
	return j.generateToken(username, email, userID, role, "access", j.cfg.JWT.Duration)
}

// GenerateRefreshToken creates a long-lived Refresh Token.
func (j *JWTManager) GenerateRefreshToken(username, email string, userID uuid.UUID, role string) (string, time.Time, error) {
	return j.generateToken(username, email, userID, role, "refresh", j.cfg.JWT.RefreshTokenDuration)
}

func (j *JWTManager) generateToken(username, email string, userID uuid.UUID, role string, tokenType string, duration time.Duration) (string, time.Time, error) {
	claims := JWTClaims{
		Username: username,
		Email:    email,
		Role:     role,
		UserID:   userID,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
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
