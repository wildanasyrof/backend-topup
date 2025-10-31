package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wildanasyrof/backend-topup/internal/config"
)

type JWTService interface {
	GenerateAccessToken(id uint64, role string) (string, error)
	// GenerateRefreshToken(userID uuid.UUID) (string, error)
	GetRefreshTokenDuration() time.Duration
	ValidateToken(token string) (uint64, string, error)
}

type jwtService struct {
	SecretKey  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		SecretKey:  cfg.JWT.AccessSecret,
		accessTTL:  time.Duration(cfg.JWT.AccessTokenMinutes) * time.Minute,
		refreshTTL: time.Duration(cfg.JWT.RefreshTokenDays) * time.Hour * 24, // Convert days to hours
	}
}

type AccessTokenClaims struct {
	Id   uint64 `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken implements JWTService.
func (j *jwtService) GenerateAccessToken(id uint64, role string) (string, error) {
	claims := AccessTokenClaims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// // GenerateRefreshToken implements JWTService.
// func (j *jwtService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
// 	claims := RefreshTokenClaims{
// 		UserID: userID.String(),
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(j.SecretKey))
// }

// ValidateToken implements JWTService.
// Returns userID if valid, otherwise error.
func (j *jwtService) ValidateToken(tokenStr string) (uint64, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return 0, "", err
	}
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims.Id, claims.Role, nil
	}
	return 0, "", errors.New("invalid token")
}

func (j *jwtService) GetRefreshTokenDuration() time.Duration {
	return j.refreshTTL
}
