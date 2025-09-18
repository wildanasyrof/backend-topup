package service

import (
	"strconv"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type UserService interface {
	GetUserByID(userID string) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepositry repository.UserRepository) UserService {
	return &userService{userRepository: userRepositry}
}

func (s *userService) GetUserByID(id string) (*entity.User, error) {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetByID(uint64(userId))
	if err != nil {
		return nil, err
	}

	return user, nil
}
