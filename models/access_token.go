package models

import "time"

type Access_Token struct {
	ID         uint `gorm:"primaryKey"`
	UserID   	 uint `gorm:"not null"`
	Revoked    bool `gorm:"default:false"`
	Expires_at time.Time
}