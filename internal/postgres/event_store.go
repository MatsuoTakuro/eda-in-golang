package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stackus/errors"

	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
)

type eventStore struct {
	tableName string
	db        DB
	registry  registry.Registry
}

var _ es.AggregateStore = (*eventStore)(nil)

func NewEventStore(tableName string, db DB, registry registry.Registry) eventStore {
	return eventStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}
}

// Load loads the aggregate from the event store and applies the events to it.
func (s eventStore) Load(ctx context.Context, aggregate es.EventSourcedAggregate) (err error) {
	const query = `SELECT stream_version, event_id, event_name, event_data, occurred_at FROM %s WHERE stream_id = $1 AND stream_name = $2 AND stream_version > $3 ORDER BY stream_version ASC`

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	var rows *sql.Rows

	rows, err = s.db.QueryContext(ctx, s.table(query), aggregateID, aggregateName, aggregate.Version())
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing event rows")
		}
	}(rows)

	for rows.Next() {
		var eventID, eventName string
		var payloadData []byte
		var aggregateVersion int
		var occurredAt time.Time
		err := rows.Scan(&aggregateVersion, &eventID, &eventName, &payloadData, &occurredAt)
		if err != nil {
			return err
		}

		var payload interface{}
		payload, err = s.registry.Deserialize(eventName, payloadData)
		if err != nil {
			return err
		}

		event := aggregateEvent{
			id:         eventID,
			name:       eventName,
			payload:    payload,
			aggregate:  aggregate,
			version:    aggregateVersion,
			occurredAt: occurredAt,
		}

		if err = es.LoadEvent(aggregate, event); err != nil {
			return err
		}
	}

	return nil
}

// Save saves all uncommitted events of the aggregate to the event store.
func (s eventStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) (err error) {
	const query = `INSERT INTO %s (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at) VALUES`

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	placeholders := make([]string, len(aggregate.Events()))
	values := make([]any, len(aggregate.Events())*7)

	for i, event := range aggregate.Events() {
		var payloadData []byte

		payloadData, err = s.registry.Serialize(event.EventName(), event.Payload())
		if err != nil {
			return err
		}

		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*7+1, // stream_id
			i*7+2, // stream_name
			i*7+3, // stream_version
			i*7+4, // event_id
			i*7+5, // event_name
			i*7+6, // event_data
			i*7+7, // occurred_at
		)

		values[i*7] = aggregateID
		values[i*7+1] = aggregateName
		values[i*7+2] = event.AggregateVersion()
		values[i*7+3] = event.ID()
		values[i*7+4] = event.EventName()
		values[i*7+5] = payloadData
		values[i*7+6] = event.OccurredAt()
	}
	// Here is an example query;
	// INSERT INTO ordering.events
	//   (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at)
	// VALUES
	//   ($1, $2, $3, $4, $5, $6, $7),
	//   ($8, $9, $10, $11, $12, $13, $14);
	if _, err = s.db.ExecContext(
		ctx,
		fmt.Sprintf("%s %s", s.table(query), strings.Join(placeholders, ",")),
		values...,
	); err != nil {
		return err
	}

	return nil
}

func (s eventStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
