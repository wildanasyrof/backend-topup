package entity

import "time"

type Product struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	SkuCode     string    `json:"sku_code" gorm:"not null;unique"`
	SellerName  string    `json:"seller_name" gorm:"varchar(255)"`
	CategoryID  int       `json:"category_id" gorm:"not null"`
	ProviderID  int64     `json:"provider_id" gorm:"not null"`
	Status      CatStatus `json:"status" gorm:"type:text;not null;default:inactive;check:cat_status_check,status IN ('inactive','active','problem')"`
	Stock       int64     `json:"stock" gorm:"not null"`
	Description string    `json:"description"`
	ImgUrl      string    `json:"img_url" gorm:"not null"`
	StartOff    string    `json:"start_off"`
	EndOff      string    `json:"end_off"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Category Category `json:"-" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Provider Provider `json:"-" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Prices []Price `gorm:"foreignKey:ProductID" json:"prices"` // <- important
}

func (Product) TableName() string { return "products" }
