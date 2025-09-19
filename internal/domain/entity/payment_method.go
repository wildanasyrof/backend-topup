package entity

import "time"

type PaymentMethod struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Type         string    `json:"type" gorm:"type:varchar(255);not null"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	ImgUrl       string    `json:"img_url" gorm:"type:varchar(255);not null"`
	Provider     string    `json:"provider" gorm:"type:varchar(255);not null"`
	ProviderCode string    `json:"provider_code" gorm:"type:varchar(255);not null"`
	Fee          *float64  `json:"fee,omitempty"`
	Percent      *float64  `json:"percent,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime;"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime;"`
}

// TableName overrides the table name used by GORM
func (PaymentMethod) TableName() string { return "payment_methods" }
