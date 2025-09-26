package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Store(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Destroy(ctx context.Context, id uint64) error
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
func (u *userRepository) Destroy(ctx context.Context, id uint64) error {
	return u.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID implements UserRepository.
func (u *userRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Store implements UserRepository.
func (u *userRepository) Store(ctx context.Context, user *entity.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// Update implements UserRepository.
func (u *userRepository) Update(ctx context.Context, user *entity.User) error {
	return u.db.WithContext(ctx).Save(user).Error
}
