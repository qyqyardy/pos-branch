package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"pos-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	DB *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name         string `json:"name"`
		Price        int64  `json:"price"`
		ImageDataURL string `json:"image_data_url"`
		IsActive     *bool  `json:"is_active"`
		Stock        int    `json:"stock"`
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	repo := repository.ProductRepo{DB: h.DB}
	err := repo.Create(req.Name, req.Price, req.ImageDataURL, isActive, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(201)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	var req struct {
		Name         string `json:"name"`
		Price        int64  `json:"price"`
		ImageDataURL string `json:"image_data_url"`
		IsActive     *bool  `json:"is_active"`
		Stock        int    `json:"stock"`
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	repo := repository.ProductRepo{DB: h.DB}
	err = repo.Update(id, req.Name, req.Price, req.ImageDataURL, isActive, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(200)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	repo := repository.ProductRepo{DB: h.DB}
	err = repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(204)
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
