package domain

type StoreV1 struct {
	Name          string
	Location      string
	Participating bool
}

func (StoreV1) SnapshotName() string { return "stores.StoreV1" }

func (s StoreV1) Key() string { return s.SnapshotName() }
