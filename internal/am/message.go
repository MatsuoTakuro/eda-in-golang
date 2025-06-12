package am

import (
	"context"

	"eda-in-golang/internal/ddd"
)

// AckableMessage represents a message that supports acknowledgment control.
type AckableMessage interface {
	Message
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

// Message is a message interface that represents a message in the system.
// It can be an event, command, or any other type of message.
type Message interface {
	// ID returns the unique identifier of the message.
	ddd.IDer
	// Subject returns the subject (topic) of the message.
	Subject() string
	// MessageName returns the type name of the message.
	MessageName() string
}

type MessageHandler[I AckableMessage] interface {
	HandleMessage(ctx context.Context, msg I) error
}

// MessagePublisher is a publisher for messages.
type MessagePublisher[T any] interface {
	Publish(ctx context.Context, topicName string, v T) error
}

// MessageSubscriber is a subscriber for messages.
type MessageSubscriber[A AckableMessage] interface {
	Subscribe(topicName string, handler MessageHandler[A], options ...SubscriberOption) error
}

// MessageStream is a stream for messages.
type MessageStream[T any, A AckableMessage] interface {
	MessagePublisher[T]
	MessageSubscriber[A]
}

type MessageHandlerFunc[A AckableMessage] func(ctx context.Context, msg A) error

var _ MessageHandler[AckableMessage] = MessageHandlerFunc[AckableMessage](nil)

func (f MessageHandlerFunc[A]) HandleMessage(ctx context.Context, msg A) error {
	return f(ctx, msg)
}
