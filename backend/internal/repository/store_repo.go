package repository

import (
	"database/sql"
	"errors"

	"pos-backend/internal/model"

	"github.com/lib/pq"
)

type StoreRepo struct {
	DB *sql.DB
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

	err := r.DB.QueryRow(
		`SELECT name, tagline, address_line1, address_line2, phone, logo_data_url, footer_message
		 FROM store_settings
		 WHERE id = 1`,
	).Scan(&s.Name, &s.Tagline, &s.AddressLine1, &s.AddressLine2, &s.Phone, &s.LogoDataURL, &s.FooterMessage)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StoreRepo) Upsert(s model.StoreSettings) error {
	_, err := r.DB.Exec(
		`INSERT INTO store_settings (id, name, tagline, address_line1, address_line2, phone, logo_data_url, footer_message, updated_at)
		 VALUES (1, $1, $2, $3, $4, $5, $6, $7, now())
		 ON CONFLICT (id) DO UPDATE
		 SET name = EXCLUDED.name,
		     tagline = EXCLUDED.tagline,
		     address_line1 = EXCLUDED.address_line1,
		     address_line2 = EXCLUDED.address_line2,
		     phone = EXCLUDED.phone,
		     logo_data_url = EXCLUDED.logo_data_url,
		     footer_message = EXCLUDED.footer_message,
		     updated_at = now()`,
		s.Name, s.Tagline, s.AddressLine1, s.AddressLine2, s.Phone, s.LogoDataURL, s.FooterMessage,
	)
	return err
}
