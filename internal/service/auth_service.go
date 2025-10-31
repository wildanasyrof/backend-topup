package service

import (
	"context"
	"time"

	"github.com/google/uuid" // <-- Import
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/hash"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterUserRequest) (*entity.User, error)
	// Login sekarang mengembalikan AccessToken dan Session
	Login(ctx context.Context, req *dto.LoginUserRequest, userAgent, clientIP string) (*entity.User, string, *entity.UserSession, error)
	// Refresh mengambil RT string (dari cookie)
	Refresh(ctx context.Context, oldRefreshToken string, userAgent, clientIP string) (string, *entity.UserSession, error)
	// Logout mengambil RT string (dari cookie)
	Logout(ctx context.Context, refreshToken string) error

	RegisterByGoogle(ctx context.Context, userInfo *dto.GoogleLoginResponse) (*entity.User, error)
	// CreateSession adalah helper baru, menggantikan GenerateToken
	CreateSession(ctx context.Context, user *entity.User, userAgent, clientIP string) (string, *entity.UserSession, error)
}

type authService struct {
	userRepository repository.UserRepository
	sessionRepo    repository.SessionRepository // <--- TAMBAHKAN
	jwtService     jwt.JWTService
}

// Modifikasi NewAuthService
func NewAuthService(
	userRepository repository.UserRepository,
	sessionRepo repository.SessionRepository, // <--- TAMBAHKAN
	jwtService jwt.JWTService,
) AuthService {
	return &authService{
		userRepository: userRepository,
		sessionRepo:    sessionRepo, // <--- TAMBAHKAN
		jwtService:     jwtService,
	}
}

// RegisterByGoogle implements AuthService.
func (a *authService) RegisterByGoogle(ctx context.Context, userInfo *dto.GoogleLoginResponse) (*entity.User, error) {
	user := &entity.User{
		GoogleID:   userInfo.Sub,
		Name:       userInfo.Name,
		IsVerified: userInfo.EmailVerified,
		Email:      userInfo.Email, // <-- Pastikan email disimpan
		Role:       "user",         // <-- Set role default
	}

	if err := a.userRepository.Store(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Modifikasi Login
func (a *authService) Login(ctx context.Context, req *dto.LoginUserRequest, userAgent, clientIP string) (*entity.User, string, *entity.UserSession, error) {
	user, err := a.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", nil, err
	}

	if err := hash.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", nil, apperror.New(apperror.CodeUnauthorized, "invalid credentials", err)
	}

	// Gunakan CreateSession
	accessToken, session, err := a.CreateSession(ctx, user, userAgent, clientIP)
	if err != nil {
		return nil, "", nil, apperror.New(apperror.CodeInternal, "issue token failed", err)
	}

	return user, accessToken, session, nil
}

// Register implements AuthService.
func (a *authService) Register(ctx context.Context, req *dto.RegisterUserRequest) (*entity.User, error) {
	hashedPassword := hash.HashPassword(req.Password)

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user", // <-- Set role default
	}

	if err := a.userRepository.Store(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Helper Baru: CreateSession
func (a *authService) CreateSession(ctx context.Context, user *entity.User, userAgent, clientIP string) (string, *entity.UserSession, error) {
	// 1. Buat Access Token
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", nil, apperror.New(apperror.CodeInternal, "failed to generate access token", err)
	}

	// 2. Buat Sesi (Refresh Token) di DB
	session := &entity.UserSession{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(a.jwtService.GetRefreshTokenDuration()),
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}

	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return "", nil, apperror.New(apperror.CodeInternal, "failed to create session", err)
	}

	// 3. Kembalikan AT string dan objek Session
	return accessToken, session, nil
}

// Fungsi Baru: Refresh
func (a *authService) Refresh(ctx context.Context, oldRefreshToken string, userAgent, clientIP string) (string, *entity.UserSession, error) {
	oldSessionID, err := uuid.Parse(oldRefreshToken)
	if err != nil {
		return "", nil, apperror.New(apperror.CodeUnauthorized, "invalid refresh token format", err)
	}

	// 1. Cari sesi lama (include User data)
	oldSession, err := a.sessionRepo.FindByID(ctx, oldSessionID)
	if err != nil {
		// Tidak ditemukan atau error DB
		return "", nil, err
	}

	// 2. Validasi sesi
	if oldSession.IsRevoked {
		_ = a.sessionRepo.Delete(ctx, oldSession.ID)
		return "", nil, apperror.New(apperror.CodeUnauthorized, "session has been revoked", nil)
	}
	if time.Now().After(oldSession.ExpiresAt) {
		// Hapus token kedaluwarsa dari DB
		_ = a.sessionRepo.Delete(ctx, oldSession.ID)
		return "", nil, apperror.New(apperror.CodeUnauthorized, "refresh token expired", nil)
	}

	// 3. Lakukan Rotasi: Hapus token lama
	// Kita gunakan Delete, bukan Revoke, karena kita akan ganti dengan yang baru
	if err := a.sessionRepo.Delete(ctx, oldSession.ID); err != nil {
		return "", nil, apperror.New(apperror.CodeInternal, "could not rotate token", err)
	}

	// 4. Buat token & sesi baru
	// Kita gunakan data User dari oldSession yg sudah di-Preload
	newAccessToken, newSession, err := a.CreateSession(ctx, &oldSession.User, userAgent, clientIP)
	if err != nil {
		return "", nil, apperror.New(apperror.CodeInternal, "could not issue new token", err)
	}

	return newAccessToken, newSession, nil
}

// Fungsi Baru: Logout
func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		// Tidak ada token untuk di-logout
		return nil
	}

	sessionID, err := uuid.Parse(refreshToken)
	if err != nil {
		// Token invalid, abaikan
		return nil
	}

	// Hapus sesi dari database
	// Kita bisa gunakan Revoke() jika ingin menyimpan history,
	// tapi Delete() lebih bersih untuk logout normal.
	return a.sessionRepo.Delete(ctx, sessionID)
}
