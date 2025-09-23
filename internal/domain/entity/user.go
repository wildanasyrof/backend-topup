package entity

import (
	"time"
)

type User struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"                json:"id"`
	Name            string    `gorm:"not null"               json:"name"`
	Email           string    `gorm:"unique;not null"              json:"email"`
	Role            string    `gorm:"type:varchar(10);default:'user'"            json:"role"`
	Balance         float64   `gorm:";not null;default:0"         json:"balance"`
	EmailVerifiedAt time.Time `json:"verified_at"`
	PasswordHash    string    `gorm:"size:255"       json:"-"`
	Whatsapp        string    `gorm:"size:255"     json:"whatsapp"`
	GoogleID        string    `gorm:"size:255"       json:"google_id"`
	GoogleType      string    `gorm:"size:255"     json:"google_type"`
	OTP             int       `         json:"-"`
	IsVerified      bool      `gorm:"not null;default:false"      json:"-"`
	RememberToken   bool      `gorm:"not null;default:false"  json:"-"`
	UserLevelID     int       `json:"user_level_id" gorm:"not null;default:1"`
	CreatedAt       time.Time `gorm:"autoCreateTime"         json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"         json:"updated_at"`

	UserLevel UserLevel `json:"-" gorm:"foreignKey:UserLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (User) TableName() string { return "users" }
