package domain

import (
	"context"
)

type StoreClient interface {
	Find(ctx context.Context, storeID string) (*Store, error)
}
