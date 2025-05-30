package am

type (
	// RawMessageStream is a stream of raw messages that can be published and subscribed to.
	// It can be a stream for event, command, or reply messages.
	RawMessageStream     = MessageStream[RawMessage, AckableRawMessage]
	RawMessageHandler    = MessageHandler[AckableRawMessage]
	RawMessagePublisher  = MessagePublisher[RawMessage]
	RawMessageSubscriber = MessageSubscriber[AckableRawMessage]
)
