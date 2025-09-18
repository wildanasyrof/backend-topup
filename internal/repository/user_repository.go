package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Store(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id uint64) (*entity.User, error)
	Update(user *entity.User) error
	Destroy(id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Destroy implements UserRepository.
func (u *userRepository) Destroy(id uint64) error {
	return u.db.Delete(&entity.User{}, id).Error
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID implements UserRepository.
func (u *userRepository) GetByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Store implements UserRepository.
func (u *userRepository) Store(user *entity.User) error {
	return u.db.Create(user).Error
}

// Update implements UserRepository.
func (u *userRepository) Update(user *entity.User) error {
	return u.db.Save(user).Error
}
