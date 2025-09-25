package dto

import "github.com/wildanasyrof/backend-topup/internal/domain/entity"

// In your DTO, keep json tags for JSON requests, and ADD form tags for multipart.
type CreatePaymentMethodRequest struct {
	Type       string   `json:"type" form:"type" validate:"required,min=3,max=100"`
	Name       string   `json:"name" form:"name" validate:"required,min=3,max=100"`
	ImgUrl     string   `json:"img_url" form:"img_url" validate:"omitempty,url"`
	ProviderID int64    `form:"provider_id" validate:"required"`
	Fee        *float64 `json:"fee,omitempty" form:"fee"`
	Percent    *float64 `json:"percent,omitempty" form:"percent"`
}

type UpdatePaymentMethodRequest struct {
	Type       string   `json:"type" form:"type" validate:"omitempty,min=3,max=100"`
	Name       string   `json:"name" form:"name" validate:"omitempty,min=3,max=100"`
	ImgUrl     string   `json:"img_url" form:"img_url" validate:"omitempty,startswith=/uploads/"`
	ProviderID *int64   `form:"provider_id"`
	Fee        *float64 `json:"fee,omitempty" form:"fee"`
	Percent    *float64 `json:"percent,omitempty" form:"percent"`
}

// in dto/update_payment_method_request.go
func (r *UpdatePaymentMethodRequest) ApplyTo(pm *entity.PaymentMethod) {
	if r.Name != "" {
		pm.Name = r.Name
	}
	if r.Type != "" {
		pm.Type = r.Type
	}
	if r.ImgUrl != "" {
		pm.ImgUrl = r.ImgUrl
	}
	if r.Fee != nil {
		pm.Fee = r.Fee
	}
	if r.Percent != nil {
		pm.Percent = r.Percent
	}
}
