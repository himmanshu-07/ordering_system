package models

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleKitchen Role = "kitchen"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
