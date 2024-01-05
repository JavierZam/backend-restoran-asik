package models

type Transaction struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	ProductID       int     `json:"product_id"`
	Quantity        int     `json:"quantity"`
	TotalAmount     float64 `json:"total_amount"`
	TransactionDate string  `json:"transaction_date"`
	PaymentType     string  `json:"payment_type"`
	IsPaid          bool    `json:"is_paid"`
}
