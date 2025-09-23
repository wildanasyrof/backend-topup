package entity

type Product struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	CategoryID  int       `json:"category_id" gorm:"not null"`
	ProviderID  int64     `json:"provider_id" gorm:"not null"`
	Status      CatStatus `json:"status" gorm:"type:text;not null;default:inactive;check:cat_status_check,status IN ('inactive','active','problem')"`
	Description string    `json:"description"`
	ImgUrl      string    `json:"img_url" gorm:"not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime"`

	Category Category `json:"-" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Provider Provider `json:"-" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Prices []Price `gorm:"foreignKey:ProductID" json:"prices"` // <- important
}

func (Product) TableName() string { return "products" }
