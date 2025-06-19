package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandStream interface {
	CommandPublisher
	CommandSubscriber
}

type CommandPublisher = MessagePublisher[ddd.Command]
type CommandSubscriber = MessageSubscriber[CommandMessage]

type commandStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var _ CommandStream = (*commandStream)(nil)

func NewCommandStream(reg registry.Registry, stream RawMessageStream) commandStream {
	return commandStream{
		reg:    reg,
		stream: stream,
	}
}

// Publish publishes a command message to the specified topic.
func (s commandStream) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	metadata, err := structpb.NewStruct(command.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(
		command.CommandName(), command.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      command.ID(),
		name:    command.CommandName(),
		subject: topicName,
		data:    data,
	})
}

// Subscribe subscribes to a topic for command messages and handles them using the provided handler.
func (s commandStream) Subscribe(topicName string, handler MessageHandler[CommandMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[AckableRawMessage](func(ctx context.Context, msg AckableRawMessage) error {
		var commandData CommandMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		commandName := msg.MessageName()

		payload, err := s.reg.Deserialize(commandName, commandData.GetPayload())
		if err != nil {
			return err
		}

		commandMsg := commandMessage{
			id:         msg.ID(),
			name:       commandName,
			payload:    payload,
			metadata:   commandData.GetMetadata().AsMap(),
			occurredAt: commandData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, commandMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (s commandStream) Unsubscribe() error {
	return s.stream.Unsubscribe()
}
