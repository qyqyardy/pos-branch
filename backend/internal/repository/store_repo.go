package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"pos-backend/internal/model"

	"github.com/lib/pq"
)

type StoreRepo struct {
	DB *sql.DB
}

func normalizePlan(plan string) string {
	p := strings.ToLower(strings.TrimSpace(plan))
	switch p {
	case "standard", "premium":
		return p
	default:
		return "premium"
	}
}

func isUndefinedColumn(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// 42703 = undefined_column
		return string(pqErr.Code) == "42703"
	}
	return false
}

func (r *StoreRepo) Get() (*model.StoreSettings, error) {
	var s model.StoreSettings

	// Newer schema includes plan + paid_until. Keep backward compatibility for old DB volumes.
	err := r.DB.QueryRow(
		`SELECT name, tagline, address_line1, address_line2, phone, logo_data_url, plan, paid_until
		 FROM store_settings
		 WHERE id = 1`,
	).Scan(&s.Name, &s.Tagline, &s.AddressLine1, &s.AddressLine2, &s.Phone, &s.LogoDataURL, &s.Plan, &s.PaidUntil)
	if err != nil {
		if isUndefinedColumn(err) {
			err2 := r.DB.QueryRow(
				`SELECT name, tagline, address_line1, address_line2, phone, logo_data_url
				 FROM store_settings
				 WHERE id = 1`,
			).Scan(&s.Name, &s.Tagline, &s.AddressLine1, &s.AddressLine2, &s.Phone, &s.LogoDataURL)
			if err2 != nil {
				return nil, err2
			}

			s.Plan = "premium"
			s.PaidUntil = time.Now().AddDate(10, 0, 0)
			return &s, nil
		}

		return nil, err
	}

	s.Plan = normalizePlan(s.Plan)
	if s.PaidUntil.IsZero() {
		s.PaidUntil = time.Now().AddDate(10, 0, 0)
	}
	return &s, nil
}

func (r *StoreRepo) Upsert(s model.StoreSettings) error {
	_, err := r.DB.Exec(
		`INSERT INTO store_settings (id, name, tagline, address_line1, address_line2, phone, logo_data_url, updated_at)
		 VALUES (1, $1, $2, $3, $4, $5, $6, now())
		 ON CONFLICT (id) DO UPDATE
		 SET name = EXCLUDED.name,
		     tagline = EXCLUDED.tagline,
		     address_line1 = EXCLUDED.address_line1,
		     address_line2 = EXCLUDED.address_line2,
		     phone = EXCLUDED.phone,
		     logo_data_url = EXCLUDED.logo_data_url,
		     updated_at = now()`,
		s.Name, s.Tagline, s.AddressLine1, s.AddressLine2, s.Phone, s.LogoDataURL,
	)
	return err
}
