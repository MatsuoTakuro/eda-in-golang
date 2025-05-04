package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application/eventhandlers"
)

type InvoiceEventHandler struct {
	invoice eventhandlers.Invoice
	logger  zerolog.Logger
}

func NewInvoiceEventHandler(
	invoice eventhandlers.Invoice,
	logger zerolog.Logger) InvoiceEventHandler {
	return InvoiceEventHandler{
		invoice: invoice,
		logger:  logger,
	}
}

func (h InvoiceEventHandler) OnOrderReadied(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Ordering.OnOrderReadied")
	defer func() { h.logger.Info().Err(err).Msg("<-- Ordering.OnOrderReadied") }()
	return h.invoice.OnOrderReadied(ctx, event)
}
