package queries

import (
	"context"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Queries interface {
	GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error)
}

type queries struct {
	GetOrderHandler
}

func New(orderRepo infra.OrderRepository) Queries {
	return &queries{
		GetOrderHandler: NewGetOrderHandler(orderRepo),
	}
}
