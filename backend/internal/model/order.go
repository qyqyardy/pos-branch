package model

import "github.com/google/uuid"

type OrderItemInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Qty       int       `json:"qty"`
	Price     int64     `json:"price"`
}

type CreateOrderRequest struct {
	Items []OrderItemInput `json:"items"`

	// Optional metadata for receipt/reporting.
	OrderType    string  `json:"order_type"`
	TableNo      *string `json:"table_no"`
	GuestCount   *int    `json:"guest_count"`
	CustomerName *string `json:"customer_name"`

	// Optional payment info for finance reporting.
	PaymentMethod string `json:"payment_method"`
	Received      *int64 `json:"received"`
	Change        *int64 `json:"change"`
}
