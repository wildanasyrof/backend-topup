package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
)

type DepositService interface {
	Create(userID uint64, req *dto.DepositRequest) (*entity.Deposit, error)
	GetByUserID(userID uint64) ([]entity.Deposit, error)
	GetByDepositID(depositID string) (*entity.Deposit, error)
}

type depositService struct {
	repo repository.DepositRepository
}

func NewDepositService(repo repository.DepositRepository) DepositService {
	return &depositService{repo: repo}
}

// Create implements DepositService.
func (d *depositService) Create(userID uint64, req *dto.DepositRequest) (*entity.Deposit, error) {
	deposit := &entity.Deposit{
		UserID:          userID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		Status:          entity.DepProcessing,
		TopupID:         utils.GenerateTopupID(),
		Fee:             req.Fee,
	}

	if err := d.repo.Create(deposit); err != nil {
		return nil, err
	}

	return deposit, nil
}

// GetByDepositID implements DepositService.
func (d *depositService) GetByDepositID(depositID string) (*entity.Deposit, error) {
	return d.repo.FindByTopupID(depositID)
}

// GetByUserID implements DepositService.
func (d *depositService) GetByUserID(userID uint64) ([]entity.Deposit, error) {
	return d.repo.FindByUserID(userID)
}
