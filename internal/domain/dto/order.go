package dto

type CreateOrder struct {
	ProductID    int    `json:"product_id" validate:"required"`
	WA           string `json:"wa,omitempty"`
	Email        string `json:"email,omitempty" validate:"required,email"`
	CustomerName string `json:"customer_name,omitempty" validate:"required"`
	CustomerID   string `json:"customer_id,omitempty" validate:"required"`
}

// UpdateOrder represents the payload to update an existing order
type UpdateOrder struct {
	WA           *string `json:"wa,omitempty"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	CustomerName *string `json:"customer_name,omitempty"`
	CustomerID   *string `json:"customer_id,omitempty"`
}

type OrderResponse struct {
	OrderRef      string  `json:"order_ref"`
	ProductID     uint64  `json:"product_id"`
	WA            string  `json:"wa"`
	Email         string  `json:"email"`
	CustomerName  string  `json:"customer_name"`
	CustomerID    string  `json:"customer_id"`
	PaymentRef    string  `json:"payment_ref"`
	PaymentStatus string  `json:"payment_status"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
}
