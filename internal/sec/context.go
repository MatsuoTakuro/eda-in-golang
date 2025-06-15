package sec

type Context[T any] struct {
	//ID is saga ID, which is unique for each saga instance.
	ID string
	// Data is usually the payload for the next command to be published.
	Data T
	// Step indicates the current step in the saga.
	Step int
	// Done indicates whether the saga has completed all its steps.
	Done bool
	// IsCompensating indicates whether the saga at the current step is compensating or not.
	IsCompensating isCompensating
}

func (s *Context[T]) advance(steps int) {
	var direction = 1
	if s.IsCompensating {
		direction = -1
	}

	s.Step += direction * steps
}

func (s *Context[T]) complete() {
	s.Done = true
}

func (s *Context[T]) markAsCompensating() {
	s.IsCompensating = true
}
