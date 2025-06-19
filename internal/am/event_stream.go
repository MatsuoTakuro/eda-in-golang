package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"errors"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	// EventPublisher is a publisher for events.
	EventPublisher = MessagePublisher[ddd.Event]
	// EventSubscriber is a subscriber for events.
	EventSubscriber = MessageSubscriber[EventMessage]
	// EventStream is a stream for events.
	EventStream = MessageStream[ddd.Event, EventMessage]
)

type eventStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var (
	_ EventStream = (*eventStream)(nil)
)

func NewEventStream(reg registry.Registry, stream RawMessageStream) *eventStream {
	return &eventStream{
		reg:    reg,
		stream: stream,
	}
}

func (s eventStream) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	metadata, err := structpb.NewStruct(event.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(
		event.EventName(), event.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      event.ID(),
		name:    event.EventName(),
		subject: topicName,
		data:    data,
	})
}

func (s eventStream) Subscribe(topicName string, handler MessageHandler[EventMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[AckableRawMessage](func(ctx context.Context, msg AckableRawMessage) error {
		var eventData EventMessageData

		// filter messages and only process those that match the filters
		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return ErrMessageSkipped
			}
		}

		// convert the raw message data to the event message data
		err := proto.Unmarshal(msg.Data(), &eventData)
		if err != nil {
			return err
		}

		eventName := msg.MessageName()

		payload, err := s.reg.Deserialize(eventName, eventData.GetPayload())
		if err != nil {
			return err
		}

		eventMsg := eventMessage{
			id:         msg.ID(),
			name:       eventName,
			payload:    payload,
			metadata:   eventData.GetMetadata().AsMap(),
			occurredAt: eventData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, eventMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

var ErrMessageSkipped = errors.New("message skipped due to filter")

func (s eventStream) Unsubscribe() error {
	return s.stream.Unsubscribe()
}
