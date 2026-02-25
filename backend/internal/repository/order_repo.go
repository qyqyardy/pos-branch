package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type OrderRepo struct {
	DB *sql.DB
}

type OrderItemData struct {
	ProductID uuid.UUID
	Qty       int
	Price     int64
}

type OrderMeta struct {
	OrderType    string
	TableNo      *string
	GuestCount   *int
	CustomerName *string

	PaymentMethod string
	Received      *int64
	Change        *int64
}

func (r *OrderRepo) Create(userID uuid.UUID, items []OrderItemData, total int64, meta OrderMeta) (uuid.UUID, error) {

	tx, err := r.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	orderID := uuid.New()

	_, err = tx.Exec(
		`INSERT INTO orders (
			id, cashier_id, total,
			order_type, table_no, guest_count, customer_name,
			payment_method, received, change
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		orderID, userID, total,
		meta.OrderType, meta.TableNo, meta.GuestCount, meta.CustomerName,
		meta.PaymentMethod, meta.Received, meta.Change,
	)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	for _, i := range items {
		_, err = tx.Exec(
			"INSERT INTO order_items (id,order_id,product_id,qty,price) VALUES ($1,$2,$3,$4,$5)",
			uuid.New(), orderID, i.ProductID, i.Qty, i.Price,
		)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return orderID, nil
}
