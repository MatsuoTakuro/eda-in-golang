package am

import (
	"context"
)

type (
	RawMessageHandlerMiddleware = func(handler RawMessageHandler) RawMessageHandler
)

func RawMessageHandlerWithMiddleware(
	handler RawMessageHandler,
	mws ...RawMessageHandlerMiddleware,
) RawMessageHandler {
	h := handler
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type RawMessageHandlerFunc func(ctx context.Context, msg AckableRawMessage) error

func (f RawMessageHandlerFunc) HandleMessage(ctx context.Context, cmd AckableRawMessage) error {
	return f(ctx, cmd)
}
