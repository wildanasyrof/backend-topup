package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wildanasyrof/backend-topup/internal/config"
)

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	GetRefreshTokenDuration() time.Duration
	ValidateToken(token string) (string, error)
}

type jwtService struct {
	SecretKey  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		SecretKey:  cfg.JWT.JWTSecret,
		accessTTL:  time.Duration(cfg.JWT.AccessTokenMinutes) * time.Minute,
		refreshTTL: time.Duration(cfg.JWT.RefreshTokenDays) * time.Hour * 24, // Convert days to hours
	}
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken implements JWTService.
func (j *jwtService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := AccessTokenClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// GenerateRefreshToken implements JWTService.
func (j *jwtService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := RefreshTokenClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// ValidateToken implements JWTService.
// Returns userID if valid, otherwise error.
func (j *jwtService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}
	return "", errors.New("invalid token")
}

func (j *jwtService) GetRefreshTokenDuration() time.Duration {
	return j.refreshTTL
}
