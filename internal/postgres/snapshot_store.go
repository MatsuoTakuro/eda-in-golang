package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
)

type snapshotStore struct {
	es.AggregateStore
	tableName string
	db        *sql.DB
	registry  registry.Registry
}

func WithSnapshotStore(tableName string, db *sql.DB, registry registry.Registry) es.AggregateStoreMiddleware {
	snapStore := snapshotStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}

	return func(aggStore es.AggregateStore) es.AggregateStore {
		snapStore.AggregateStore = aggStore
		return snapStore
	}
}

// Load loads the aggregate from the store and applies the snapshot to it.
func (s snapshotStore) Load(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	const query = `SELECT stream_version, snapshot_name, snapshot_data FROM %s WHERE stream_id = $1 AND stream_name = $2 LIMIT 1`

	var entityVersion int
	var snapshotName string
	var snapshotData []byte

	if err := s.db.QueryRowContext(ctx, s.table(query), aggregate.ID(), aggregate.AggregateName()).Scan(&entityVersion, &snapshotName, &snapshotData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.AggregateStore.Load(ctx, aggregate)
		}
		return err
	}

	v, err := s.registry.Deserialize(snapshotName, snapshotData, registry.ValidateImplements((*es.Snapshot)(nil)))
	if err != nil {
		return err
	}

	if err := es.LoadSnapshot(aggregate, v.(es.Snapshot), entityVersion); err != nil {
		return err
	}

	return s.AggregateStore.Load(ctx, aggregate)
}

// Save saves the aggregate to the store and then saves its snapshot if needed.
func (s snapshotStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	const query = `INSERT INTO %s (stream_id, stream_name, stream_version, snapshot_name, snapshot_data) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT (stream_id, stream_name) DO
UPDATE SET stream_version = EXCLUDED.stream_version, snapshot_name = EXCLUDED.snapshot_name, snapshot_data = EXCLUDED.snapshot_data`

	if err := s.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}

	if !s.shouldSnapshot(aggregate) {
		return nil
	}

	sser, ok := aggregate.(es.Snapshotter)
	if !ok {
		return fmt.Errorf("%T does not implelement es.Snapshotter", aggregate)
	}

	snapshot := sser.ToSnapshot()

	data, err := s.registry.Serialize(snapshot.SnapshotName(), snapshot)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, s.table(query), aggregate.ID(), aggregate.AggregateName(), aggregate.PendingVersion(), snapshot.SnapshotName(), data)

	return err
}

// TODO use injected & configurable strategies
func (snapshotStore) shouldSnapshot(aggregate es.EventSourcedAggregate) bool {
	var maxChanges = 3 // low for demonstration; production envs should use higher values 50, 75, 100...
	var pendingVersion = aggregate.PendingVersion()
	var pendingChanges = len(aggregate.Events())

	return pendingVersion >= maxChanges &&
		((pendingChanges >= maxChanges) ||
			(pendingVersion%maxChanges < pendingChanges) ||
			(pendingVersion%maxChanges == 0))
}

func (s snapshotStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
