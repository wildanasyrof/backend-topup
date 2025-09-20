package entity

import "time"

// DepositStatus is a custom type to represent the status of a deposit.
type DepositStatus string

const (
	DepPending    DepositStatus = "pending"
	DepProcessing DepositStatus = "processing"
	DepSuccess    DepositStatus = "success"
	DepCanceled   DepositStatus = "canceled"
)

// Deposit represents the `deposit` table in the database.
type Deposit struct {
	ID              uint64        `gorm:"primaryKey;autoIncrement;column:id"`
	TopupID         string        `gorm:"type:varchar(255);not null;column:topup_id"`
	UserID          uint64        `gorm:"not null;column:user_id"`
	PaymentMethodID uint64        `gorm:"not null;column:payment_method_id"`
	Amount          float64       `gorm:"not null;column:amount"`
	Status          DepositStatus `gorm:"type:text;not null;default=pending;check:deposit_status_check,status IN ('pending','processing','success','canceled')"`
	Fee             float64       `gorm:"not null;column:fee"`
	Payment         string        `gorm:"type:text;not null;column:payment"`
	CreatedAt       time.Time     `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime;column:updated_at"`

	User   User          `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Method PaymentMethod `json:"-" gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM.
func (Deposit) TableName() string {
	return "deposit"
}
