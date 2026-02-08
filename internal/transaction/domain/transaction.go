package domain

import (
	"time"
)

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type SalesSummary struct {
	TotalRevenue      int             `json:"total_revenue"`
	TotalTransactions int             `json:"total_transactions"`
	MostSoldProduct   MostSoldProduct `json:"most_sold_product"`
}

type MostSoldProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}
