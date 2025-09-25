package entity

import "time"

type OrderStatus string

const (
	StatusSuccess    OrderStatus = "success"
	StatusCanceled   OrderStatus = "canceled"
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "processing"
)

// Standardize ID to uint64 for consistency
type Order struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderRef string `gorm:"size:50;not null;uniqueIndex:ux_orders_order_ref" json:"order_ref"`

	UserID uint64 `gorm:"not null;index:idx_orders_user_id" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`

	// Changed ProductID to uint64 to match Product.ID
	ProductID uint64   `gorm:"not null;index:idx_orders_product_id" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product,omitempty"`

	WA           string `gorm:"size:20" json:"wa"`
	Email        string `gorm:"size:100" json:"email"`
	CustomerName string `gorm:"size:100" json:"customer_name"`
	CustomerID   string `gorm:"size:100" json:"customer_id"`

	PaymentRef    string      `gorm:"size:100" json:"payment_ref"`
	PaymentStatus OrderStatus `json:"payment_status" gorm:"type:text;not null;default:processing;check:order_status_check,status IN ('success','canceled','pending','processing')"`
	Status        OrderStatus `json:"status" gorm:"type:text;not null;default:processing;check:order_status_check,status IN ('success','canceled','pending','processing')"`

	Amount float64 `gorm:"not null" json:"amount"`
	Fee    float64 `gorm:"not null;default:0" json:"fee"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Order) TableName() string { return "orders" }
