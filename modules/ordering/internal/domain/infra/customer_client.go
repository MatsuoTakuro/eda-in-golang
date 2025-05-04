package infra

import (
	"context"
)

type CustomerClient interface {
	Authorize(ctx context.Context, customerID string) error
}
