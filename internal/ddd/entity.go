package ddd

type Entity interface {
	IDer
	EntityNamer
	IDSetter
	NameSetter
}

type IDer interface {
	ID() string
}

type EntityNamer interface {
	EntityName() string
}

var _ Entity = (*entity)(nil)

type entity struct {
	id   string
	name string
}

func NewEntity(id, name string) *entity {
	return &entity{
		id:   id,
		name: name,
	}
}

func (e entity) ID() string             { return e.id }
func (e entity) EntityName() string     { return e.name }
func (e entity) Equals(other IDer) bool { return e.id == other.ID() }
func (e *entity) setID(id string)       { e.id = id }
func (e *entity) setName(name string)   { e.name = name }
