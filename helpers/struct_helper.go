package helpers

import "time"

type UpdateUserRequest struct {
	DefaultCurrency string `json:"default_currency"`
	Name            string `json:"name"`
}

type BudgetData struct {
	CategoryName *string `json:"category_name"`
	Amount       *int    `json:"amount"`
	Threshold    *int    `json:"threshold"`
}

type BudgetResponse struct {
	ID              uint   `json:"id"`
	UserID          uint   `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	DefaultCurrency string `json:"default_currency"`
	CategoryID      uint   `json:"category_id"`
	CategoryName    string `json:"category_name"`
	Amount          int    `json:"amount"`
	Threshold       int    `json:"threshold"`
}

type TransactionData struct {
	TransactionType    string     `json:"transaction_type"`
	CategoryName       *string    `json:"category_name,omitempty"`
	Amount             uint       `json:"amount"`
	TransactionPartner *string    `json:"transaction_partner,omitempty"`
	PaymentDueDate     *time.Time `json:"payment_due_date,omitempty" gorm:"type:date"`
}

type TransactionResponse struct {
	ID                 uint       `json:"id"`
	UserID             uint       `json:"user_id"`
	Name               string     `json:"name"`
	DefaultCurrency    string     `json:"default_currency"`
	TransactionType    string     `json:"transaction_type"`
	CategoryID         *uint      `json:"category_id,omitempty"`
	CategoryName       *string    `json:"category_name,omitempty"`
	Amount             uint       `json:"amount"`
	TransactionPartner *string    `json:"transaction_partner,omitempty"`
	PaymentDueDate     *time.Time `json:"payment_due_date,omitempty" gorm:"type:date"`
}

type TransactionUpdate struct {
	TransactionType *string `json:"transaction_type,omitempty"`
	CategoryName    *string `json:"category_name,omitempty"`
	Amount          *uint   `json:"amount,omitempty"`
}

type TransactionPartnerData struct {
	PartnerName string `json:"partner_name"`
}

type TransactionPartnerResponse struct {
	PartnerName     string    `json:"partner_name"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	DueDate         time.Time `json:"payment_due_date" gorm:"type:date"`
}
