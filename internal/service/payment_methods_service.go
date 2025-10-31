package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

type PaymentMethodsService interface {
	FindAll(ctx context.Context, q dto.PaymentMethodListQuery) ([]*entity.PaymentMethod, pagination.Meta, error)
	Create(ctx context.Context, req *dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error)
	Delete(ctx context.Context, id uint64) (*entity.PaymentMethod, error)
	Update(ctx context.Context, id uint64, req *dto.UpdatePaymentMethodRequest) (*entity.PaymentMethod, error)
	FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error)
}

type paymentMethodsService struct {
	repo repository.PaymentMethodsRepository
}

func NewPaymentMethodsService(repo repository.PaymentMethodsRepository) PaymentMethodsService {
	return &paymentMethodsService{repo: repo}
}

// Create implements PaymentMethodsService.
func (p *paymentMethodsService) Create(ctx context.Context, req *dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error) {
	data := &entity.PaymentMethod{
		Name:       req.Name,
		Type:       req.Type,
		ImgUrl:     req.ImgUrl,
		ProviderID: req.ProviderID,
		Fee:        req.Fee,
		Percent:    req.Percent,
	}

	if err := p.repo.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

// Delete implements PaymentMethodsService.
func (p *paymentMethodsService) Delete(ctx context.Context, id uint64) (*entity.PaymentMethod, error) {
	data, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(ctx, data.ID); err != nil {
		return nil, err
	}

	return data, nil
}

// FindAll implements PaymentMethodsService.
func (p *paymentMethodsService) FindAll(ctx context.Context, q dto.PaymentMethodListQuery) ([]*entity.PaymentMethod, pagination.Meta, error) {
	return p.repo.FindAll(ctx, q)
}

// FindByID implements PaymentMethodsService.
func (p *paymentMethodsService) FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error) {
	return p.repo.FindByID(ctx, id)
}

// Update implements PaymentMethodsService.
func (p *paymentMethodsService) Update(ctx context.Context, id uint64, req *dto.UpdatePaymentMethodRequest) (*entity.PaymentMethod, error) {
	data, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(data)

	if err := p.repo.Update(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}
