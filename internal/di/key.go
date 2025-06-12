package di

// Key is a unique identifier for a dependency.
// You can add your application-specific types as keys.
type Key string

// common singleton keys
const (
	Registry         Key = "registry"
	Logger           Key = "logger"
	Stream           Key = "stream"
	DomainDispatcher Key = "domain_dispatcher"
	DB               Key = "sql_db"
	GRPCConn         Key = "grpc_conn"
	OutboxProcessor  Key = "outbox_processor"
)

// common scoped keys
const (
	TX                      Key = "sql_transaction"
	TXStream                Key = "sql_transaction_stream"
	EventStream             Key = "event_stream"
	CommandStream           Key = "command_stream"
	ReplyStream             Key = "reply_stream"
	InboxMiddleware         Key = "inbox_middleware"
	Application             Key = "application"
	DomainEventHandler      Key = "domain_event_handler"
	IntegrationEventHandler Key = "integration_event_handler"
	CommandHandler          Key = "command_handler"
)
