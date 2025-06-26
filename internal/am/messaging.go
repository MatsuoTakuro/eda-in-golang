package am

import (
	"context"
	"errors"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg IncomingMessage) error
}

type MessageHandlerFunc func(ctx context.Context, msg IncomingMessage) error

var _ MessageHandler = MessageHandlerFunc(nil)

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	return f(ctx, msg)
}

type MessageSubscriber interface {
	Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (unsubscribe func() error, err error)
	UnsubscribeAll() error
}

type MessageStream interface {
	MessagePublisher
	MessageSubscriber
}

type MessagePublisher interface {
	Publish(ctx context.Context, topicName string, msg Message) error
}

var ErrMessageSkipped = errors.New("message skipped due to filter")
