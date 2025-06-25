package system

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/runner"
)

// Service represents a microservice with access to shared infrastructure and common communication layers.
type Service interface {
	Config() config.AppConfig
	DB() *sql.DB
	Nats() *nats.Conn
	JS() jetstream.JetStream
	Mux() *chi.Mux
	RPC() *grpc.Server
	Runner() runner.Runner
	RunFuncs
	Logger() zerolog.Logger
}

type RunFuncs interface {
	RunWeb(context.Context) error
	RunRPC(context.Context) error
	RunStream(context.Context) error
}
