package models

import "gorm.io/gorm"

type Categories struct {
	gorm.Model
	Name   string `json:"name"`
	Type   string `json:"type"`
	UserID *uint  `json:"user_id"`
	User   User   `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	Active bool   `json:"active"`
}
