package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         string
}