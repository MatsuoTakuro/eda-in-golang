package sec

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type Step[T any] interface {
	Action(fn StepActionFunc[T]) Step[T]
	Compensation(fn StepActionFunc[T]) Step[T]
	OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T]
	OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T]
	isInvocable(compensating bool) bool
	execute(ctx context.Context, sagaCtx *Context[T]) stepResult[T]
	handle(ctx context.Context, sagaCtx *Context[T], reply ddd.Reply) error
}

type StepActionFunc[T any] func(ctx context.Context, data T) am.Command
type StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error

type stepResult[T any] struct {
	ctx *Context[T]
	cmd am.Command
	err error
}
type step[T any] struct {
	actions  map[bool]StepActionFunc[T]
	handlers map[bool]map[string]StepReplyHandlerFunc[T]
}

var _ Step[any] = (*step[any])(nil)

func (s *step[T]) Action(fn StepActionFunc[T]) Step[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *step[T]) Compensation(fn StepActionFunc[T]) Step[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *step[T]) OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *step[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s step[T]) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s step[T]) execute(ctx context.Context, sagaCtx *Context[T]) stepResult[T] {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		return stepResult[T]{
			ctx: sagaCtx,
			cmd: action(ctx, sagaCtx.Data),
		}
	}

	return stepResult[T]{ctx: sagaCtx}
}

func (s step[T]) handle(ctx context.Context, sagaCtx *Context[T], reply ddd.Reply) error {
	if handler := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaCtx.Data, reply)
	}
	return nil
}
