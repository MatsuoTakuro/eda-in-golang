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
	// AddStep returns a new step that can be added to the saga.
	AddStep() Step[T]
	// Name returns the name of the saga.
	Name() string
	// ReplyTopic returns the topic to which replies for this saga should be sent.
	ReplyTopic() string
	// getSteps returns the steps added to the saga.
	getSteps() []Step[T]
}

type saga[T any] struct {
	name       string
	replyTopic string
	// steps holds the steps to be executed in order how you added them.
	steps []Step[T]
}

var _ Saga[any] = (*saga[any])(nil)

func NewSaga[T any](name, replyTopic string) *saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

func (s *saga[T]) AddStep() Step[T] {
	step := &step[T]{
		actions: map[isCompensating]StepActionFunc[T]{
			normal:       nil,
			compensating: nil,
		},
		replyHandlers: map[isCompensating]map[string]StepReplyHandlerFunc[T]{
			normal:       {},
			compensating: {},
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
