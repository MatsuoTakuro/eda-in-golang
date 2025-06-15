package domain

import "context"

type OrderRequestRepository interface {
	FindOrInsert(ctx context.Context, idemKey string, request OrderRequest, command any) (
		orderID string, inserted bool, err error,
	)
}

type OrderRequest string

const (
	CreateOrderRequest OrderRequest = "create_order_request"
)
