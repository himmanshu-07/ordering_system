package repository

import (
	"database/sql"

	"github.com/yourname/qr-ordering-system/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAvailableByCategory(categoryID int) ([]models.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, category_id, name, description, price, image_url, is_available
		 FROM products WHERE category_id = $1 AND is_available = true`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.IsAvailable); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		`SELECT id, category_id, name, description, price, image_url, is_available
		 FROM products WHERE id = $1`, id,
	).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.IsAvailable)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(categoryID int, name, description string, price float64, imageURL string) (*models.Product, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO products (category_id, name, description, price, image_url, is_available)
		 VALUES ($1, $2, $3, $4, $5, true) RETURNING id`,
		categoryID, name, description, price, imageURL,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.Product{
		ID: id, CategoryID: categoryID, Name: name,
		Description: description, Price: price, ImageURL: imageURL,
		IsAvailable: true,
	}, nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) ToggleAvailability(id int) (*models.Product, error) {
	_, err := r.db.Exec(
		`UPDATE products SET is_available = NOT is_available WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}


func (r *ProductRepository) GetAll() ([]models.Product, error) {
	rows, err := r.db.Query(`SELECT id, category_id, name, description, price, image_url, is_available FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.IsAvailable); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
