package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type LedgerRepo struct {
	DB *sql.DB
}

type LedgerEntry struct {
	ID          uuid.UUID
	EntryDate   time.Time
	Type        string
	Amount      int64
	Payment     string
	Category    string
	Description sql.NullString
	CreatedAt   time.Time

	CreatedByID    sql.NullString
	CreatedByName  sql.NullString
	CreatedByEmail sql.NullString
}

func (r *LedgerRepo) ListByDate(date string) ([]LedgerEntry, error) {
	rows, err := r.DB.Query(
		`SELECT l.id, l.entry_date, l.type, l.amount, l.payment_method, l.category, l.description, l.created_at,
		        u.id, u.name, u.email
		   FROM cash_ledger l
		   LEFT JOIN users u ON u.id = l.created_by
		  WHERE l.entry_date = $1::date
		  ORDER BY l.created_at DESC`,
		date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []LedgerEntry
	for rows.Next() {
		var e LedgerEntry
		if err := rows.Scan(
			&e.ID, &e.EntryDate, &e.Type, &e.Amount, &e.Payment, &e.Category, &e.Description, &e.CreatedAt,
			&e.CreatedByID, &e.CreatedByName, &e.CreatedByEmail,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *LedgerRepo) Create(entryDate string, typ string, amount int64, paymentMethod string, category string, description *string, createdBy uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.DB.Exec(
		`INSERT INTO cash_ledger (id, entry_date, type, amount, payment_method, category, description, created_by)
		 VALUES ($1, $2::date, $3, $4, $5, $6, $7, $8)`,
		id, entryDate, typ, amount, paymentMethod, category, description, createdBy,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *LedgerRepo) Delete(id uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM cash_ledger WHERE id=$1", id)
	return err
}
