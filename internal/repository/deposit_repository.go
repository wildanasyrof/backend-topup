package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type DepositRepository interface {
	Create(req *entity.Deposit) error
	FindByTopupID(topupID string) (*entity.Deposit, error)
	Update(deposit *entity.Deposit) error
	FindByUserID(userID uint64) ([]entity.Deposit, error)
	FindAll() ([]entity.Deposit, error)
}

type depositRepository struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &depositRepository{db: db}
}

// Create implements DepositRepository.
func (d *depositRepository) Create(req *entity.Deposit) error {
	return d.db.Create(req).Error
}

// FindAll implements DepositRepository.
func (d *depositRepository) FindAll() ([]entity.Deposit, error) {
	var deposits []entity.Deposit
	if err := d.db.Find(&deposits).Error; err != nil {
		return nil, err
	}

	return deposits, nil
}

// FindByTopupID implements DepositRepository.
func (d *depositRepository) FindByTopupID(topupID string) (*entity.Deposit, error) {
	var deposit entity.Deposit
	if err := d.db.Where("topup_id = ?", topupID).First(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

// FindByUserID implements DepositRepository.
func (d *depositRepository) FindByUserID(userID uint64) ([]entity.Deposit, error) {
	var deposits []entity.Deposit
	if err := d.db.Where("user_id = ?", userID).Find(&deposits).Error; err != nil {
		return nil, err
	}

	return deposits, nil
}

// Update implements DepositRepository.
func (d *depositRepository) Update(deposit *entity.Deposit) error {
	return d.db.Save(deposit).Error
}
