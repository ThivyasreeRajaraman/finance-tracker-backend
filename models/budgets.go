package models

import "gorm.io/gorm"

type Budgets struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	User       User       `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	CategoryID uint       `json:"category_id"`
	Category   Categories `gorm:"foreignkey:CategoryID;association_foreignkey:CategoryID"`
	Amount     int        `json:"amount"`
	Threshold  int        `json:"threshold"`
	Active     bool       `json:"active"`
}
