package infra

import (
	"context"
	"eda-in-golang/modules/ordering/internal/domain"
)

type ShoppingClient interface {
	Create(ctx context.Context, order *domain.Order) (string, error)
	Cancel(ctx context.Context, shoppingID string) error
}
