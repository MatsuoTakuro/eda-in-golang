package queries

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type GetOrder struct {
	ID string
}

type GetOrderHandler struct {
	orderRepo infra.OrderRepository
}

func NewGetOrderHandler(orderRepo infra.OrderRepository) GetOrderHandler {
	return GetOrderHandler{orderRepo: orderRepo}
}

func (h GetOrderHandler) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.orderRepo.Find(ctx, query.ID)

	return order, errors.Wrap(err, "get order query")
}
