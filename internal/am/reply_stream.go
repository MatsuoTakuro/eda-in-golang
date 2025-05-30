package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	ReplyPublisher  = MessagePublisher[ddd.Reply]
	ReplySubscriber = MessageSubscriber[ReplyMessage]
	ReplyStream     = MessageStream[ddd.Reply, ReplyMessage]
)

type replyStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var _ ReplyStream = (*replyStream)(nil)

func NewReplyStream(reg registry.Registry, stream RawMessageStream) replyStream {
	return replyStream{
		reg:    reg,
		stream: stream,
	}
}

func (s replyStream) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {
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

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:   reply.ID(),
		name: reply.ReplyName(),
		data: data,
	})
}

func (s replyStream) Subscribe(topicName string, handler MessageHandler[ReplyMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[AckableRawMessage](func(ctx context.Context, msg AckableRawMessage) error {
		var replyData ReplyMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &replyData)
		if err != nil {
			return err
		}

		replyName := msg.MessageName()

		var payload any

		if replyName != SuccessReply && replyName != FailureReply {
			payload, err = s.reg.Deserialize(replyName, replyData.GetPayload())
			if err != nil {
				return err
			}
		}

		replyMsg := replyMessage{
			id:         msg.ID(),
			name:       replyName,
			payload:    payload,
			metadata:   replyData.GetMetadata().AsMap(),
			occurredAt: replyData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, replyMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}
