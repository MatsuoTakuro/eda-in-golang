package sec

import (
	"eda-in-golang/internal/am"
)

const (
	SagaCommandIDHdr   = am.CommandHdrPrefix + "SAGA_ID"
	SagaCommandNameHdr = am.CommandHdrPrefix + "SAGA_NAME"

	SagaReplyIDHdr   = am.ReplyHdrPrefix + "SAGA_ID"
	SagaReplyNameHdr = am.ReplyHdrPrefix + "SAGA_NAME"
)

type Saga[T any] interface {
	AddStep() Step[T]
	Name() string
	ReplyTopic() string
	getSteps() []Step[T]
}

type saga[T any] struct {
	name       string
	replyTopic string
	steps      []Step[T]
}

var _ Saga[any] = (*saga[any])(nil)

func NewSaga[T any](name, replyTopic string) *saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

const (
	notCompensating bool = false
	isCompensating  bool = true
)

func (s *saga[T]) AddStep() Step[T] {
	step := &step[T]{
		actions: map[bool]StepActionFunc[T]{
			notCompensating: nil,
			isCompensating:  nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc[T]{
			notCompensating: {},
			isCompensating:  {},
		},
	}

	s.steps = append(s.steps, step)

	return step
}

func (s *saga[T]) Name() string {
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga[T]) getSteps() []Step[T] {
	return s.steps
}
