package sec

type Context[T any] struct {
	ID           string
	Data         T
	Step         int
	Done         bool
	Compensating bool
}

func (s *Context[T]) advance(steps int) {
	var dir = 1
	if s.Compensating {
		dir = -1
	}

	s.Step += dir * steps
}

func (s *Context[T]) complete() {
	s.Done = true
}

func (s *Context[T]) compensate() {
	s.Compensating = true
}
