package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"pos-backend/internal/middleware"
	"pos-backend/internal/repository"
	"pos-backend/internal/service"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", 400)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", 500)
		return
	}

	// Prevent privilege escalation: role is assigned server-side.
	role := "cashier"

	_, err = h.DB.Exec(
		"INSERT INTO users (id,name,email,password_hash,role) VALUES (uuid_generate_v4(),$1,$2,$3,$4)",
		req.Name, req.Email, string(hash), role,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(201)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing email or password", 400)
		return
	}

	repo := repository.UserRepo{DB: h.DB}
	user, err := repo.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", 401)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Wrong password", 401)
		return
	}

	token, err := service.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {

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

	repo := repository.UserRepo{DB: h.DB}
	user, err := repo.FindByID(userID)
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":    user.ID.String(),
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
