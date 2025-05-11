package ddd

type EventOption interface {
	apply(*event)
}
