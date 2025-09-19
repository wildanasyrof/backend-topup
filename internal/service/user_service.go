package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type UserService interface {
	GetUserByID(userID uint64) (*entity.User, error)
	Update(userID uint64, req *dto.UpdateUserRequest) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepositry repository.UserRepository) UserService {
	return &userService{userRepository: userRepositry}
}

func (s *userService) GetUserByID(id uint64) (*entity.User, error) {
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// update implements UserService.
func (s *userService) Update(userID uint64, req *dto.UpdateUserRequest) (*entity.User, error) {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Whatsapp != nil {
		user.Whatsapp = *req.Whatsapp
	}

	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
