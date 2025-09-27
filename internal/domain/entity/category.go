package entity

import "time"

type Type string
type CatStatus string

const (
	TypePrabayar   Type = "prabayar"
	TypePascabayar Type = "pascabayar"
)

const (
	CatActive   CatStatus = "active"
	CatInactive CatStatus = "inactive"
	CatProblem  CatStatus = "problem"
)

type Category struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	Type        Type      `json:"type" gorm:"type:text;not null;default:prabayar;check:cat_type_check,type IN ('prabayar','pascabayar')"`
	MenuID      int64     `json:"menu_id" gorm:"not null"`
	ProviderID  int64     `json:"provider_id" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Status      CatStatus `json:"status" gorm:"type:text;not null;default:inactive;check:cat_status_check,status IN ('inactive','active','problem')"`
	Description string    `json:"description"`
	InputType   string    `json:"input_type" gorm:"type:varchar(255);not null"`
	ImgUrl      string    `json:"img_url" gorm:"not null"`
	IsLogin     bool      `json:"is_login"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Menu     Menu     `json:"-" gorm:"foreignKey:MenuID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Provider Provider `json:"-" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products"` // <- important
}

// TableName overrides the table name used by GORM
func (Category) TableName() string { return "categories" }
