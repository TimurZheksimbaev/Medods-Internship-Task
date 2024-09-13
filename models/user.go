package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Email    string
	Password string
}

type RefreshToken struct {
	gorm.Model
	UserID string
	Token  string
}