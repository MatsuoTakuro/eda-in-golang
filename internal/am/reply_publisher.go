package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReplyPublisher interface {
	Publish(ctx context.Context, topicName string, reply ddd.Reply) error
}

type replyPublisher struct {
	reg       registry.Registry
	publisher MessagePublisher
}

var _ ReplyPublisher = (*replyPublisher)(nil)

func NewReplyPublisher(reg registry.Registry, msgPublisher MessagePublisher, mws ...MessagePublisherMiddleware) ReplyPublisher {
	return &replyPublisher{
		reg:       reg,
		publisher: messagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s replyPublisher) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {
	var err error
	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       reply.ID(),
		name:     reply.ReplyName(),
		subject:  topicName,
		data:     data,
		metadata: reply.Metadata(),
		sentAt:   time.Now(),
	})
}
