package dto

type DepositRequest struct {
	Amount          float64 `json:"amount" validate:"required,min=1000"`
	PaymentMethodID uint64  `json:"payment_method_id" validate:"required"`
	Fee             float64 `json:"fee,omitempty"`
}
