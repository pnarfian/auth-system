package models

import "time"

type ResetToken struct {
	ID         uint `gorm:"primaryKey"`
	Token   	 string `gorm:"not null"`
	UserID		 uint	`gorm:"not null"`
	CreatedAt  time.Time 
	ExpiresAt	 time.Time
	IsUsed		 bool `gorm:"default:false"`
}