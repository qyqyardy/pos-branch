package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"pos-backend/internal/model"
)

type ProductRepo struct {
	DB *sql.DB
}

func (r *ProductRepo) Create(name string, price int64) error {
	_, err := r.DB.Exec(
		"INSERT INTO products (id,name,price) VALUES ($1,$2,$3)",
		uuid.New(), name, price,
	)
	return err
}

func (r *ProductRepo) GetAll() ([]model.Product, error) {
	rows, err := r.DB.Query("SELECT id,name,price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ProductRepo) FindByID(id uuid.UUID) (*model.Product, error) {
	var p model.Product
	err := r.DB.QueryRow(
		"SELECT id,name,price FROM products WHERE id=$1",
		id,
	).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
