package service

import "errors"

var (
	ErrEmptyCart       = errors.New("cart is empty")
	ErrInvalidQty      = errors.New("invalid quantity")
	ErrInvalidProduct  = errors.New("invalid product")
	ErrProductNotFound = errors.New("product not found")
)
