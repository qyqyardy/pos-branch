package repository

import (
	"database/sql"

	"pos-backend/internal/model"
)

type StoreRepo struct {
	DB *sql.DB
}

func (r *StoreRepo) Get() (*model.StoreSettings, error) {
	var s model.StoreSettings
	err := r.DB.QueryRow(
		`SELECT name, tagline, address_line1, address_line2, phone, logo_data_url
		 FROM store_settings
		 WHERE id = 1`,
	).Scan(&s.Name, &s.Tagline, &s.AddressLine1, &s.AddressLine2, &s.Phone, &s.LogoDataURL)
	if err != nil {
		return nil, err
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
