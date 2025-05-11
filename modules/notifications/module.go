package notifications

import (
	"context"

	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/notifications/internal/application"
	"eda-in-golang/modules/notifications/internal/grpc"
	"eda-in-golang/modules/notifications/internal/logging"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Server) error {
	// setup Driven adapters
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := grpc.NewCustomerRepository(conn)

	// setup application
	var app application.App
	app = application.New(customers)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}

	return nil
}
