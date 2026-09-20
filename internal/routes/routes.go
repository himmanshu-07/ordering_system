package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/yourname/qr-ordering-system/internal/admin"
	"github.com/yourname/qr-ordering-system/internal/handlers"
	appmiddleware "github.com/yourname/qr-ordering-system/internal/middleware"
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/services"
)

func New(
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	adminHandler *admin.AdminHandler,
	kitchenHandler *handlers.KitchenHandler,
	categoryHandler *handlers.CategoryHandler,
	authHandler *handlers.AuthHandler,
	authService *services.AuthService,
	tableHandler *handlers.TableHandler,
	configHandler *handlers.ConfigHandler,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Global safety net — no single IP should be able to hammer the whole app
	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Route("/api", func(r chi.Router) {
		r.Get("/menu", productHandler.GetMenu)
		r.Get("/products/{id}", productHandler.GetProduct)
		r.Get("/categories", categoryHandler.GetCategories)
		r.Get("/config", configHandler.GetPublicConfig)
		r.Get("/orders/{orderID}", orderHandler.GetOrder)
		r.Patch("/orders/{orderID}/confirm-served", orderHandler.ConfirmServed)
		r.Get("/tables/next-available", tableHandler.GetNextAvailable)

		// Order placement — a real customer places at most one order per table at a time,
		// so this is generous for genuine use but blocks spam/bot abuse
		r.With(httprate.LimitByIP(10, time.Minute)).Post("/orders", orderHandler.PlaceOrder)

		// Login — the classic brute-force target, keep this tight
		r.With(httprate.LimitByIP(5, time.Minute)).Post("/auth/login", authHandler.Login)

		r.Route("/admin", func(r chi.Router) {
			r.Use(appmiddleware.RequireAuth(authService, models.RoleAdmin))
			r.Post("/tables", adminHandler.CreateTable)
			r.Post("/tables/batch", adminHandler.GenerateBatch)
			r.Get("/settings/tables", tableHandler.GetSettings)
			r.Patch("/settings/tables", tableHandler.UpdateTableCount)
			r.Post("/categories", categoryHandler.CreateCategory)
			r.Post("/products", productHandler.CreateProduct)
			r.Get("/products", productHandler.GetAllProducts)
			r.Delete("/products/{id}", productHandler.DeleteProduct)
			r.Delete("/orders/flush", orderHandler.FlushOrders)
			r.Get("/dashboard", adminHandler.GetDashboard)
			r.Patch("/products/{id}/toggle-availability", productHandler.ToggleAvailability)

		})

		r.Route("/kitchen", func(r chi.Router) {
			r.Use(appmiddleware.RequireAuth(authService, models.RoleAdmin, models.RoleKitchen))
			r.Get("/orders", kitchenHandler.GetActiveOrders)
			r.Patch("/orders/{orderID}/status", kitchenHandler.UpdateOrderStatus)
		})
	})

	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/home.html")
	})

	r.Get("/order/{tableID}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/menu.html")
	})

	r.Get("/kitchen-dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/kitchen.html")
	})

	r.Get("/admin-dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/admin.html")
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/login.html")
	})

	return r
}
