package postgres

import (
	"context"
	"database/sql"
)

// DB defines a simplified interface that both *sql.DB and *sql.Tx from the database/sql package implement.
// It lets you write functions that work with either a normal DB connection or a transaction.
type DB interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
