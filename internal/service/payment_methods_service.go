package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type PaymentMethodsService interface {
	FindAll() ([]*entity.PaymentMethod, error)
	Create(req *dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error)
	Delete(id uint64) (*entity.PaymentMethod, error)
	Update(id uint64, req *dto.UpdatePaymentMethodRequest) (*entity.PaymentMethod, error)
	FindByID(id uint64) (*entity.PaymentMethod, error)
}

type paymentMethodsService struct {
	repo repository.PaymentMethodsRepository
}

func NewPaymentMethodsService(repo repository.PaymentMethodsRepository) PaymentMethodsService {
	return &paymentMethodsService{repo: repo}
}

// Create implements PaymentMethodsService.
func (p *paymentMethodsService) Create(req *dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error) {
	data := &entity.PaymentMethod{
		Name:         req.Name,
		Type:         req.Type,
		ImgUrl:       req.ImgUrl,
		Provider:     req.Provider,
		ProviderCode: req.ProviderCode,
		Fee:          req.Fee,
		Percent:      req.Percent,
	}

	if err := p.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

// Delete implements PaymentMethodsService.
func (p *paymentMethodsService) Delete(id uint64) (*entity.PaymentMethod, error) {
	data, err := p.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := p.repo.Delete(data.ID); err != nil {
		return nil, err
	}

	return data, nil
}

// FindAll implements PaymentMethodsService.
func (p *paymentMethodsService) FindAll() ([]*entity.PaymentMethod, error) {
	return p.repo.FindAll()
}

// FindByID implements PaymentMethodsService.
func (p *paymentMethodsService) FindByID(id uint64) (*entity.PaymentMethod, error) {
	return p.repo.FindByID(id)
}

// Update implements PaymentMethodsService.
func (p *paymentMethodsService) Update(id uint64, req *dto.UpdatePaymentMethodRequest) (*entity.PaymentMethod, error) {
	data, err := p.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(data)

	if err := p.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}
