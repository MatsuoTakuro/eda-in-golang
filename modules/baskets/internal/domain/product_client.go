package domain

import (
	"context"
)

type ProductClient interface {
	Find(ctx context.Context, productID string) (*Product, error)
}
