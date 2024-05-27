package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	DefaultCurrency *string `json:"default_currency"`
}
