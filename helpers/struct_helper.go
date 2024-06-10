package helpers

type UpdateUserRequest struct {
	DefaultCurrency string `json:"default_currency"`
	Name            string `json:"name"`
}

type BudgetData struct {
	CategoryName *string `json:"category_name"`
	Amount       *uint   `json:"amount"`
	Threshold    *uint   `json:"threshold"`
	Currency     string  `json:"currency"`
}

type BudgetResponse struct {
	ID              uint   `json:"id"`
	UserID          uint   `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	DefaultCurrency string `json:"default_currency"`
	CategoryID      uint   `json:"category_id"`
	CategoryName    string `json:"category_name"`
	Amount          uint   `json:"amount"`
	Threshold       uint   `json:"threshold"`
	Currency        string `json:"currency"`
}

type TransactionData struct {
	TransactionType    string  `json:"transaction_type"`
	CategoryName       *string `json:"category_name,omitempty"`
	Amount             uint    `json:"amount"`
	TransactionPartner *string `json:"transaction_partner,omitempty"`
	PaymentDueDate     *string `json:"payment_due_date,omitempty"`
}

type TransactionResponse struct {
	ID                 uint    `json:"id"`
	UserID             uint    `json:"user_id"`
	Name               string  `json:"name"`
	DefaultCurrency    string  `json:"default_currency"`
	TransactionType    string  `json:"transaction_type"`
	CategoryID         *uint   `json:"category_id,omitempty"`
	CategoryName       *string `json:"category_name,omitempty"`
	Amount             uint    `json:"amount"`
	TransactionPartner *string `json:"transaction_partner,omitempty"`
	PaymentDueDate     *string `json:"payment_due_date,omitempty"`
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
	PartnerName     string `json:"partner_name"`
	Amount          uint   `json:"amount"`
	TransactionType string `json:"transaction_type"`
	DueDate         string `json:"payment_due_date"`
}

type RecurringExpenseData struct {
	CategoryName    string `json:"category_name"`
	Amount          uint   `json:"amount"`
	Frequency       string `json:"frequency"`
	NextExpenseDate string `json:"next_expense_date"`
	Currency        string `json:"currency"`
}

type UpdateRecurringExpenseData struct {
	CategoryName    *string `json:"category_name"`
	Amount          *uint   `json:"amount"`
	Frequency       *string `json:"frequency"`
	NextExpenseDate *string `json:"next_expense_date"`
	Currency        *string `json:"currency"`
	Active          *bool   `json:"active"`
}

type RecurringExpenseResponse struct {
	Category        string `json:"category"`
	Amount          uint   `json:"amount"`
	Frequency       string `json:"frequency"`
	NextExpenseDate string `json:"next_expense_date"`
	Currency        string `json:"currency"`
	Active          bool   `json:"active"`
}
