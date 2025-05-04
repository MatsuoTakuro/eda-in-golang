package infra

import (
	"context"
	"eda-in-golang/modules/ordering/internal/domain"
)

type OrderRepository interface {
	Find(ctx context.Context, orderID string) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
}
