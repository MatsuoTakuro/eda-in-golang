package am

import (
	"context"

	"eda-in-golang/internal/ddd"
)

// Message is a message interface that represents a message in the system.
// It can be an event, command, or any other type of message.
type Message interface {
	ddd.IDer
	MessageName() string
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type MessageHandler[M Message] interface {
	HandleMessage(ctx context.Context, msg M) error
}

// MessagePublisher is a publisher for messages.
type MessagePublisher[T any] interface {
	Publish(ctx context.Context, topicName string, v T) error
}

// MessageSubscriber is a subscriber for messages.
type MessageSubscriber[M Message] interface {
	Subscribe(topicName string, handler MessageHandler[M], options ...SubscriberOption) error
}

// MessageStream is a stream for messages.
type MessageStream[T any, M Message] interface {
	MessagePublisher[T]
	MessageSubscriber[M]
}

type MessageHandlerFunc[M Message] func(ctx context.Context, msg M) error

func (f MessageHandlerFunc[M]) HandleMessage(ctx context.Context, msg M) error {
	return f(ctx, msg)
}
