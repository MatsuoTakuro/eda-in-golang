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

var _ RawMessageHandler = (*eventMsgHandler)(nil)

func NewEventMessageHandler(reg registry.Registry, handler ddd.EventHandler[ddd.Event]) eventMsgHandler {
	return eventMsgHandler{
		reg:     reg,
		handler: handler,
	}
}

func (h eventMsgHandler) HandleMessage(ctx context.Context, msg AckableRawMessage) error {
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

	eventMsg := eventMessage{
		id:         msg.ID(),
		name:       eventName,
		payload:    payload,
		metadata:   eventData.GetMetadata().AsMap(),
		occurredAt: eventData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleEvent(ctx, eventMsg)
}
