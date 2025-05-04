package ddd

type Entity interface {
	GetID() string
}

var _ Entity = (*EntityBase)(nil)

type EntityBase struct {
	ID string
}

func (e EntityBase) GetID() string {
	return e.ID
}
