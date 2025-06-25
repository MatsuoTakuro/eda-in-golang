package depot

import (
	"context"

	"eda-in-golang/internal/system"
)

type Module struct{}

var _ system.Module = (*Module)(nil)

func (Module) Startup(ctx context.Context, svc system.Service) error {
	return Root(ctx, svc)
}
