package models

import (
	"time"
)

type Transaction struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	IsNewCustomer bool      `json:"is_new_customer"`
	CustomerID    string    `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductSize   string    `json:"product_size"`
	Flavor        string    `json:"flavor"`
	Quantity      int       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
	TransactionAt time.Time `json:"transaction_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"update_at"`
}
