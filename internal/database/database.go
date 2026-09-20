package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrating db: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tables (
		id SERIAL PRIMARY KEY,
		table_number INTEGER NOT NULL UNIQUE,
		qr_code TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		display_order INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES categories(id),
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		image_url TEXT,
		is_available BOOLEAN DEFAULT TRUE
	);

	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		table_id INTEGER NOT NULL REFERENCES tables(id),
		status TEXT NOT NULL DEFAULT 'pending',
		total REAL NOT NULL DEFAULT 0,
		payment_method TEXT NOT NULL DEFAULT 'cash',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id TEXT NOT NULL REFERENCES orders(id),
		product_id INTEGER NOT NULL REFERENCES products(id),
		quantity INTEGER NOT NULL,
		unit_price REAL NOT NULL,
		notes TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO settings (key, value) VALUES ('total_tables', '20') ON CONFLICT (key) DO NOTHING`,
	); err != nil {
		return err
	}

	return nil
}
