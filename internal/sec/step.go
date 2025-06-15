package sec

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type Step[T any] interface {
	// Normal adds a normal action.
	Normal(fn StepActionFunc[T]) Step[T]
	// Compensation adds a compensating action.
	Compensation(fn StepActionFunc[T]) Step[T]
	// OnNormalReply registers a reply handler for the normal action.
	OnNormalReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T]
	// OnCompensationReply registers a reply handler for the compensating action.
	OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) Step[T]
	// hasActionFor checks if the step has an action registered for the given isCompensating state.
	hasActionFor(isCompensating isCompensating) bool
	// executeAction executes the action for the step based on the isCompensating state.
	executeAction(ctx context.Context, sagaCtx *Context[T]) stepResult[T]
	// handleReply handles the reply for the step based on the isCompensating state.
	handleReply(ctx context.Context, sagaCtx *Context[T], reply ddd.Reply) error
}

type StepActionFunc[T any] func(ctx context.Context, data T) am.Command
type StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error

type stepResult[T any] struct {
	ctx *Context[T]
	// next holds the next command to be executed.
	next am.Command
	err  error
}

type step[T any] struct {
	// actions holds the action functions for both:
	// - normal action
	// - compensating action
	actions map[isCompensating]StepActionFunc[T]
	// replyHandlers holds the reply handlers for both:
	// - normal reply handler
	// - compensating reply handler
	replyHandlers map[isCompensating]map[string]StepReplyHandlerFunc[T]
}

// isCompensating indicates whether a step is isCompensating or not.
type isCompensating bool

const (
	normal       isCompensating = false
	compensating isCompensating = true
)

var _ Step[any] = (*step[any])(nil)

func (s *step[T]) Normal(action StepActionFunc[T]) Step[T] {
	s.actions[normal] = action
	return s
}

func (s *step[T]) Compensation(action StepActionFunc[T]) Step[T] {
	s.actions[compensating] = action
	return s
}

func (s *step[T]) OnNormalReply(replyName string, replyHandler StepReplyHandlerFunc[T]) Step[T] {
	s.replyHandlers[normal][replyName] = replyHandler
	return s
}

func (s *step[T]) OnCompensationReply(replyName string, replyHandler StepReplyHandlerFunc[T]) Step[T] {
	s.replyHandlers[compensating][replyName] = replyHandler
	return s
}

func (s step[T]) hasActionFor(isCompensating isCompensating) bool {
	return s.actions[isCompensating] != nil
}

func (s step[T]) executeAction(ctx context.Context, sagaCtx *Context[T]) stepResult[T] {
	if action := s.actions[sagaCtx.IsCompensating]; action != nil {
		return stepResult[T]{
			ctx:  sagaCtx,
			next: action(ctx, sagaCtx.Data),
		}
	}

	return stepResult[T]{ctx: sagaCtx}
}

func (s step[T]) handleReply(ctx context.Context, sagaCtx *Context[T], reply ddd.Reply) error {
	if replyHandler := s.replyHandlers[sagaCtx.IsCompensating][reply.ReplyName()]; replyHandler != nil {
		return replyHandler(ctx, sagaCtx.Data, reply)
	}
	return nil
}
