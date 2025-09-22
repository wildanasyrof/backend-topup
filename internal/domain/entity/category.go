package entity

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
	Type        Type      `json:"type" gorm:"type:text;not null;default=prabayar;check:cat_status_check,status IN ('prabayar','pascabayar')"`
	MenuID      int64     `json:"menu_id" gorm:"not null;column:menu_id"`
	ProviderID  int64     `json:"provider_id" gorm:"primaryKey;autoIncrement"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);not null;unique"`
	Status      CatStatus `json:"status" gorm:"type:text;not null;default=inactive;check:cat_status_check,status IN ('inactive','active','problem')"`
	Description string    `json:"description" gorm:"column:description"`
	InputType   string    `json:"input_type" gorm:"type:varchar(255);not null"`
	ImgUrl      string    `json:"img_url" gorm:"not null;url"`
	IsLogin     bool      `json:"is_login"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime"`

	Menu     `json:"-" gorm:"foreignKey:MenuID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Provider `json:"-" gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName overrides the table name used by GORM
func (Category) TableName() string { return "categories" }
