package service

import (
	"context" // Import the context package

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

// UserService interface updated to include context.Context
type UserService interface {
	GetUserByID(ctx context.Context, userID uint64) (*entity.User, error)
	FindUserByGoogleID(ctx context.Context, id string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, userID uint64, req *dto.UpdateUserRequest) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepositry repository.UserRepository) UserService {
	return &userService{userRepository: userRepositry}
}

// FindUserByEmail implements UserService.
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*entity.User, error) {
	// Pass ctx to the repository call
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// update implements UserService.
func (s *userService) Update(ctx context.Context, userID uint64, req *dto.UpdateUserRequest) (*entity.User, error) {
	// Pass ctx to the repository call
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Whatsapp != nil {
		user.Whatsapp = *req.Whatsapp
	}

	// Pass ctx to the repository call
	if err := s.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByGoogleID implements AuthService.
func (a *userService) FindUserByGoogleID(ctx context.Context, id string) (*entity.User, error) {
	user, err := a.userRepository.FindByGoogleID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, err
}
