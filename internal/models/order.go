package models

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusServed    OrderStatus = "served"
)


type Order struct {
	ID            string      `json:"id"`
	TableID       int         `json:"table_id"`
	Status        OrderStatus `json:"status"`
	Total         float64     `json:"total"`
	PaymentMethod string      `json:"payment_method"`
	CreatedAt     time.Time   `json:"created_at"`
	Items         []OrderItem `json:"items,omitempty"`
}
