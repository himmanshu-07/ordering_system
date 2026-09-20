package models

type Table struct {
	ID          int    `json:"id"`
	TableNumber int    `json:"table_number"`
	QRCode      string `json:"qr_code"`
}
