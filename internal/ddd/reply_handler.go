package ddd

import "context"

type ReplyHandler[T Reply] interface {
	HandleReply(ctx context.Context, reply T) error
}

type ReplyHandlerFunc[T Reply] func(ctx context.Context, reply T) error

var _ ReplyHandler[Reply] = (*ReplyHandlerFunc[Reply])(nil)

func (f ReplyHandlerFunc[T]) HandleReply(ctx context.Context, reply T) error {
	return f(ctx, reply)
}
