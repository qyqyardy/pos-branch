package repository

import (
	"database/sql"
	"errors"

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

	PaymentStatus string
	PaymentToken  string
	KitchenStatus string
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
			payment_method, received, change,
			payment_status, payment_token, kitchen_status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		orderID, userID, total,
		meta.OrderType, meta.TableNo, meta.GuestCount, meta.CustomerName,
		meta.PaymentMethod, meta.Received, meta.Change,
		meta.PaymentStatus, meta.PaymentToken, meta.KitchenStatus,
	)

	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	for _, i := range items {
		// Deduct stock and verify availability
		res, err := tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1",
			i.Qty, i.ProductID,
		)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			tx.Rollback()
			return uuid.Nil, errors.New("insufficient stock for product " + i.ProductID.String())
		}

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

func (r *OrderRepo) UpdateStatus(id uuid.UUID, status string) error {
	_, err := r.DB.Exec("UPDATE orders SET payment_status = $1 WHERE id = $2", status, id)
	return err
}

func (r *OrderRepo) UpdateKitchenStatus(id uuid.UUID, status string) error {
	_, err := r.DB.Exec("UPDATE orders SET kitchen_status = $1 WHERE id = $2", status, id)
	return err
}

func (r *OrderRepo) GetByID(id uuid.UUID) (*OrderMeta, error) {
	var meta OrderMeta
	err := r.DB.QueryRow(
		"SELECT payment_status, payment_token FROM orders WHERE id = $1",
		id,
	).Scan(&meta.PaymentStatus, &meta.PaymentToken)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}
