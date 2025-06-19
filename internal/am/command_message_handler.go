package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"strings"

	"google.golang.org/protobuf/proto"
)

type commandMsgHandler struct {
	reg       registry.Registry
	publisher ReplyPublisher
	handler   ddd.CommandHandler[ddd.Command]
}

var _ RawMessageHandler = (*commandMsgHandler)(nil)

func NewCommandMessageHandler(
	reg registry.Registry,
	publisher ReplyPublisher,
	handler ddd.CommandHandler[ddd.Command],
) commandMsgHandler {
	return commandMsgHandler{
		reg:       reg,
		publisher: publisher,
		handler:   handler,
	}
}

// HandleMessage converts the raw message into a command message,
// deserializes the payload, and invokes the command handler.
func (h commandMsgHandler) HandleMessage(ctx context.Context, msg AckableRawMessage) error {
	var commandData CommandMessageData

	err := proto.Unmarshal(msg.Data(), &commandData)
	if err != nil {
		return err
	}

	commandName := msg.MessageName()

	payload, err := h.reg.Deserialize(commandName, commandData.GetPayload())
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

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, replyChannel, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, replyChannel, h.success(reply, commandMsg))
}

// publishReply publishes a reply message to the specified reply channel.
func (h commandMsgHandler) publishReply(ctx context.Context, replyChannel string, reply ddd.Reply) error {
	return h.publisher.Publish(ctx, replyChannel, reply)
}

// failure creates a failure reply.
func (h commandMsgHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeFailure)

	return h.applyCorrelationHeaders(reply, cmd)
}

// success creates a success reply.
func (h commandMsgHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeSuccess)

	return h.applyCorrelationHeaders(reply, cmd)
}

// applyCorrelationHeaders copies correlation headers from the command to the reply.
func (h commandMsgHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
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
