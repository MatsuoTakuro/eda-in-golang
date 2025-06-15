package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandStream interface {
	CommandPublisher
	CommandSubscriber
}

type CommandPublisher = MessagePublisher[ddd.Command]

type CommandSubscriber interface {
	Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error
}

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
func (s commandStream) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
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

		replyChannel := commandMsg.Metadata().Get(CommandReplyChannelHdr).(string)

		var reply ddd.Reply
		reply, err = handler.HandleMessage(ctx, commandMsg)
		if err != nil {
			return s.publishReply(ctx, replyChannel, s.failure(reply, commandMsg))
		}

		return s.publishReply(ctx, replyChannel, s.success(reply, commandMsg))
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

// publishReply publishes a reply message to the specified reply channel.
func (s commandStream) publishReply(ctx context.Context, replyChannel string, reply ddd.Reply) error {
	metadata, err := structpb.NewStruct(reply.Metadata())
	if err != nil {
		return err
	}

	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(
			reply.ReplyName(), reply.Payload(),
		)
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, replyChannel, rawMessage{
		id:      reply.ID(),
		name:    reply.ReplyName(),
		subject: replyChannel,
		data:    data,
	})
}

// failure creates a failure reply.
func (s commandStream) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeFailure)

	return s.applyCorrelationHeaders(reply, cmd)
}

// success creates a success reply.
func (s commandStream) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeSuccess)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandStream) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHdr {
			continue
		}

		if strings.HasPrefix(key, CommandHdrPrefix) {
			hdr := ReplyHdrPrefix + key[len(CommandHdrPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}
