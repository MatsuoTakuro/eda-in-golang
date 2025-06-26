package am

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"

	"google.golang.org/protobuf/proto"
)

type eventMsgHandler struct {
	reg     registry.Registry
	handler ddd.EventHandler[ddd.Event]
}

var _ MessageHandler = (*eventMsgHandler)(nil)

func NewEventMessageHandler(
	reg registry.Registry,
	handler ddd.EventHandler[ddd.Event],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return messageHandlerWithMiddleware(eventMsgHandler{
		reg:     reg,
		handler: handler,
	}, mws...)
}

func (h eventMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var eventData EventMessageData

	err := proto.Unmarshal(msg.Data(), &eventData)
	if err != nil {
		return err
	}

	eventName := msg.MessageName()

	payload, err := h.reg.Deserialize(eventName, eventData.GetPayload())
	if err != nil {
		return err
	}

	// TODO either this should be a ddd.Event or the handler is a HandleMessage[am.EventMessage]
	eventMsg := eventMessage{
		id:         msg.ID(),
		name:       eventName,
		payload:    payload,
		occurredAt: eventData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleEvent(ctx, eventMsg)
}
