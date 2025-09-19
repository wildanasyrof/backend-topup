package entity

type Settings struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Value     string `gorm:"type:text;not null" json:"value"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}
