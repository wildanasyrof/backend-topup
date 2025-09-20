package entity

import "time"

type Banner struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ImgUrl    string    `json:"img_url" gorm:"not null;url"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
