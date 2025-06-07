package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/sec"
	"eda-in-golang/modules/cosec/internal/models"
)

func SubscribeReplies(subscriber am.ReplySubscriber, orchestrator sec.Orchestrator[*models.CreateOrderData]) error {
	replyMsgHandler := am.MessageHandlerFunc[am.ReplyMessage](func(ctx context.Context, replyMsg am.ReplyMessage) error {
		return orchestrator.HandleReply(ctx, replyMsg)
	})
	return subscriber.Subscribe(orchestrator.ReplyTopic(), replyMsgHandler, am.GroupName("cosec-replies"))
}
