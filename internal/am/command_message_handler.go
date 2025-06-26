package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"strings"

	"google.golang.org/protobuf/proto"
)

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandNameHdr         = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"

	ReplyHdrPrefix  = "REPLY_"
	ReplyNameHdr    = ReplyHdrPrefix + "NAME"
	ReplyOutcomeHdr = ReplyHdrPrefix + "OUTCOME"
)

type commandMsgHandler struct {
	reg            registry.Registry
	replyPublisher ReplyPublisher
	cmdHandler     ddd.CommandHandler[ddd.Command]
}

var _ MessageHandler = (*commandMsgHandler)(nil)

func NewCommandMessageHandler(
	reg registry.Registry,
	replyPublisher ReplyPublisher,
	cmdHandler ddd.CommandHandler[ddd.Command],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return messageHandlerWithMiddleware(commandMsgHandler{
		reg:            reg,
		replyPublisher: replyPublisher,
		cmdHandler:     cmdHandler,
	}, mws...)
}

// HandleMessage converts the raw message into a command message,
// deserializes the payload, and invokes the command handler.
func (h commandMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
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
		occurredAt: commandData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	replyChannel := commandMsg.Metadata().Get(CommandReplyChannelHdr).(string)

	reply, err := h.cmdHandler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, replyChannel, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, replyChannel, h.success(reply, commandMsg))
}

// publishReply publishes a reply message to the specified reply channel.
func (h commandMsgHandler) publishReply(ctx context.Context, replyChannel string, reply ddd.Reply) error {
	return h.replyPublisher.Publish(ctx, replyChannel, reply)
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
