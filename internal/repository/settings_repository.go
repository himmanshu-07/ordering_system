package repository

import (
	"database/sql"
	"strconv"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetTotalTables() (int, error) {
	var value string
	err := r.db.QueryRow(`SELECT value FROM settings WHERE key = 'total_tables'`).Scan(&value)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (r *SettingsRepository) SetTotalTables(count int) error {
	_, err := r.db.Exec(
		`INSERT INTO settings (key, value) VALUES ('total_tables', $1)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		strconv.Itoa(count),
	)
	return err
}
