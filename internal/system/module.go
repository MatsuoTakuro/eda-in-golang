package system

import "context"

type Module interface {
	Startup(context.Context, Service) error
}
