package service

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/hash"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterUserRequest) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginUserRequest) (*entity.User, string, error)
	RegisterByGoogle(ctx context.Context, userInfo *dto.GoogleLoginResponse) (*entity.User, error)
	GenerateToken(id uint64, role string) (*dto.TokenResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtService     jwt.JWTService
}

func NewAuthService(userRepository repository.UserRepository, jwtService jwt.JWTService) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

// RegisterByGoogle implements AuthService.
func (a *authService) RegisterByGoogle(ctx context.Context, userInfo *dto.GoogleLoginResponse) (*entity.User, error) {
	user := &entity.User{
		GoogleID:   userInfo.Sub,
		Name:       userInfo.Name,
		IsVerified: userInfo.EmailVerified,
	}

	if err := a.userRepository.Store(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login implements AuthService.
func (a *authService) Login(ctx context.Context, req *dto.LoginUserRequest) (*entity.User, string, error) {
	user, err := a.userRepository.GetByEmail(ctx, req.Email)

	if err != nil && user == nil {
		return nil, "", errors.New("invalid credential")
	}

	if err := hash.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", errors.New("invalid credential")
	}

	token, err := a.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token.AccessToken, nil
}

// Register implements AuthService.
func (a *authService) Register(ctx context.Context, req *dto.RegisterUserRequest) (*entity.User, error) {
	hashedPassword := hash.HashPassword(req.Password)

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := a.userRepository.Store(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authService) GenerateToken(id uint64, role string) (*dto.TokenResponse, error) {
	accessToken, err := a.jwtService.GenerateAccessToken(id, role)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken: accessToken,
	}, nil
}
