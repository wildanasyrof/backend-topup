package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSession merepresentasikan refresh token yang valid untuk seorang user.
// ID-nya (UUID) adalah refresh token itu sendiri.
type UserSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	IsRevoked bool      `gorm:"not null;default:false" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	ClientIP  string    `gorm:"type:varchar(50)" json:"client_ip"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relasi
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}

// TableName menentukan nama tabel
func (UserSession) TableName() string {
	return "user_sessions"
}

// BeforeCreate akan men-generate UUID baru secara otomatis
func (s *UserSession) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
