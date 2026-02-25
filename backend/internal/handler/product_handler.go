package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"pos-backend/internal/repository"
)

type ProductHandler struct {
	DB *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Price <= 0 {
		http.Error(w, "Invalid product data", 400)
		return
	}

	repo := repository.ProductRepo{DB: h.DB}
	err := repo.Create(req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(201)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

	repo := repository.ProductRepo{DB: h.DB}
	list, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
