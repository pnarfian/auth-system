package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"type:varchar(100);not null"`
	FirstName string `gorm:"type:varchar(350);not null"`
	LastName  string `gorm:"type:varchar(500);not null"`
	TelephoneNo string `gorm:"type:varchar(15);not null"`
	Email			string	`gorm:"type:varchar(350);not null"`
	Password	string	`gorm:"type:varchar(200);not null"`
	CreatedAt time.Time
	UpdateAt 	time.Time
	IsDeleted	bool	`gorm:"default:false"`
}