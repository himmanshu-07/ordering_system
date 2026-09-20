package main

import (
	"log"
	"net/http"

	"github.com/yourname/qr-ordering-system/internal/admin"
	"github.com/yourname/qr-ordering-system/internal/config"
	"github.com/yourname/qr-ordering-system/internal/database"
	"github.com/yourname/qr-ordering-system/internal/handlers"
	"github.com/yourname/qr-ordering-system/internal/kitchen"
	"github.com/yourname/qr-ordering-system/internal/repository"
	"github.com/yourname/qr-ordering-system/internal/routes"
	"github.com/yourname/qr-ordering-system/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	tableRepo := repository.NewTableRepository(db)
	userRepo := repository.NewUserRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	orderService := services.NewOrderService(orderRepo, productRepo)
	qrService := services.NewQRService(tableRepo, cfg.BaseURL)
	dashboardService := kitchen.NewDashboardService(orderRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	productService := services.NewProductService(productRepo)
	tableService := services.NewTableService(orderRepo, settingsRepo)

	dashboardSvc := admin.NewDashboardService(orderRepo)

	productHandler := handlers.NewProductHandler(productRepo, categoryRepo, productService)
	kitchenHandler := handlers.NewKitchenHandler(dashboardService)
	adminHandler := admin.NewAdminHandler(qrService, dashboardSvc)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	authHandler := handlers.NewAuthHandler(authService)
	orderHandler := handlers.NewOrderHandler(orderService, orderRepo)
	tableHandler := handlers.NewTableHandler(tableService)
	configHandler := handlers.NewConfigHandler(cfg.UPIID, cfg.WhatsAppNumber)


	router := routes.New(productHandler, orderHandler, adminHandler, kitchenHandler, categoryHandler, authHandler, authService, tableHandler, configHandler)

	log.Printf("server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
