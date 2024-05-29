package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionPartner struct {
	gorm.Model
	PartnerName    string    `json:"partner_name"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	DueDate        time.Time `json:"due_date" gorm:"type:date"`
	ClosingBalance int       `json:"closing_balance"`
}
