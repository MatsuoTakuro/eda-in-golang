package monolith

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/runner"
)

type Server interface {
	Config() config.AppConfig
	DB() *sql.DB
	// TODO: update the legacy API of nats.JetStreamContext to new JetStream API
	// (see: https://github.com/nats-io/nats.go/blob/main/jetstream/README.md)
	JS() nats.JetStreamContext
	Logger() zerolog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Runner() runner.Runner
}

type Module interface {
	Startup(context.Context, Server) error
}
