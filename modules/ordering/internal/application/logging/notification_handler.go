package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application/eventhandlers"
)

type NotificaionEventHandler struct {
	notificaion eventhandlers.Notification
	logger      zerolog.Logger
}

func NewNotificationEventHandler(
	notificaion eventhandlers.Notification,
	logger zerolog.Logger) NotificaionEventHandler {
	return NotificaionEventHandler{
		notificaion: notificaion,
		logger:      logger,
	}
}

func (h NotificaionEventHandler) OnOrderCreated(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Ordering.OnOrderCreated")
	defer func() { h.logger.Info().Err(err).Msg("<-- Ordering.OnOrderCreated") }()
	return h.notificaion.OnOrderCreated(ctx, event)
}

func (h NotificaionEventHandler) OnOrderReadied(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Ordering.OnOrderReadied")
	defer func() { h.logger.Info().Err(err).Msg("<-- Ordering.OnOrderReadied") }()
	return h.notificaion.OnOrderReadied(ctx, event)
}

func (h NotificaionEventHandler) OnOrderCanceled(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Ordering.OnOrderCanceled")
	defer func() { h.logger.Info().Err(err).Msg("<-- Ordering.OnOrderCanceled") }()
	return h.notificaion.OnOrderCanceled(ctx, event)
}
