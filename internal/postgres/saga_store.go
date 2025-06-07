package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/sec"
)

type sagaStore struct {
	tableName string
	db        *sql.DB
	registry  registry.Registry
}

var _ sec.Store = (*sagaStore)(nil)

func NewSagaStore(tableName string, db *sql.DB, registry registry.Registry) sagaStore {
	return sagaStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}
}

func (s sagaStore) Load(ctx context.Context, sagaName, sagaID string) (*sec.Context[[]byte], error) {
	const query = "SELECT data, step, done, compensating FROM %s WHERE name = $1 AND id = $2"

	sagaCtx := &sec.Context[[]byte]{
		ID: sagaID,
	}
	err := s.db.QueryRowContext(ctx, s.table(query), sagaName, sagaID).
		Scan(&sagaCtx.Data, &sagaCtx.Step, &sagaCtx.Done, &sagaCtx.IsCompensating)

	return sagaCtx, err
}

func (s sagaStore) Save(ctx context.Context, sagaName string, sagaCtx *sec.Context[[]byte]) error {
	// Upsert saga state: insert if new, or update existing saga context by name and ID
	const query = `INSERT INTO %s (name, id, data, step, done, compensating) 
VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (name, id) DO
UPDATE SET data = EXCLUDED.data, step = EXCLUDED.step, done = EXCLUDED.done, compensating = EXCLUDED.compensating`

	_, err := s.db.ExecContext(ctx, s.table(query), sagaName, sagaCtx.ID, sagaCtx.Data, sagaCtx.Step, sagaCtx.Done, sagaCtx.IsCompensating)

	return err
}

func (s sagaStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
