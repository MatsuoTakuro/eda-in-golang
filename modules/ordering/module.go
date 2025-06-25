package ordering

import (
	"context"

	"eda-in-golang/internal/system"
)

type Module struct{}

var _ system.Module = (*Module)(nil)

// Startup starts up the ordering service as one of the modular monolith modules.
func (Module) Startup(ctx context.Context, svc system.Service) (err error) {
	return Root(ctx, svc)
}
