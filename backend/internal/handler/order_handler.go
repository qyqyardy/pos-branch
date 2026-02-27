package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"pos-backend/internal/middleware"
	"pos-backend/internal/model"
	"pos-backend/internal/repository"
	"pos-backend/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	resp := map[string]any{
		"order_id": orderID.String(),
		"total":    total,
	}

	// Fetch saved payment info
	saved, err := h.Service.OrderRepo.GetByID(orderID)
	if err == nil {
		resp["payment_status"] = saved.PaymentStatus
		resp["payment_token"] = saved.PaymentToken
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
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

	// Optimization: Use date range instead of ::date to allow standard index usage
	start := date + " 00:00:00"
	end := date + " 23:59:59.999"

	rows, err := h.Service.OrderRepo.DB.Query(
		`SELECT o.id, o.total, o.created_at,
		        o.order_type, o.table_no, o.guest_count, o.customer_name,
		        o.payment_method, o.received, o.change, o.kitchen_status,
		        u.id, u.name, u.email
		   FROM orders o
		   LEFT JOIN users u ON u.id = o.cashier_id
		  WHERE o.created_at >= $1 AND o.created_at <= $2
		  ORDER BY o.created_at DESC`,
		start, end,
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
			kitchen      string
			cashierID    sql.NullString
			cashierName  sql.NullString
			cashierEmail sql.NullString
		)

		if err := rows.Scan(
			&id, &total, &createdAt,
			&orderType, &tableNo, &guestCount, &customerName,
			&payMethod, &received, &change, &kitchen,
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
			"received":       nullIfZeroInt64(received),
			"change":         nullIfZeroInt64(change),
			"kitchen_status": kitchen,
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
		kitchen      string
		cashierName  sql.NullString
		cashierEmail sql.NullString
	)

	err = h.Service.OrderRepo.DB.QueryRow(
		`SELECT o.total, o.created_at,
		        o.order_type, o.table_no, o.guest_count, o.customer_name,
		        o.payment_method, o.received, o.change, o.kitchen_status,
		        u.name, u.email
		   FROM orders o
		   LEFT JOIN users u ON u.id = o.cashier_id
		  WHERE o.id = $1`,
		id,
	).Scan(
		&total, &createdAt,
		&orderType, &tableNo, &guestCount, &customerName,
		&payMethod, &received, &change, &kitchen,
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
		"received":       nullIfZeroInt64(received),
		"change":         nullIfZeroInt64(change),
		"kitchen_status": kitchen,
		"cashier": map[string]any{
			"name":  nullIfBlank(cashierName.String),
			"email": nullIfBlank(cashierEmail.String),
		},
		"items": items,
	})
}

func (h *OrderHandler) MidtransWebhook(w http.ResponseWriter, r *http.Request) {
	var notification struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid notification body", 400)
		return
	}

	orderID, err := uuid.Parse(notification.OrderID)
	if err != nil {
		parts := strings.Split(notification.OrderID, "-")
		if len(parts) > 1 {
			orderID, err = uuid.Parse(parts[len(parts)-1])
		}
	}

	if err != nil {
		w.WriteHeader(200)
		return
	}

	status := "pending"
	if notification.TransactionStatus == "capture" || notification.TransactionStatus == "settlement" {
		if notification.FraudStatus == "accept" || notification.FraudStatus == "" {
			status = "completed"
		}
	} else if notification.TransactionStatus == "deny" || notification.TransactionStatus == "expire" || notification.TransactionStatus == "cancel" {
		status = "failed"
	}

	if err := h.Service.OrderRepo.UpdateStatus(orderID, status); err != nil {
		http.Error(w, "Failed to update order status", 500)
		return
	}

	w.WriteHeader(200)
}

func (h *OrderHandler) UpdateKitchenStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", 400)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	valid := false
	for _, s := range []string{"pending", "preparing", "ready", "done"} {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "Invalid kitchen status", 400)
		return
	}

	if err := h.Service.OrderRepo.UpdateKitchenStatus(id, status); err != nil {
		http.Error(w, "Failed to update kitchen status", 500)
		return
	}

	w.WriteHeader(204)
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
