package search

import (
	"context"

	"eda-in-golang/internal/system"
)

type Module struct{}

var _ system.Module = (*Module)(nil)

func (m Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}
