package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourname/qr-ordering-system/internal/config"
	"github.com/yourname/qr-ordering-system/internal/database"
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
	"github.com/yourname/qr-ordering-system/internal/services"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: seedadmin <username> <password> <role: admin|kitchen>")
		os.Exit(1)
	}
	username, password, role := os.Args[1], os.Args[2], os.Args[3]

	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	user, err := authService.CreateUser(username, password, models.Role(role))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created user: %s (role: %s, id: %d)\n", user.Username, user.Role, user.ID)
}
