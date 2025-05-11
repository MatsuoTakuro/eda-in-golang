package infra

import (
	"context"
	"eda-in-golang/modules/ordering/internal/domain"
)

type ShoppingClient interface {
	Create(ctx context.Context, orderID string, items []domain.Item) (string, error)
	Cancel(ctx context.Context, shoppingID string) error
}
