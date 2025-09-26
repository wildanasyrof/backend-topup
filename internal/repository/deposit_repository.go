package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type DepositRepository interface {
	Create(ctx context.Context, req *entity.Deposit) error
	FindByTopupID(ctx context.Context, topupID string) (*entity.Deposit, error)
	Update(ctx context.Context, deposit *entity.Deposit) error
	FindByUserID(ctx context.Context, userID uint64) ([]entity.Deposit, error)
	FindAll(ctx context.Context) ([]entity.Deposit, error)
}

type depositRepository struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &depositRepository{db: db}
}

// Create implements DepositRepository.
func (d *depositRepository) Create(ctx context.Context, req *entity.Deposit) error {
	return d.db.WithContext(ctx).Create(req).Error
}

// FindAll implements DepositRepository.
func (d *depositRepository) FindAll(ctx context.Context) ([]entity.Deposit, error) {
	var deposits []entity.Deposit
	if err := d.db.WithContext(ctx).Find(&deposits).Error; err != nil {
		return nil, err
	}

	return deposits, nil
}

// FindByTopupID implements DepositRepository.
func (d *depositRepository) FindByTopupID(ctx context.Context, topupID string) (*entity.Deposit, error) {
	var deposit entity.Deposit
	if err := d.db.WithContext(ctx).Where("topup_id = ?", topupID).First(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

// FindByUserID implements DepositRepository.
func (d *depositRepository) FindByUserID(ctx context.Context, userID uint64) ([]entity.Deposit, error) {
	var deposits []entity.Deposit
	if err := d.db.WithContext(ctx).Where("user_id = ?", userID).Find(&deposits).Error; err != nil {
		return nil, err
	}

	return deposits, nil
}

// Update implements DepositRepository.
func (d *depositRepository) Update(ctx context.Context, deposit *entity.Deposit) error {
	return d.db.WithContext(ctx).Save(deposit).Error
}
