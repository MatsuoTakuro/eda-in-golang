package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventPublisher interface {
	Publish(ctx context.Context, topicName string, event ddd.Event) error
}

type eventPublisher struct {
	reg          registry.Registry
	msgPublisher MessagePublisher
}

func NewEventPublisher(reg registry.Registry, msgPublisher MessagePublisher, mws ...MessagePublisherMiddleware) *eventPublisher {
	return &eventPublisher{
		reg:          reg,
		msgPublisher: messagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s *eventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	payload, err := s.reg.Serialize(event.EventName(), event.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
	})
	if err != nil {
		return err
	}
	return s.msgPublisher.Publish(ctx, topicName, message{
		id:       event.ID(),
		name:     event.EventName(),
		subject:  topicName,
		data:     data,
		metadata: event.Metadata(),
		sentAt:   time.Now(),
	})
}
