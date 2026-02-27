package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type AnalyticsHandler struct {
	DB *sql.DB
}

func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

type DailySale struct {
	Date  string `json:"date"`
	Total int64  `json:"total"`
}

type TopProduct struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

func (h *AnalyticsHandler) GetSalesSummary(w http.ResponseWriter, r *http.Request) {
	// Daily revenue for last 7 days
	rows, err := h.DB.Query(`
		SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as day, SUM(total)
		FROM orders
		WHERE created_at >= CURRENT_DATE - INTERVAL '6 days'
		  AND payment_status = 'completed'
		GROUP BY day
		ORDER BY day ASC
	`)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	defer rows.Close()

	salesMap := make(map[string]int64)
	for rows.Next() {
		var day string
		var total int64
		if err := rows.Scan(&day, &total); err != nil {
			http.Error(w, "Server error", 500)
			return
		}
		salesMap[day] = total
	}

	// Ensure all 7 days are present
	results := make([]DailySale, 0)
	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		results = append(results, DailySale{
			Date:  day,
			Total: salesMap[day],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *AnalyticsHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT p.name, SUM(oi.qty) as total_qty
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON p.id = oi.product_id
		WHERE o.payment_status = 'completed'
		  AND o.created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 5
	`)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	defer rows.Close()

	results := make([]TopProduct, 0)
	for rows.Next() {
		var p TopProduct
		if err := rows.Scan(&p.Name, &p.Qty); err != nil {
			http.Error(w, "Server error", 500)
			return
		}
		results = append(results, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
