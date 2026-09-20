package repository

import (
	"database/sql"

	"github.com/yourname/qr-ordering-system/internal/models"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(tableNumber int, qrCode string) (*models.Table, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO tables (table_number, qr_code) VALUES ($1, $2) RETURNING id`,
		tableNumber, qrCode,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.Table{ID: id, TableNumber: tableNumber, QRCode: qrCode}, nil
}

func (r *TableRepository) GetAll() ([]models.Table, error) {
	rows, err := r.db.Query(`SELECT id, table_number, qr_code FROM tables ORDER BY table_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var t models.Table
		if err := rows.Scan(&t.ID, &t.TableNumber, &t.QRCode); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func (r *TableRepository) GetByTableNumber(tableNumber int) (*models.Table, error) {
	var t models.Table
	err := r.db.QueryRow(
		`SELECT id, table_number, qr_code FROM tables WHERE table_number = $1`, tableNumber,
	).Scan(&t.ID, &t.TableNumber, &t.QRCode)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
