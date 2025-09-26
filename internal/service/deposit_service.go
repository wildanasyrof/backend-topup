package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/utils"
)

type DepositService interface {
	Create(ctx context.Context, userID uint64, req *dto.DepositRequest) (*entity.Deposit, error)
	GetByUserID(ctx context.Context, userID uint64) ([]entity.Deposit, error)
	GetByDepositID(ctx context.Context, depositID string) (*entity.Deposit, error)
}

type depositService struct {
	repo repository.DepositRepository
}

func NewDepositService(repo repository.DepositRepository) DepositService {
	return &depositService{repo: repo}
}

// Create implements DepositService.
func (d *depositService) Create(ctx context.Context, userID uint64, req *dto.DepositRequest) (*entity.Deposit, error) {
	deposit := &entity.Deposit{
		UserID:          userID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		Status:          entity.DepProcessing,
		TopupID:         utils.GenerateTopupID(),
		Fee:             req.Fee,
	}

	if err := d.repo.Create(ctx, deposit); err != nil {
		return nil, err
	}

	return deposit, nil
}

// GetByDepositID implements DepositService.
func (d *depositService) GetByDepositID(ctx context.Context, depositID string) (*entity.Deposit, error) {
	return d.repo.FindByTopupID(ctx, depositID)
}

// GetByUserID implements DepositService.
func (d *depositService) GetByUserID(ctx context.Context, userID uint64) ([]entity.Deposit, error) {
	return d.repo.FindByUserID(ctx, userID)
}
