package dto

import "github.com/wildanasyrof/backend-topup/internal/domain/entity"

type CreatePaymentMethodRequest struct {
	Type         string   `json:"type" validate:"required,min=3,max=100"`
	Name         string   `json:"name" validate:"required,min=3,max=100"`
	ImgUrl       string   `json:"img_url" validate:"omitempty,url"`
	Provider     string   `json:"provider" validate:"required"`
	ProviderCode string   `json:"provider_code" validate:"required"`
	Fee          *float64 `json:"fee,omitempty"`
	Percent      *float64 `json:"percent,omitempty"`
}

type UpdatePaymentMethodRequest struct {
	Type         string   `json:"type" validate:"omitempty,min=3,max=100"`
	Name         string   `json:"name" validate:"omitempty,min=3,max=100"`
	ImgUrl       string   `json:"img_url" validate:"omitempty,url"`
	Provider     string   `json:"provider" validate:"omitempty"`
	ProviderCode string   `json:"provider_code" validate:"omitempty"`
	Fee          *float64 `json:"fee,omitempty"`
	Percent      *float64 `json:"percent,omitempty"`
}

// in dto/update_payment_method_request.go
func (r *UpdatePaymentMethodRequest) ApplyTo(pm *entity.PaymentMethod) {
	if r.Name != "" {
		pm.Name = r.Name
	}
	if r.Type != "" {
		pm.Type = r.Type
	}
	if r.Provider != "" {
		pm.Provider = r.Provider
	}
	if r.ProviderCode != "" {
		pm.ProviderCode = r.ProviderCode
	}
	if r.Fee != nil {
		pm.Fee = r.Fee
	}
	if r.Percent != nil {
		pm.Percent = r.Percent
	}
}
