package entity

import "time"

type Price struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID   int       `json:"product_id" gorm:"not null"`
	UserLevelID int       `json:"user_level_id" gorm:"not null;default:1"`
	Price       float64   `gorm:"not null;column:amount"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	UserLevel *UserLevel `json:"-" gorm:"foreignKey:UserLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Price) TableName() string { return "prices" }
