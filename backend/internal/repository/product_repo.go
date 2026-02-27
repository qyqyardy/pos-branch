package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"pos-backend/internal/model"
)

type ProductRepo struct {
	DB *sql.DB
}

func (r *ProductRepo) Create(name string, price int64, image string) error {
	_, err := r.DB.Exec(
		"INSERT INTO products (id,name,price,image_data_url) VALUES ($1,$2,$3,$4)",
		uuid.New(), name, price, image,
	)
	return err
}

func (r *ProductRepo) GetAll() ([]model.Product, error) {
	rows, err := r.DB.Query("SELECT id,name,price,image_data_url FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Product
	for rows.Next() {
		var p model.Product
		var img sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &img); err != nil {
			return nil, err
		}
		p.ImageDataURL = img.String
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ProductRepo) FindByID(id uuid.UUID) (*model.Product, error) {
	var p model.Product
	var img sql.NullString
	err := r.DB.QueryRow(
		"SELECT id,name,price,image_data_url FROM products WHERE id=$1",
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &img)
	if err != nil {
		return nil, err
	}
	p.ImageDataURL = img.String
	return &p, nil
}

func (r *ProductRepo) Update(id uuid.UUID, name string, price int64, image string) error {
	_, err := r.DB.Exec(
		"UPDATE products SET name=$1, price=$2, image_data_url=$3 WHERE id=$4",
		name, price, image, id,
	)
	return err
}

func (r *ProductRepo) Delete(id uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}
