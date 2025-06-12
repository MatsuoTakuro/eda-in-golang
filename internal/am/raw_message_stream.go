package am

type (
	// RawMessageStream is a stream of raw messages that can be published and subscribed to.
	// It can be a stream for event, command, or reply messages.
	RawMessageStream     = MessageStream[RawMessage, AckableRawMessage]
	RawMessageHandler    = MessageHandler[AckableRawMessage]
	RawMessagePublisher  = MessagePublisher[RawMessage]
	RawMessageSubscriber = MessageSubscriber[AckableRawMessage]
)

type (
	RawMessageStreamMiddleware = func(stream RawMessageStream) RawMessageStream
)

func WithRawMessageStreamMiddlewares(
	stream RawMessageStream,
	mws ...RawMessageStreamMiddleware,
) RawMessageStream {
	s := stream
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}
	return s
}
