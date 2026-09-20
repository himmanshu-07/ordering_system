package models

type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Notes       string  `json:"notes,omitempty"`
}
