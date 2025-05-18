package es

import (
	"fmt"
)

type Snapshot interface {
	SnapshotName() string
}

type Snapshotter interface {
	SnapshotApplier
	ToSnapshot() Snapshot
}

type SnapshotApplier interface {
	// ApplySnapshot applies a snapshot to an aggregate, depending on the type of the snapshot.
	ApplySnapshot(snapshot Snapshot) error
}

// LoadSnapshot applies a snapshot to an aggregate. It is used to rebuild the state of an aggregate from its snapshot.
func LoadSnapshot(v interface{}, snapshot Snapshot, version int) error {
	type loader interface {
		SnapshotApplier
		VersionSetter
	}

	agg, ok := v.(loader)
	if !ok {
		return fmt.Errorf("%T does not have the methods implemented to load snapshots", v)
	}

	if err := agg.ApplySnapshot(snapshot); err != nil {
		return err
	}
	agg.setVersion(version)

	return nil
}
