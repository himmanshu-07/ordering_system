package repository

import (
	"database/sql"

	"github.com/yourname/qr-ordering-system/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id, name, display_order FROM categories ORDER BY display_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.DisplayOrder); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(name string, displayOrder int) (*models.Category, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO categories (name, display_order) VALUES ($1, $2) RETURNING id`,
		name, displayOrder,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.Category{ID: id, Name: name, DisplayOrder: displayOrder}, nil
}
