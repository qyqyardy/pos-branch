package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"pos-backend/internal/middleware"
	"pos-backend/internal/repository"
)

type UserAdminHandler struct {
	Repo      *repository.UserRepo
	StoreRepo *repository.StoreRepo
}

func NewUserAdminHandler(db *sql.DB) *UserAdminHandler {
	return &UserAdminHandler{
		Repo:      &repository.UserRepo{DB: db},
		StoreRepo: &repository.StoreRepo{DB: db},
	}
}

func normalizeRole(role string) string {
	r := strings.ToLower(strings.TrimSpace(role))
	if r == "superadmin" {
		return "admin"
	}
	return r
}

func isValidRole(role string) bool {
	switch role {
	case "admin", "cashier", "finance":
		return true
	default:
		return false
	}
}

func (h *UserAdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.Repo.List()
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	out := make([]map[string]any, 0, len(list))
	for _, u := range list {
		out = append(out, map[string]any{
			"id":         u.ID.String(),
			"name":       u.Name,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *UserAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)
	role := normalizeRole(req.Role)
	if role == "" {
		role = "cashier"
	}

	if name == "" || email == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Missing required fields", 400)
		return
	}
	if !isValidRole(role) {
		http.Error(w, "Invalid role", 400)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", 500)
		return
	}

	id := uuid.New()
	err = h.Repo.Create(id, name, email, string(hash), role)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id.String(),
	})
}

func (h *UserAdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user id", 400)
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	user, err := h.Repo.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", 404)
			return
		}
		http.Error(w, "Server error", 500)
		return
	}

	if req.Name != nil {
		user.Name = strings.TrimSpace(*req.Name)
	}
	if req.Email != nil {
		user.Email = strings.TrimSpace(*req.Email)
	}
	if req.Role != nil {
		newRole := normalizeRole(*req.Role)
		if !isValidRole(newRole) {
			http.Error(w, "Invalid role", 400)
			return
		}

		// Prevent locking out: must keep at least 1 admin.
		if user.Role == "admin" && newRole != "admin" {
			adminCount, err := h.Repo.CountByRole("admin")
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}
			if adminCount <= 1 {
				http.Error(w, "Cannot remove the last admin", 400)
				return
			}
		}

		user.Role = newRole
	}

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", 500)
			return
		}
		user.PasswordHash = string(hash)
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Name and email are required", 400)
		return
	}

	err = h.Repo.Update(user)
	if err != nil {
		http.Error(w, err.Error(), 400)
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

func (h *UserAdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user id", 400)
		return
	}

	currentUserIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if ok && currentUserIDStr != "" && currentUserIDStr == id.String() {
		http.Error(w, "Cannot delete your own account", 400)
		return
	}

	user, err := h.Repo.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", 404)
			return
		}
		http.Error(w, "Server error", 500)
		return
	}

	if user.Role == "admin" {
		adminCount, err := h.Repo.CountByRole("admin")
		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}
		if adminCount <= 1 {
			http.Error(w, "Cannot delete the last admin", 400)
			return
		}
	}

	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete user", 400)
		return
	}

	w.WriteHeader(204)
}
