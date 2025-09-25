package entity

import "time"

type Menu struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Categories []Category `gorm:"foreignKey:MenuID" json:"categories"` // <- important
}

// TableName overrides the table name used by GORM
func (Menu) TableName() string { return "menus" }
