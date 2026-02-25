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
	"pos-backend/internal/repository"
)

type LedgerHandler struct {
	Repo *repository.LedgerRepo
}

func NewLedgerHandler(db *sql.DB) *LedgerHandler {
	return &LedgerHandler{
		Repo: &repository.LedgerRepo{DB: db},
	}
}

func (h *LedgerHandler) List(w http.ResponseWriter, r *http.Request) {
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "Invalid date (expected YYYY-MM-DD)", 400)
		return
	}

	list, err := h.Repo.ListByDate(date)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	out := make([]map[string]any, 0, len(list))
	for _, e := range list {
		out = append(out, map[string]any{
			"id":          e.ID.String(),
			"entry_date":  e.EntryDate.Format("2006-01-02"),
			"type":        e.Type,
			"amount":      e.Amount,
			"payment_method": func() string {
				if e.Payment == "" {
					return "cash"
				}
				return e.Payment
			}(),
			"category": func() string {
				if strings.TrimSpace(e.Category) == "" {
					return "general"
				}
				return e.Category
			}(),
			"description": nullIfEmpty(e.Description.String),
			"created_at":  e.CreatedAt,
			"created_by": map[string]any{
				"id":    nullIfEmpty(e.CreatedByID.String),
				"name":  nullIfEmpty(e.CreatedByName.String),
				"email": nullIfEmpty(e.CreatedByEmail.String),
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *LedgerHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		EntryDate   string  `json:"entry_date"`
		Type        string  `json:"type"`
		Amount      int64   `json:"amount"`
		PaymentMethod string `json:"payment_method"`
		Category    string  `json:"category"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	date := strings.TrimSpace(req.EntryDate)
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "Invalid entry_date (expected YYYY-MM-DD)", 400)
		return
	}

	typ := strings.ToLower(strings.TrimSpace(req.Type))
	if typ != "income" && typ != "expense" {
		http.Error(w, "Invalid type (income|expense)", 400)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "Amount must be > 0", 400)
		return
	}

	paymentMethod := strings.ToLower(strings.TrimSpace(req.PaymentMethod))
	if paymentMethod == "" {
		paymentMethod = "cash"
	}
	if paymentMethod != "cash" && paymentMethod != "bank" {
		http.Error(w, "Invalid payment_method (cash|bank)", 400)
		return
	}

	category := strings.TrimSpace(req.Category)
	if category == "" {
		category = "general"
	}
	if len(category) > 64 {
		http.Error(w, "Category too long", 400)
		return
	}

	var desc *string
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		if d != "" {
			desc = &d
		}
	}

	id, err := h.Repo.Create(date, typ, req.Amount, paymentMethod, category, desc, userID)
	if err != nil {
		http.Error(w, "Failed to create entry", 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id.String(),
	})
}

func (h *LedgerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid id", 400)
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete entry", 400)
		return
	}

	w.WriteHeader(204)
}

func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
