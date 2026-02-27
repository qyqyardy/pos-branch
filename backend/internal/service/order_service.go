package service

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"

	"pos-backend/internal/model"
	"pos-backend/internal/repository"
)

type OrderService struct {
	OrderRepo   *repository.OrderRepo
	ProductRepo *repository.ProductRepo
}

func NewOrderService(orderRepo *repository.OrderRepo, productRepo *repository.ProductRepo) *OrderService {
	return &OrderService{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(userID uuid.UUID, req model.CreateOrderRequest) (uuid.UUID, int64, error) {

	if len(req.Items) == 0 {
		return uuid.Nil, 0, ErrEmptyCart
	}

	var items []repository.OrderItemData
	var total int64

	for _, i := range req.Items {

		if i.ProductID == uuid.Nil {
			return uuid.Nil, 0, ErrInvalidProduct
		}

		if i.Qty <= 0 {
			return uuid.Nil, 0, ErrInvalidQty
		}

		p, err := s.ProductRepo.FindByID(i.ProductID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return uuid.Nil, 0, ErrProductNotFound
			}
			return uuid.Nil, 0, err
		}

		total += int64(i.Qty) * p.Price

		items = append(items, repository.OrderItemData{
			ProductID: i.ProductID,
			Qty:       i.Qty,
			Price:     p.Price,
		})
	}

	// Metadata (optional)
	orderType := strings.TrimSpace(req.OrderType)
	if orderType == "" {
		orderType = "dine_in"
	}
	if orderType != "dine_in" && orderType != "take_away" {
		return uuid.Nil, 0, errors.New("invalid order_type")
	}

	payMethod := strings.TrimSpace(req.PaymentMethod)
	if payMethod == "" {
		payMethod = "cash"
	}
	if payMethod != "cash" && payMethod != "qris" && payMethod != "midtrans" {
		return uuid.Nil, 0, errors.New("invalid payment_method")
	}

	var received *int64
	var change *int64
	if payMethod == "cash" {
		rv := total
		if req.Received != nil {
			if *req.Received < total {
				return uuid.Nil, 0, errors.New("received is less than total")
			}
			rv = *req.Received
		}
		cv := rv - total
		received = &rv
		change = &cv
	}

	var paymentStatus string = "completed"
	var paymentToken string

	if payMethod == "midtrans" {
		paymentStatus = "pending"
		// Generate snap token (Simplified)
		token, err := s.generateMidtransToken(total)
		if err != nil {
			return uuid.Nil, 0, err
		}
		paymentToken = token
	}

	meta := repository.OrderMeta{
		OrderType:     orderType,
		TableNo:       req.TableNo,
		GuestCount:    req.GuestCount,
		CustomerName:  req.CustomerName,
		PaymentMethod: payMethod,
		Received:      received,
		Change:        change,
		PaymentStatus: paymentStatus,
		PaymentToken:  paymentToken,
	}

	orderID, err := s.OrderRepo.Create(userID, items, total, meta)
	if err != nil {
		return uuid.Nil, 0, err
	}
	return orderID, total, nil
}

func (s *OrderService) generateMidtransToken(total int64) (string, error) {
	// In a real app, use Midtrans SDK or direct HTTP call to Snap API here.
	// For this demo, we'll return a mock token if no key is found.
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return "mock_snap_token_" + uuid.New().String(), nil
	}

	// Pseudocode for direct HTTP call:
	// payload := map[string]interface{}{
	//     "transaction_details": map[string]interface{}{"order_id": uuid.New().String(), "gross_amount": total},
	// }
	// req, _ := http.NewRequest("POST", "https://app.sandbox.midtrans.com/snap/v1/transactions", body)
	// req.SetBasicAuth(serverKey, "")
	// ...

	return "mock_snap_token_" + uuid.New().String(), nil
}
