package postgres

import (
	"context"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/tm"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
)

type inboxStore struct {
	tableName string
	db        DB
}

var _ tm.InboxStore = (*inboxStore)(nil)

func NewInboxStore(tableName string, db DB) inboxStore {
	return inboxStore{
		tableName: tableName,
		db:        db,
	}
}

func (s inboxStore) Save(ctx context.Context, msg am.RawMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data, received_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)"

	_, err := s.db.ExecContext(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
	}

	return err
}

func (s inboxStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
