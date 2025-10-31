package entity

import "time"

type Price struct {
	ID          int `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID   int `json:"product_id" gorm:"not null"`
	UserLevelID int `json:"user_level_id" gorm:"not null;default:1"`

	// --- UBAH JSON TAG DARI "Price" -> "price" ---
	Price float64 `json:"price" gorm:"not null;column:amount"`
	// ---

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// --- UBAH JSON TAG DARI "-" -> "user_level,omitempty" ---
	UserLevel *UserLevel `json:"user_level,omitempty" gorm:"foreignKey:UserLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	// ---

	// --- TAMBAHKAN RELASI KE PRODUK ---
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	// ---
}

func (Price) TableName() string { return "prices" }
