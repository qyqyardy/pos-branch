package handler

import (
	"database/sql"
	"net/http"

	"pos-backend/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) http.Handler {

	r := mux.NewRouter()

	// CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	auth := NewAuthHandler(db)
	settings := NewSettingsHandler(db)
	userAdmin := NewUserAdminHandler(db)
	ledger := NewLedgerHandler(db)
	product := NewProductHandler(db)
	order := NewOrderHandler(db)
	analytics := NewAnalyticsHandler(db)

	r.HandleFunc("/register", auth.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", auth.Login).Methods("POST", "OPTIONS")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWT)

	api.HandleFunc("/me", auth.Me).Methods("GET", "OPTIONS")
	api.HandleFunc("/settings/store", settings.GetStore).Methods("GET", "OPTIONS")
	api.HandleFunc("/products", product.GetProducts).Methods("GET", "OPTIONS")
	api.HandleFunc("/analytics/sales", analytics.GetSalesSummary).Methods("GET", "OPTIONS")
	api.HandleFunc("/analytics/top-products", analytics.GetTopProducts).Methods("GET", "OPTIONS")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRoles("admin"))
	admin.HandleFunc("/products", product.CreateProduct).Methods("POST", "OPTIONS")
	admin.HandleFunc("/products/{id}", product.UpdateProduct).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/products/{id}", product.DeleteProduct).Methods("DELETE", "OPTIONS")
	admin.HandleFunc("/users", userAdmin.ListUsers).Methods("GET", "OPTIONS")

	admin.HandleFunc("/users", userAdmin.CreateUser).Methods("POST", "OPTIONS")
	admin.HandleFunc("/users/{id}", userAdmin.UpdateUser).Methods("PATCH", "OPTIONS")
	admin.HandleFunc("/users/{id}", userAdmin.DeleteUser).Methods("DELETE", "OPTIONS")

	api.Handle("/orders", middleware.RequireRoles("admin", "cashier")(http.HandlerFunc(order.CreateOrder))).Methods("POST", "OPTIONS")
	api.Handle("/orders", middleware.RequireRoles("admin", "finance")(http.HandlerFunc(order.ListOrders))).Methods("GET", "OPTIONS")
	api.Handle("/orders/{id}", middleware.RequireRoles("admin", "finance")(http.HandlerFunc(order.GetOrder))).Methods("GET", "OPTIONS")

	api.Handle("/settings/store", middleware.RequireRoles("admin")(http.HandlerFunc(settings.UpdateStore))).Methods("PUT", "OPTIONS")

	api.Handle("/ledger", middleware.RequireRoles("admin", "finance")(http.HandlerFunc(ledger.List))).Methods("GET", "OPTIONS")
	api.Handle("/ledger", middleware.RequireRoles("admin", "finance")(http.HandlerFunc(ledger.Create))).Methods("POST", "OPTIONS")
	api.Handle("/ledger/{id}", middleware.RequireRoles("admin", "finance")(http.HandlerFunc(ledger.Delete))).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/webhooks/midtrans", order.MidtransWebhook).Methods("POST", "OPTIONS")

	return r
}
