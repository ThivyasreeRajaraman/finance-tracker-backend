package helpers

type UpdateUserRequest struct {
	DefaultCurrency string `json:"default_currency"`
	Name            string `json:"name"`
}

type BudgetData struct {
	CategoryName *string `json:"category_name"`
	Amount       *uint   `json:"amount"`
	Threshold    *uint   `json:"threshold"`
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
}
