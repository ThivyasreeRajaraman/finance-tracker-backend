package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID               uint               `json:"user_id"`
	User                 User               `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	TransactionType      string             `json:"transaction_type"`
	CategoryID           *uint              `json:"category_id"`
	Category             Categories         `gorm:"foreignkey:CategoryID;association_foreignkey:CategoryID"`
	Amount               uint               `json:"amount"`
	TransactionPartnerID *uint              `json:"transaction_partner_id"`
	TransactionPartner   TransactionPartner `gorm:"foreignkey:TransactionPartnerID;association_foreignkey:TransactionPartnerID"`
}
