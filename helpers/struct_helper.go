package helpers

type UpdateUserRequest struct {
	DefaultCurrency string `json:"default_currency"`
	Name            string `json:"name"`
}
