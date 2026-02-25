package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"pos-backend/internal/model"
)

type UserRepo struct {
	DB *sql.DB
}

type UserListItem struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRow(
		"SELECT id,name,email,password_hash,role FROM users WHERE email=$1",
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindByID(id uuid.UUID) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRow(
		"SELECT id,name,email,password_hash,role FROM users WHERE id=$1",
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)",
		id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) List() ([]UserListItem, error) {
	rows, err := r.DB.Query(
		"SELECT id,name,email,role,created_at FROM users ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []UserListItem
	for rows.Next() {
		var u UserListItem
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *UserRepo) CountByRole(role string) (int, error) {
	var c int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role=$1", role).Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *UserRepo) Create(id uuid.UUID, name, email, passwordHash, role string) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (id,name,email,password_hash,role) VALUES ($1,$2,$3,$4,$5)",
		id, name, email, passwordHash, role,
	)
	return err
}

func (r *UserRepo) Update(user *model.User) error {
	_, err := r.DB.Exec(
		"UPDATE users SET name=$1, email=$2, password_hash=$3, role=$4 WHERE id=$5",
		user.Name, user.Email, user.PasswordHash, user.Role, user.ID,
	)
	return err
}

func (r *UserRepo) Delete(id uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
