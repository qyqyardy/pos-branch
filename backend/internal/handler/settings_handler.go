package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"pos-backend/internal/model"
	"pos-backend/internal/repository"
)

type SettingsHandler struct {
	StoreRepo *repository.StoreRepo
}

func NewSettingsHandler(db *sql.DB) *SettingsHandler {
	return &SettingsHandler{
		StoreRepo: &repository.StoreRepo{DB: db},
	}
}

func (h *SettingsHandler) GetStore(w http.ResponseWriter, r *http.Request) {
	s, err := h.StoreRepo.Get()
	if err != nil {
		// Fresh DB without seed row (or manual deletion). Return a safe default.
		if err == sql.ErrNoRows {
			s = &model.StoreSettings{
				Name:          "WARKOP",
				Tagline:       "Point of Sale",
				AddressLine1:  "",
				AddressLine2:  "",
				Phone:         "",
				LogoDataURL:   "",
				FooterMessage: "",
			}
		} else {
			http.Error(w, "Server error", 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"name":           s.Name,
		"tagline":        s.Tagline,
		"address_lines":  []string{s.AddressLine1, s.AddressLine2},
		"phone":          s.Phone,
		"logo_data_url":  s.LogoDataURL,
		"footer_message": s.FooterMessage,
	})
}

func (h *SettingsHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name          string   `json:"name"`
		Tagline       string   `json:"tagline"`
		AddressLines  []string `json:"address_lines"`
		AddressLine1  string   `json:"address_line1"`
		AddressLine2  string   `json:"address_line2"`
		Phone         string   `json:"phone"`
		LogoDataURL   *string  `json:"logo_data_url"`
		FooterMessage string   `json:"footer_message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	name := strings.TrimSpace(req.Name)
	tagline := strings.TrimSpace(req.Tagline)
	phone := strings.TrimSpace(req.Phone)

	addr1 := strings.TrimSpace(req.AddressLine1)
	addr2 := strings.TrimSpace(req.AddressLine2)
	if len(req.AddressLines) > 0 {
		addr1 = strings.TrimSpace(req.AddressLines[0])
	}
	if len(req.AddressLines) > 1 {
		addr2 = strings.TrimSpace(req.AddressLines[1])
	}

	if name == "" {
		http.Error(w, "Store name is required", 400)
		return
	}

	logo := ""
	if req.LogoDataURL == nil {
		current, err := h.StoreRepo.Get()
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Server error", 500)
			return
		}
		if current != nil {
			logo = current.LogoDataURL
		}
	} else {
		logo = strings.TrimSpace(*req.LogoDataURL)

		// Guard: keep logo small enough for DB/localStorage and printing.
		if len(logo) > 450000 {
			http.Error(w, "Logo too large (max ~450KB data URL)", 400)
			return
		}

		if logo != "" && !strings.HasPrefix(logo, "data:image/") {
			http.Error(w, "Invalid logo format", 400)
			return
		}
	}

	err := h.StoreRepo.Upsert(model.StoreSettings{
		Name:          name,
		Tagline:       tagline,
		AddressLine1:  addr1,
		AddressLine2:  addr2,
		Phone:         phone,
		LogoDataURL:   logo,
		FooterMessage: strings.TrimSpace(req.FooterMessage),
	})
	if err != nil {
		http.Error(w, "Failed to save settings", 400)
		return
	}

	h.GetStore(w, r)
}
