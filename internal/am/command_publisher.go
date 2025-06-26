package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandPublisher interface {
	// Publish sends a command message to the specified topic.
	Publish(ctx context.Context, topicName string, cmd ddd.Command) error
}

type commandPublisher struct {
	reg          registry.Registry
	msgPublisher MessagePublisher
}

func NewCommandPublisher(reg registry.Registry, msgPublisher MessagePublisher, mws ...MessagePublisherMiddleware) CommandPublisher {
	return commandPublisher{
		reg:          reg,
		msgPublisher: messagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

// Publish publishes a command message to the specified topic.
// It converts the command into a message format and publishes it.
func (s commandPublisher) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	payload, err := s.reg.Serialize(
		command.CommandName(), command.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return s.msgPublisher.Publish(ctx, topicName, message{
		id:       command.ID(),
		name:     command.CommandName(),
		subject:  topicName,
		data:     data,
		metadata: command.Metadata(),
		sentAt:   time.Now(),
	})
}
