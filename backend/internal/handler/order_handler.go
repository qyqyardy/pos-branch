package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"pos-backend/internal/middleware"
	"pos-backend/internal/model"
	"pos-backend/internal/repository"
	"pos-backend/internal/service"
)

type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(db *sql.DB) *OrderHandler {

	orderRepo := &repository.OrderRepo{DB: db}
	productRepo := &repository.ProductRepo{DB: db}
	svc := service.NewOrderService(orderRepo, productRepo)

	return &OrderHandler{Service: svc}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	// Token might still be valid while the user row is gone (e.g. DB reset).
	userRepo := repository.UserRepo{DB: h.Service.OrderRepo.DB}
	exists, err := userRepo.ExistsByID(userID)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if !exists {
		http.Error(w, "Unauthorized", 401)
		return
	}

	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	orderID, total, err := h.Service.CreateOrder(userID, req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]any{
		"order_id": orderID.String(),
		"total":    total,
	})
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "Invalid date (expected YYYY-MM-DD)", 400)
		return
	}

	rows, err := h.Service.OrderRepo.DB.Query(
		`SELECT o.id, o.total, o.created_at,
		        o.order_type, o.table_no, o.guest_count, o.customer_name,
		        o.payment_method, o.received, o.change,
		        u.id, u.name, u.email
		   FROM orders o
		   LEFT JOIN users u ON u.id = o.cashier_id
		  WHERE o.created_at::date = $1::date
		  ORDER BY o.created_at DESC`,
		date,
	)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	defer rows.Close()

	out := make([]map[string]any, 0)

	for rows.Next() {
		var (
			id           uuid.UUID
			total        int64
			createdAt    time.Time
			orderType    string
			tableNo      sql.NullString
			guestCount   sql.NullInt32
			customerName sql.NullString
			payMethod    string
			received     sql.NullInt64
			change       sql.NullInt64
			cashierID    sql.NullString
			cashierName  sql.NullString
			cashierEmail sql.NullString
		)

		if err := rows.Scan(
			&id, &total, &createdAt,
			&orderType, &tableNo, &guestCount, &customerName,
			&payMethod, &received, &change,
			&cashierID, &cashierName, &cashierEmail,
		); err != nil {
			http.Error(w, "Server error", 500)
			return
		}

		out = append(out, map[string]any{
			"id":            id.String(),
			"total":         total,
			"created_at":    createdAt,
			"order_type":    orderType,
			"table_no":      nullIfBlank(tableNo.String),
			"guest_count":   nullIfZeroInt32(guestCount),
			"customer_name": nullIfBlank(customerName.String),
			"payment_method": func() string {
				if payMethod == "" {
					return "cash"
				}
				return payMethod
			}(),
			"received": nullIfZeroInt64(received),
			"change":   nullIfZeroInt64(change),
			"cashier": map[string]any{
				"id":    nullIfBlank(cashierID.String),
				"name":  nullIfBlank(cashierName.String),
				"email": nullIfBlank(cashierEmail.String),
			},
		})
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", 400)
		return
	}

	var (
		total        int64
		createdAt    time.Time
		orderType    string
		tableNo      sql.NullString
		guestCount   sql.NullInt32
		customerName sql.NullString
		payMethod    string
		received     sql.NullInt64
		change       sql.NullInt64
		cashierName  sql.NullString
		cashierEmail sql.NullString
	)

	err = h.Service.OrderRepo.DB.QueryRow(
		`SELECT o.total, o.created_at,
		        o.order_type, o.table_no, o.guest_count, o.customer_name,
		        o.payment_method, o.received, o.change,
		        u.name, u.email
		   FROM orders o
		   LEFT JOIN users u ON u.id = o.cashier_id
		  WHERE o.id = $1`,
		id,
	).Scan(
		&total, &createdAt,
		&orderType, &tableNo, &guestCount, &customerName,
		&payMethod, &received, &change,
		&cashierName, &cashierEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", 404)
			return
		}
		http.Error(w, "Server error", 500)
		return
	}

	rows, err := h.Service.OrderRepo.DB.Query(
		`SELECT p.name, oi.qty, oi.price
		   FROM order_items oi
		   JOIN products p ON p.id = oi.product_id
		  WHERE oi.order_id = $1
		  ORDER BY oi.id ASC`,
		id,
	)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var (
			name  string
			qty   int
			price int64
		)
		if err := rows.Scan(&name, &qty, &price); err != nil {
			http.Error(w, "Server error", 500)
			return
		}
		items = append(items, map[string]any{
			"name":       name,
			"qty":        qty,
			"price":      price,
			"line_total": int64(qty) * price,
		})
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":            id.String(),
		"total":         total,
		"created_at":    createdAt,
		"order_type":    orderType,
		"table_no":      nullIfBlank(tableNo.String),
		"guest_count":   nullIfZeroInt32(guestCount),
		"customer_name": nullIfBlank(customerName.String),
		"payment_method": func() string {
			if payMethod == "" {
				return "cash"
			}
			return payMethod
		}(),
		"received": nullIfZeroInt64(received),
		"change":   nullIfZeroInt64(change),
		"cashier": map[string]any{
			"name":  nullIfBlank(cashierName.String),
			"email": nullIfBlank(cashierEmail.String),
		},
		"items": items,
	})
}

func nullIfBlank(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullIfZeroInt32(v sql.NullInt32) any {
	if !v.Valid || v.Int32 <= 0 {
		return nil
	}
	return v.Int32
}

func nullIfZeroInt64(v sql.NullInt64) any {
	if !v.Valid || v.Int64 == 0 {
		return nil
	}
	return v.Int64
}
