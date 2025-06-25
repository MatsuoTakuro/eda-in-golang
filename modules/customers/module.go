package customers

import (
	"context"

	"eda-in-golang/internal/system"
)

type Module struct{}

var _ system.Module = (*Module)(nil)

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}
