package entity

type Type string

const (
	TypePrabayar   Type = "prabayar"
	TypePascabayar Type = "pascabayar"
)

type Category struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:varchar(255);not null;unique"`
	Type Type   `json:"type" gorm:"type:enum('prabayar','pascabayar');not null"`
	Menu []Menu `json:"menu,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Slug string `json:"slug" gorm:"type:varchar(255);not null;unique"`

	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}
