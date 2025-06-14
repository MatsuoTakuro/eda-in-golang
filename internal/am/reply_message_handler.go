package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"

	"google.golang.org/protobuf/proto"
)

type replyMsgHandler struct {
	reg     registry.Registry
	handler ddd.ReplyHandler[ddd.Reply]
}

var _ RawMessageHandler = (*replyMsgHandler)(nil)

func NewReplyMessageHandler(reg registry.Registry, handler ddd.ReplyHandler[ddd.Reply]) replyMsgHandler {
	return replyMsgHandler{
		reg:     reg,
		handler: handler,
	}
}

func (h replyMsgHandler) HandleMessage(ctx context.Context, msg AckableRawMessage) error {
	var replyData ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {
		return err
	}

	replyName := msg.MessageName()

	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.GetPayload())
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

	return h.handler.HandleReply(ctx, replyMsg)
}
