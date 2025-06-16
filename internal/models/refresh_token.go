package models

import (
	"time"
)

type RefreshToken struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"not null;unique"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Revoked   bool      `json:"revoked" gorm:"not null;default:false"`
}
