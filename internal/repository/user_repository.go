package repository

import (
	"database/sql"

	"github.com/yourname/qr-ordering-system/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}


func (r *UserRepository) Create(username, passwordHash string, role models.Role) (*models.User, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		username, passwordHash, role,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.User{ID: id, Username: username, Role: role}, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
