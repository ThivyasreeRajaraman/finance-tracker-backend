package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	Income  TransactionType = "Income"
	Expense TransactionType = "Expense"
	Lend    TransactionType = "Lend"
	Borrow  TransactionType = "Borrow"
)

type Frequencies string

const (
	Daily   Frequencies = "Daily"
	Weekly  Frequencies = "Weekly"
	Monthly Frequencies = "Monthly"
)

type Transaction struct {
	gorm.Model
	UserID             uint            `json:"user_id"`
	User               User            `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	TransactionType    TransactionType `json:"transaction_type_id"`
	CategoryID         *uint           `json:"category_id"`
	Category           Categories      `gorm:"foreignkey:CategoryID;association_foreignkey:CategoryID"`
	Amount             int             `json:"amount"`
	TransactionPartner string          `json:"transaction_partner"`
	Frequency          Frequencies     `json:"frequency"`
	PaymentDueDate     *time.Time      `json:"due_date"`
}
