package models

import (
	"gorm.io/gorm"
)

type TransactionPartner struct {
	gorm.Model
	PartnerName    string `json:"partner_name"`
	UserID         uint   `json:"user_id"`
	User           User   `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	DueDate        string `json:"due_date"`
	ClosingBalance int    `json:"closing_balance"`
}
